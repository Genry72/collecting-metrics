package agent

import (
	"context"
	"errors"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
)

type Agent struct {
	httpClient    *resty.Client
	hostPort      string
	log           *zap.Logger
	keyHash       *string
	ratelimitChan chan struct{} // Количество одновременно исходящих запросов на сервер
}

func NewAgent(hostPort string, log *zap.Logger, keyHash *string, rateLimit uint64) *Agent {
	restyClient := resty.New()
	restyClient.SetTimeout(time.Second)
	return &Agent{
		httpClient:    restyClient,
		hostPort:      hostPort,
		log:           log,
		keyHash:       keyHash,
		ratelimitChan: make(chan struct{}, rateLimit),
	}
}

// SendMetrics Отправка метрик с заданным интервалом
func (a *Agent) SendMetrics(ctx context.Context, metric *Metrics, reportInterval time.Duration) {
	go func() {
		for {
			time.Sleep(reportInterval)

			select {
			case <-ctx.Done():
				a.log.Info("Stop SendMetrics process")
				return
			default:
			}

			go func() {
				metrics, err := metric.getMetrics()
				if err != nil {
					a.log.Error("metric.getMetrics", zap.Error(err))
				}

				if len(metrics) == 0 {
					return
				}

				if err := a.sendByJSONBatch(ctx, metrics); err != nil {
					a.log.Error("sendByJSONBatch", zap.Error(err))
					return
				}

			}()

		}

	}()

}

func (a *Agent) sendByJSONBatch(ctx context.Context, metric models.Metrics) error {

	defer func() {
		<-a.ratelimitChan
	}()

	select {
	case <-ctx.Done():
		close(a.ratelimitChan)
		a.log.Info("Stop sendByJSONBatch process")
		return nil
	case a.ratelimitChan <- struct{}{}:
		url := "/updates"

		// Индекс - количество выполненных повторов. Значение пауза в секундах
		retry := []time.Duration{0, 1, 3, 5}

		var (
			rErr error
		)

		for i := 0; i < len(retry); i++ {
			sleepTime := retry[i]
			time.Sleep(sleepTime * time.Second)

			// Добавляем заголовок с хешем тела запроса, если передан ключ
			if a.keyHash != nil {
				hash, err := metric.Encode(*a.keyHash)
				if err != nil {
					return fmt.Errorf("metric.Encode: %w", err)
				}
				a.httpClient.R().SetHeader(models.HeaderHash, hash)
			}
			client := a.httpClient.R().SetContext(ctx)
			resp, err := client.SetBody(metric).Post(a.hostPort + url)
			if err != nil {
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

		a.log.Info("metrics send success")

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
