package agent

import (
	"fmt"
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
		for _, url := range metric.getUrlsMetric() {
			if err := a.sendByURL(url); err != nil {
				a.log.Error(err.Error())
			}
		}
		a.log.Info("metrics send success")
	}
}

func (a *Agent) sendByURL(url string) error {
	resp, err := a.httpClient.R().Post(a.hostPort + url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%s :%s", resp.Status(), string(resp.Body()))
	}

	return nil
}
