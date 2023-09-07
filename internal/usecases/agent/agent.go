package agent

import (
	"errors"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"time"
)

type Agent struct {
	httpClient *resty.Client
	hostPort   string
	log        *zap.Logger
}

func NewAgent(hostPort string, log *zap.Logger) *Agent {
	restyClient := resty.New()
	restyClient.SetTimeout(time.Second)
	return &Agent{
		httpClient: restyClient,
		hostPort:   hostPort,
		log:        log,
	}
}

// SendMetrics Отправка метрик с заданным интервалом
func (a *Agent) SendMetrics(metric *Metrics, reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval)
		metrics, err := metric.getMetrics()
		if err != nil {
			a.log.Error(err.Error())
		}

		if err := a.sendByJSONBatch(metrics); err != nil {
			a.log.Error("sendByJSONBatch", zap.Error(err))
			continue
		}
		a.log.Info("metrics send success")
	}
}

func (a *Agent) sendByJSONBatch(metric []*models.Metric) error {
	url := "/updates"

	// Индекс - количество выполненных повторов. Значение пауза в секундах
	retry := []time.Duration{0, 1, 3, 5}

	var (
		rErr error
	)

	for i := 0; i < len(retry); i++ {
		sleepTime := retry[i]
		time.Sleep(sleepTime * time.Second)

		resp, err := a.httpClient.R().SetBody(metric).Post(a.hostPort + url)
		if err != nil {
			a.log.Error(err.Error())
			// или сеть или тело ответа
			continue
		}

		if err := checkStatus(resp.StatusCode()); err != nil {
			a.log.Error(err.Error())
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

	return rErr
}

func checkStatus(statusCode int) error {
	switch {
	case statusCode >= 200 && statusCode < 400:
		return nil
	case statusCode >= 400 && statusCode < 500:
		// повтор не нужен
		return fmt.Errorf("status not ok: %d", statusCode)
	case statusCode >= 500:
		// нужен повтор
		err := fmt.Errorf("status not ok: %d", statusCode)
		return models.NewRetryError(err)
	default:
		return nil
	}
}
