package usecases

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

type metricer interface {
	getUrlsMetric() chan string
}

func NewAgent(hostPort string) *Agent {
	return &Agent{
		httpClient: resty.New(),
		hostPort:   hostPort,
	}
}

// SendMetrics Отправка метрик с заданным интервалом
func (a *Agent) SendMetrics(metric metricer, reportInterval time.Duration) {
	for {
		time.Sleep(reportInterval)
		for url := range metric.getUrlsMetric() {
			if err := a.send(url); err != nil {
				fmt.Println(err)
			}
		}
		log.Println("metrics send success")
	}
}

func (a *Agent) send(url string) error {
	resp, err := a.httpClient.R().Post(a.hostPort + url)
	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("%s :%s", resp.Status(), string(resp.Body()))
	}

	return nil
}
