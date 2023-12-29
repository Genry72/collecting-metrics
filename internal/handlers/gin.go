package handlers

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/handlers/midlware/access"
	"github.com/Genry72/collecting-metrics/internal/handlers/midlware/cryptor"
	"github.com/Genry72/collecting-metrics/internal/handlers/midlware/gzip"
	"github.com/Genry72/collecting-metrics/internal/handlers/midlware/log"
	cryptor2 "github.com/Genry72/collecting-metrics/internal/usecases/cryptor"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

// Имена файлов при создании ключей
const (
	certFileName = "sertificate.cert"
	pubFileName  = "public.key"
)

// RunServer Функция использует фреймворк Gin для обработки HTTP-запросов и логирования.
// Входные параметры:
// - hostPort string: хост и порт, на котором будет запущен сервер.
// - password *string: пароль для доступа к серверу (необязательный параметр).
// - organization: Имя организации для получения сертификата (в случае его отсутствия по указанному пути).
// - privateKeyPath: Путь до закрытого ключа
// Возвращаемое значение:
// - error: ошибка, возникающая при запуске сервера.
func (h *Handler) RunServer(
	hostPort *string,
	password *string,
	privateKeyPath *string,
	organization string,
	trustedSubnet *string,
) error {
	if hostPort == nil || *hostPort == "" {
		return fmt.Errorf("empry address")
	}

	gin.SetMode(gin.ReleaseMode)

	g := gin.New()

	if trustedSubnet != nil {
		g.Use(access.CheckIPAddress(h.log, *trustedSubnet))
	}

	g.Use(log.ResponseLogger(h.log))
	g.Use(log.RequestLogger(h.log))

	g.Use(gzip.Gzip(h.log))
	// Проверка хеша тела запроса, если передан пароль
	if password != nil && *password != "" {
		g.Use(cryptor.CheckHashFromHeader(h.log, *password))
	}

	// При передаче ключа, подключаем обработчик по расшифровке тела запроса приватным ключем
	if privateKeyPath != nil && *privateKeyPath != "" {
		// Проверка наличия ключа по указанному пути, если его нет, то генерируем новый набор
		privKey, err := cryptor2.GetPrivateKeyFromFile(*privateKeyPath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				h.log.Info("The path to the private key provided by the file was not found. Generating a new keychain.")

				privKey, err = h.generateCert(*privateKeyPath, organization)
				if err != nil {
					return fmt.Errorf("h.generateCert: %w", err)
				}
			} else {
				return fmt.Errorf("cert.GetPrivateKeyFromFile: %w", err)
			}
		}

		g.Use(cryptor.DecryptBodyWithPrivateKey(h.log, privKey))
	}

	pprof.Register(g)

	h.setupRoute(g)

	if err := g.Run(*hostPort); err != nil {
		return fmt.Errorf("g.Run: %w", err)
	}

	return nil
}

// setupRoute Установка хендлеров
func (h *Handler) setupRoute(g *gin.Engine) {
	g.GET("/", h.getAllMetrics)
	g.GET("/ping", h.pingDatabase)

	update := g.Group("update")
	update.POST("/", h.setMetricJSON)
	update.POST("/:type/:name/:value", h.setMetricsText)

	updates := g.Group("updates")
	updates.POST("/", h.setMetricsJSON)

	value := g.Group("value")
	value.POST("/", h.getMetricsJSON)
	value.GET("/:type/:name", h.getMetricText)
}

// generateCert Получение и сохранение информации по связке ключей в файлы. Возвращает созданный приватный ключ
func (h *Handler) generateCert(privateKeyPath, organization string) (*rsa.PrivateKey, error) {
	cert, pub, private, err := cryptor2.GetCertsAndKeys(organization)
	if err != nil {
		return nil, fmt.Errorf("cert.GetCerts: %w", err)
	}

	if err := os.WriteFile(privateKeyPath, private, 0644); err != nil {
		return nil, fmt.Errorf("WriteFile.private: %w", err)
	}

	certPath := filepath.Join(filepath.Dir(privateKeyPath), certFileName)
	if err := os.WriteFile(certPath, cert, 0644); err != nil {
		return nil, fmt.Errorf("WriteFile.public: %w", err)
	}

	publicPath := filepath.Join(filepath.Dir(privateKeyPath), pubFileName)
	if err := os.WriteFile(publicPath, pub, 0644); err != nil {
		return nil, fmt.Errorf("WriteFile.public: %w", err)
	}

	return cryptor2.GetPrivateKeyFromFile(privateKeyPath)
}
