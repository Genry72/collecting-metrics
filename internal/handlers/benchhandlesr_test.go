package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/logger"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/repositories/memstorage"
	"github.com/Genry72/collecting-metrics/internal/usecases/server"
	"github.com/gin-gonic/gin"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"
)

func BenchmarkHandlers(b *testing.B) {
	runtime.MemProfileRate = 0
	zapLogger := logger.NewZapLogger("info")
	repo := memstorage.NewMemStorage(zapLogger)

	gin.SetMode(gin.ReleaseMode)
	ginRoute := gin.New()

	h := Handler{
		useCases: server.NewServerUc(repo, nil, nil, zapLogger),
		log:      zapLogger,
	}

	h.setupRoute(ginRoute, nil) // todo добавить пароль

	type args struct {
		metricCount int // Количество метрик, необходимых для теста
		method      string
		// Функция для генерации URL запроса, необходима для получения случайного значения типа и значения метрики
		generateUrl func([]*models.Metric) string
		// Генерация тела запроса
		body func([]*models.Metric) io.Reader
		// true говорит о том что это проверка загруженных метрик, для этого данные метрики сначала нужно загрузить
		getMetric bool
	}

	tests := []args{
		{ // Получение всех метрик
			method: http.MethodGet,
			generateUrl: func(_ []*models.Metric) string {
				return "/"
			},
			metricCount: 1,
		},
		{ // Добавление списка метрик
			method: http.MethodPost,
			generateUrl: func(_ []*models.Metric) string {
				return "/updates/"
			},
			body: func(metrics []*models.Metric) io.Reader {
				bByte, err := json.Marshal(metrics)
				if err != nil {
					err = fmt.Errorf("marshal body: %w", err)
					fmt.Println(err, metrics)
				}
				return bytes.NewReader(bByte)
			},
			metricCount: 10,
		},
		{ // Добавление одной метрики json
			method: http.MethodPost,
			generateUrl: func(_ []*models.Metric) string {
				return "/update/"
			},
			body: func(metrics []*models.Metric) io.Reader {
				bByte, err := json.Marshal(metrics[0])
				if err != nil {
					err = fmt.Errorf("marshal body: %w", err)
					fmt.Println(err, metrics)
				}
				return bytes.NewReader(bByte)
			},
			metricCount: 1,
		},
		{ // Получение метрики
			method: http.MethodPost,
			generateUrl: func(_ []*models.Metric) string {
				return "/value/"
			},
			body: func(metrics []*models.Metric) io.Reader {
				bByte, err := json.Marshal(metrics[0])
				if err != nil {
					err = fmt.Errorf("marshal body: %w", err)
					fmt.Println(err, metrics)
				}
				return bytes.NewReader(bByte)
			},
			metricCount: 1,
			getMetric:   true,
		},
		{ // Получение метрики
			method: http.MethodGet,
			generateUrl: func(metrics []*models.Metric) string {
				return fmt.Sprintf("/value/%s/%s", metrics[0].MType, metrics[0].ID)
			},
			metricCount: 1,
			getMetric:   true,
		},
	}

	runtime.MemProfileRate = 1
	b.ResetTimer()

	var m []*models.Metric // Хранение загруженных метрик в последней итерации тестов. Для тестов получения метрик

	b.Run("new", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, tt := range tests {
				runtime.MemProfileRate = 0
				// рандомные метрики
				metrics := generateMetric(tt.metricCount, rand.New(rand.NewSource(time.Now().UnixNano())))
				runtime.MemProfileRate = 1
				if tt.getMetric {
					runBanchHandlers(ginRoute, m, tt.method, tt.generateUrl, tt.body)
					continue
				}
				runBanchHandlers(ginRoute, metrics, tt.method, tt.generateUrl, tt.body)
				m = metrics
			}
		}
	})

}

func runBanchHandlers(
	ginRoute *gin.Engine,
	metrics []*models.Metric,
	method string,
	generateUrl func([]*models.Metric) string,
	bodyFunc func(metrics []*models.Metric) io.Reader) {
	runtime.MemProfileRate = 0

	var body io.Reader

	if bodyFunc != nil {
		body = bodyFunc(metrics)
	}

	w := httptest.NewRecorder()

	r := httptest.NewRequest(
		method,
		generateUrl(metrics),
		body)

	runtime.MemProfileRate = 1

	ginRoute.ServeHTTP(w, r)
}

func generateMetric(count int, r *rand.Rand) []*models.Metric {
	metrics := make([]*models.Metric, count)
	// список типов метрик
	metricTypes := []models.MetricType{models.MetricTypeGauge, models.MetricTypeCounter}

	gaugeVal := r.Float64()
	counterValue := r.Int63()
	metricName := generateString(r)
	// мапа для получения случайной метрики по ее типу
	metricMap := map[models.MetricType]func() *models.Metric{
		models.MetricTypeGauge: func() *models.Metric {
			return &models.Metric{
				ID:        models.MetricName(metricName),
				MType:     models.MetricTypeGauge,
				Delta:     nil,
				Value:     &gaugeVal,
				ValueText: fmt.Sprint(gaugeVal),
			}
		},
		models.MetricTypeCounter: func() *models.Metric {
			return &models.Metric{
				ID:        models.MetricName(metricName),
				MType:     models.MetricTypeCounter,
				Delta:     &counterValue,
				Value:     nil,
				ValueText: fmt.Sprint(counterValue),
			}
		},
	}

	for i := 0; i < count; i++ {
		// случайный тип метрики
		mType := metricTypes[r.Intn(len(metricTypes))]
		// метрика со случайным значение по ее типу
		metrics[i] = metricMap[mType]()
	}

	return metrics
}

func generateString(r *rand.Rand) string {
	// Формируем случайное имя метрики и значение типа Float64 преобразованное в string
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 10

	randomMetricName := make([]byte, length)
	for i := range randomMetricName {
		randomMetricName[i] = charset[r.Intn(len(charset))]
	}
	return string(randomMetricName)
}
