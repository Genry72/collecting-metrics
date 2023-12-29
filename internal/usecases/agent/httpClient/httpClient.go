package httpClient

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/Genry72/collecting-metrics/helpers"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/usecases/cryptor"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
	"time"
)

type HttpClient struct {
	httpClient *resty.Client
	hostPort   string
	log        *zap.Logger
	keyHash    *string
	publicKey  *rsa.PublicKey
}

func NewHttpClient(hostPort string, log *zap.Logger, keyHash *string, publicKeyPath *string) (*HttpClient, error) {
	restyClient := resty.New()
	restyClient.SetTimeout(time.Second)

	var (
		publicLey *rsa.PublicKey
		err       error
	)

	if publicKeyPath != nil && *publicKeyPath != "" {
		publicLey, err = cryptor.GetPubKeyFromFile(*publicKeyPath)
		if err != nil {
			return nil, fmt.Errorf("cryptor.GetPubKeyFromFile: %w", err)
		}
	}
	return &HttpClient{
		httpClient: restyClient,
		hostPort:   hostPort,
		log:        log,
		keyHash:    keyHash,
		publicKey:  publicLey,
	}, nil
}

func (c *HttpClient) Send(ctx context.Context, metric models.Metrics) error {
	json := jsoniter.ConfigCompatibleWithStandardLibrary

	metricJSON, err := json.Marshal(metric)
	if err != nil {
		return err
	}
	// Шифроуем тело запроса, если передан публичный ключ
	if c.publicKey != nil {
		metricJSON, err = cryptor.EncryptBodyWithPublicKey(metricJSON, c.publicKey)
		if err != nil {
			return fmt.Errorf("cryptor.EncryptBodyWithPublicKey: %w", err)
		}
	}

	url := "/updates"

	client := c.httpClient.R().SetContext(ctx)

	// Добавляем заголовок с хешем тела запроса, если передан ключ
	if c.keyHash != nil {
		hash, err := cryptor.Encrypt(metricJSON, *c.keyHash)
		if err != nil {
			return fmt.Errorf("cryptor.Encrypt: %w", err)
		}

		client.SetHeader(models.HeaderHash, hash)
	}

	// Добавляем заголовок с локальным ip
	localIP, err := helpers.GetLocalIP()
	if err != nil {
		c.log.Error("helpers.GetLocalIP", zap.Error(err))
	} else {
		client.SetHeader(models.HeaderTrustedSubnet, localIP.String())
	}

	resp, err := client.SetBody(metricJSON).Post(c.hostPort + url)
	if err != nil {
		if ctx.Err() != nil {
			return nil
		}

		c.log.Error("resp", zap.Error(err))
		// или сеть или тело ответа
		return models.NewRetryError(err)
	}

	if err := checkStatus(resp.StatusCode(), string(resp.Body())); err != nil {
		c.log.Error("checkStatus", zap.Error(err))
		return err
	}

	return nil
}

func (c *HttpClient) Stop() error {
	return nil
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
