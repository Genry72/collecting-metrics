package usecases

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type Agent struct {
	httpClient *http.Client
	hostPort   string
}

type metricer interface {
	getUrlsMetric() chan string
}

func NewAgent(hostPort string) *Agent {
	return &Agent{
		httpClient: &http.Client{},
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
	request, err := http.NewRequest("POST", a.hostPort+url, nil)
	if err != nil {
		return err
	}

	resp, err := a.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s :%s", resp.Status, string(body))
	}

	return nil
}
