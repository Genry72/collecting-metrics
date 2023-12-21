package handlers

import (
	"fmt"
	"github.com/Genry72/collecting-metrics/internal/handlers/midlware/access"
	"github.com/Genry72/collecting-metrics/internal/logger"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/Genry72/collecting-metrics/internal/repositories/filestorage"
	"github.com/Genry72/collecting-metrics/internal/repositories/memstorage"
	"github.com/Genry72/collecting-metrics/internal/repositories/postgre"
	"github.com/Genry72/collecting-metrics/internal/usecases/server"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_setMetrics(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	type args struct {
		method   string
		url      string
		headers  map[string]string
		midlware func(*gin.Engine)
	}

	zapLogger := logger.NewZapLogger("info")

	repo := memstorage.NewMemStorage(zapLogger)
	ps, err := filestorage.NewFileStorage(&filestorage.StorageConf{
		StoreInterval:   0,
		FileStorageFile: "./fs",
		Restore:         false,
	}, zapLogger)

	assert.NoError(t, err)

	dsn := "postgres://postgres:pass@localhost:5432/postgres?sslmode=disable"

	pg, _ := postgre.NewPGStorage(&dsn, zapLogger)

	uc := server.NewServerUc(repo, ps, pg, zapLogger)

	h := Handler{
		useCases: uc,
		log:      zapLogger,
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "positive test #1",
			args: args{
				method: http.MethodPost,
				url:    "/update/gauge/name/11",
			},
			want: want{
				code:     http.StatusOK,
				response: "",
			},
		},
		{
			name: "negative test #1",
			args: args{
				method: http.MethodGet,
				url:    "/update/gauge/name/11",
			},
			want: want{
				code:     http.StatusNotFound,
				response: "",
			},
		},
		{
			name: "negative test #2",
			args: args{
				method: http.MethodPost,
				url:    "/update/gauge/name",
			},
			want: want{
				code:     http.StatusNotFound,
				response: "",
			},
		},
		{
			name: "negative test #3",
			args: args{
				method: http.MethodPost,
				url:    "/update/gauge/name/aa",
			},
			want: want{
				code:     http.StatusBadRequest,
				response: "",
			},
		},
		{
			name: "negative test #4",
			args: args{
				method: http.MethodPost,
				url:    "/update/gauge/11",
			},
			want: want{
				code:     http.StatusNotFound,
				response: "",
			},
		},
		{
			name: "positive CheckIPAddress #1",
			args: args{
				method: http.MethodPost,
				url:    "/update/gauge/name/11",
				midlware: func(g *gin.Engine) {
					g.Use(access.CheckIPAddress(zapLogger, "192.168.0.1/24"))
				},
				headers: map[string]string{
					models.HeaderTrustedSubnet: "192.168.0.5",
				},
			},
			want: want{
				code:     http.StatusOK,
				response: "",
			},
		},
		{
			name: "negative CheckIPAddress #1",
			args: args{
				method: http.MethodPost,
				url:    "/update/gauge/name/11",
				midlware: func(g *gin.Engine) {
					g.Use(access.CheckIPAddress(zapLogger, "192.168.0.1/24"))
				},
				headers: map[string]string{
					models.HeaderTrustedSubnet: "192.169.0.5",
				},
			},
			want: want{
				code:     http.StatusForbidden,
				response: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.ReleaseMode)

			g := gin.New()
			if tt.args.midlware != nil {
				tt.args.midlware(g)
			}
			//g.Use(access.CheckIPAddress(zapLogger, "192.168.0.1/24"))
			h.setupRoute(g)

			w := httptest.NewRecorder()

			r := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			// Добавляем заголовки
			for k, v := range tt.args.headers {
				r.Header.Set(k, v)
			}
			g.ServeHTTP(w, r)
			// проверяем код ответа
			assert.Equal(t, tt.want.code, w.Code, w.Body.String())
		})
	}

}

func TestNewServer(t *testing.T) {
	type args struct {
		uc *server.Server
	}
	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "positive",
			args: args{
				uc: &server.Server{},
			},
			want: nil,
		},
	}

	zapLogger := logger.NewZapLogger("info")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.IsTypef(t, tt.want, NewServer(tt.args.uc, zapLogger), "NewServer(%v)", tt.args.uc)
		})
	}
}

// Получение всех метрик в формате html
func Example_getAllMetrics() {
	url := "http://localhost:8080"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

// Example_ping Проверка доступности базы данных
func Example_ping() {
	url := "http://localhost:8080/ping"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body))
}

// Example_ping Отправка значения метрики в формате JSON
func Example_setMetricJSON() {
	url := "http://localhost:8080/update"
	method := "POST"

	payload := strings.NewReader(`{
    "id": "name",
    "type":"counter",
    "delta":123
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

// Example_setMetricsText Отправка значения метрики в формате JSON в параметрах запроса
func Example_setMetricsText() {
	url := "http://localhost:8080/update/counter/name/123"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

// Example_setMetricsJSON Отправка значений метрик списком в формате JSON в параметрах запроса
func Example_setMetricsJSON() {
	url := "http://localhost:8080/updates"
	method := "POST"

	payload := strings.NewReader(`[
    {
        "id": "name",
        "type": "counter",
        "delta": 123
    },
        {
        "id": "name2",
        "type": "counter",
        "delta": 123
    }
]`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

// Example_getMetricsJSON Получение значения метрики POST запрос
func Example_getMetricsJSON() {
	url := "http://localhost:8080/value"
	method := "POST"

	payload := strings.NewReader(`{
    "id": "name",
    "type": "counter"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

// Example_getMetricText Получение значения метрики GET запрос
func Example_getMetricText() {
	url := "http://localhost:8080/value/counter/name"
	method := "GET"

	payload := strings.NewReader(`{
    "id": "name",
    "type": "counter"
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
