package agent

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/usecases/agent/grpcClient"
	"github.com/Genry72/collecting-metrics/internal/usecases/agent/httpclients"
	"github.com/Genry72/collecting-metrics/internal/usecases/cryptor"
	interceptor "github.com/Genry72/collecting-metrics/internal/usecases/interceptor/agent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Agent struct {
	client        SenderMetrics
	log           *zap.Logger
	keyHash       *string
	publicKey     *rsa.PublicKey
	ratelimitChan chan struct{} // Количество одновременно исходящих запросов на сервер
}

// NewAgent Получение агента для сбора и отправки метрик
func NewAgent(hostPort string, grpcHostPort *string, log *zap.Logger, keyHash *string, publicKeyPath *string, rateLimitPtr *int) (*Agent, error) {
	var client SenderMetrics
	// устанавливаем соединение с сервером
	var (
		grpcconn  *grpc.ClientConn
		publicLey *rsa.PublicKey
		err       error
	)

	if publicKeyPath != nil && *publicKeyPath != "" {
		publicLey, err = cryptor.GetPubKeyFromFile(*publicKeyPath)
		if err != nil {
			return nil, fmt.Errorf("cryptor.GetPubKeyFromFile: %w", err)
		}
	}

	if grpcHostPort != nil && *grpcHostPort != "" {
		interceptors := make([]grpc.UnaryClientInterceptor, 0)

		// логирование запросов
		interceptors = append(interceptors, interceptor.Logging(log))

		// Передача ip адреса в метаданных
		interceptors = append(interceptors, interceptor.SetIpToHeader(log))

		// шифрование тела запроса
		if publicLey != nil {
			interceptors = append(interceptors, interceptor.EncryptBodyWithPublicKey(log, publicLey))
		}

		// Добавление метаданных с хешем тела запроса
		if keyHash != nil && *keyHash != "" {
			interceptors = append(interceptors, interceptor.SetHashToHeader(log, *keyHash))
		}

		grpcconn, err = grpc.Dial(*grpcHostPort, grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithChainUnaryInterceptor(interceptors...))

		if err != nil {
			log.Fatal("grpc.Dial", zap.Error(err))
		}
		client, err = grpcClient.NewGrpcClient(grpcconn, log, keyHash, publicKeyPath)
	}

	if grpcconn == nil {
		client, err = httpclients.NewHTTPClient(hostPort, log, keyHash, publicKeyPath)
	}

	rateLimit := 0
	if rateLimitPtr == nil || *rateLimitPtr == 0 {
		rateLimit = 1
	} else {
		rateLimit = *rateLimitPtr
	}

	return &Agent{
		client:        client,
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

			if err := a.client.Stop(); err != nil {
				a.log.Error("a.connGrpc.Close", zap.Error(err))
			}

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

			if err := a.send(ctx, metrics); err != nil {
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
func (a *Agent) send(ctx context.Context, metric models.Metrics) error {
	select {
	case <-ctx.Done():
		return nil
	case a.ratelimitChan <- struct{}{}:
		defer func() {
			<-a.ratelimitChan
		}()

		var (
			rErr error
		)

		// Индекс - количество выполненных повторов. Значение пауза в секундах
		retry := []time.Duration{0, 1, 3, 5}

		for i := 0; i < len(retry); i++ {
			sleepTime := retry[i]
			time.Sleep(sleepTime * time.Second)

			if err := a.client.Send(ctx, metric); err != nil {
				rErr = err
				var e *models.RetryError
				if errors.As(err, &e) {
					// ошибка, при которой нужно повторить запрос
					continue
				} else {
					return err
				}
			}
		}

		//a.log.Info("metrics send success")

		return rErr
	}

}

type SenderMetrics interface {
	Send(ctx context.Context, metric models.Metrics) error
	Stop() error
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
