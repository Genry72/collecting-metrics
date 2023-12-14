package agent

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/usecases/cryptor"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"time"
)

type Agent struct {
	httpClient    *resty.Client
	hostPort      string
	log           *zap.Logger
	keyHash       *string
	publicKey     *rsa.PublicKey
	ratelimitChan chan struct{} // Количество одновременно исходящих запросов на сервер
}

// NewAgent Получение агента для сбора и отправки метрик
func NewAgent(hostPort string, log *zap.Logger, keyHash *string, publicKeyPath string, rateLimit uint64) (*Agent, error) {
	restyClient := resty.New()

	restyClient.SetTimeout(time.Second)

	var (
		publicLey *rsa.PublicKey
		err       error
	)

	if publicKeyPath != "" {
		publicLey, err = cryptor.GetPubKeyFromFile(publicKeyPath)
		if err != nil {
			return nil, fmt.Errorf("cryptor.GetPubKeyFromFile: %w", err)
		}
	}

	return &Agent{
		httpClient:    restyClient,
		hostPort:      hostPort,
		log:           log,
		keyHash:       keyHash,
		publicKey:     publicLey,
		ratelimitChan: make(chan struct{}, rateLimit),
	}, nil
}

// SendMetrics Отправка метрик с заданным интервалом
func (a *Agent) SendMetrics(ctx context.Context, metric *Metrics, reportInterval time.Duration) {
	defer func() {
		close(a.ratelimitChan)
	}()
	t := time.NewTicker(reportInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			//a.log.Info("Stop SendMetrics process")
			return
		case <-t.C:
			metrics, err := metric.getMetrics()
			if err != nil {
				a.log.Error("metric.getMetrics", zap.Error(err))
				return
			}

			if len(metrics) == 0 {
				return
			}

			if err := a.sendByJSONBatch(ctx, metrics); err != nil {
				a.log.Error("sendByJSONBatch", zap.Error(err))
				return
			}

		}

	}

}

/*
sendByJSONBatch отправляет метрики через HTTP POST запросы в формате JSON.
Функция использует рейт-лимит для ограничения количества запросов в единицу времени.
Если задан ключ, то функция добавляет заголовок с хешем тела запроса.
Функция выполняет повторные запросы в случае ошибки, используя заданные интервалы повторов.
Возвращает ошибку, если все повторные запросы неудачны или если статус ответа не является успешным.
*/
func (a *Agent) sendByJSONBatch(ctx context.Context, metric models.Metrics) error {
	select {
	case <-ctx.Done():
		return nil
	case a.ratelimitChan <- struct{}{}:
		defer func() {
			<-a.ratelimitChan
		}()
		url := "/updates"

		// Индекс - количество выполненных повторов. Значение пауза в секундах
		retry := []time.Duration{0, 1, 3, 5}

		var (
			rErr error
		)

		for i := 0; i < len(retry); i++ {
			sleepTime := retry[i]
			time.Sleep(sleepTime * time.Second)

			client := a.httpClient.R().SetContext(ctx)

			json := jsoniter.ConfigCompatibleWithStandardLibrary

			metricJSON, err := json.Marshal(metric)
			if err != nil {
				return err
			}
			// Шифроуем тело запроса, если передан публичный ключ
			if a.publicKey != nil {
				metricJSON, err = cryptor.EncryptBodyWithPublicKey(metricJSON, a.publicKey)
				if err != nil {
					return fmt.Errorf("cryptor.EncryptBodyWithPublicKey: %w", err)
				}
			}

			// Добавляем заголовок с хешем тела запроса, если передан ключ
			if a.keyHash != nil {
				hash, err := cryptor.Encrypt(metricJSON, *a.keyHash)
				if err != nil {
					return fmt.Errorf("cryptor.Encrypt: %w", err)
				}

				client.SetHeader(models.HeaderHash, hash)
			}

			resp, err := client.SetBody(metricJSON).Post(a.hostPort + url)
			if err != nil {
				if ctx.Err() != nil {
					return nil
				}
				a.log.Error("resp", zap.Error(err))
				// или сеть или тело ответа
				continue
			}

			if err := checkStatus(resp.StatusCode(), string(resp.Body())); err != nil {
				a.log.Error("checkStatus", zap.Error(err))
				rErr = err
				var e *models.RetryError
				if errors.As(err, &e) {
					// ошибка, при которой нужно повторить запрос
					continue
				}

				return err
			}
			// если дошли до сюда, то запрос выполнился корректно
			rErr = nil
			break
		}

		//a.log.Info("metrics send success")

		return rErr
	}

}

func checkStatus(statusCode int, body string) error {
	switch {
	case statusCode >= 200 && statusCode < 400:
		return nil
	case statusCode >= 400 && statusCode < 500:
		// повтор не нужен
		return fmt.Errorf("status not ok: %d body: %s", statusCode, body)
	case statusCode >= 500:
		// нужен повтор
		err := fmt.Errorf("status not ok: %d body: %s", statusCode, body)
		return models.NewRetryError(err)
	default:
		return nil
	}
}
