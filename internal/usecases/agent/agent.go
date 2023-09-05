package agent

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
	"net/http"
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
		//for i := range metrics {
		//	if err := a.sendByURL(metrics[i]); err != nil {
		//		a.log.Error(err.Error())
		//	}
		//}
		if err := a.sendByJSONBatch(metrics); err != nil {
			a.log.Error(err.Error())
		}
		a.log.Info("metrics send success")
	}
}

func (a *Agent) sendByURL(metric *models.Metric) error {
	var url string

	if metric == nil {
		return models.ErrBadMetricType
	}

	switch metric.MType {
	case models.MetricTypeGauge:
		url = fmt.Sprintf("/update/%s/%s/%v", models.MetricTypeGauge, metric.ID, *metric.Value)
	case models.MetricTypeCounter:
		url = fmt.Sprintf("/update/%s/%s/%v", models.MetricTypeCounter, metric.ID, *metric.Delta)
	}
	resp, err := a.httpClient.R().Post(a.hostPort + url)
	if err != nil {
		a.log.Error(err.Error())
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("%s :%s", resp.Status(), string(resp.Body()))
		a.log.Error(err.Error())
		return err
	}

	return nil
}

func (a *Agent) sendByJSONBatch(metric []*models.Metric) error {
	url := "/updates"
	resp, err := a.httpClient.R().SetBody(metric).Post(a.hostPort + url)
	if err != nil {
		a.log.Error(err.Error())
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("%s :%s", resp.Status(), string(resp.Body()))
		a.log.Error(err.Error())
		return err
	}

	return nil
}
