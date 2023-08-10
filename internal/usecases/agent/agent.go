package agent

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"time"
)

type Agent struct {
	httpClient *resty.Client
	hostPort   string
}

func NewAgent(hostPort string) *Agent {
	restyClient := resty.New()
	restyClient.SetTimeout(time.Second)
	return &Agent{
		httpClient: restyClient,
		hostPort:   hostPort,
	}
}

// SendMetrics Отправка метрик с заданным интервалом
func (a *Agent) SendMetrics(metric *Metrics, reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval)
		for _, url := range metric.getUrlsMetric() {
			if err := a.sendByURL(url); err != nil {
				fmt.Println(err)
			}
		}
		log.Println("metrics send success")
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
