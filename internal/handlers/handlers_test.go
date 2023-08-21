package handlers

import (
	"github.com/Genry72/collecting-metrics/internal/logger"
	"github.com/Genry72/collecting-metrics/internal/repositories/memstorage"
	"github.com/Genry72/collecting-metrics/internal/usecases/server"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_setMetrics(t *testing.T) {
	type want struct {
		code     int
		response string
	}

	type fields struct {
		useCases *server.Server
	}

	type args struct {
		method string
		url    string
	}
	zapLogger := logger.NewZapLogger("info")

	repo := memstorage.NewMemStorage(zapLogger)

	uc := server.NewServerUc(repo, zapLogger)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "positive test #1",
			fields: fields{
				useCases: uc,
			},
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
			fields: fields{
				useCases: uc,
			},
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
			fields: fields{
				useCases: uc,
			},
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
			fields: fields{
				useCases: uc,
			},
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
			fields: fields{
				useCases: uc,
			},
			args: args{
				method: http.MethodPost,
				url:    "/update/gauge/11",
			},
			want: want{
				code:     http.StatusNotFound,
				response: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{
				useCases: tt.fields.useCases,
				log:      zapLogger,
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			//h.setMetricsText(w, r)
			g := gin.Default()
			h.setupRoute(g)
			g.ServeHTTP(w, r)
			//res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, w.Code, w.Body.String())
			//if err := w.Body.; err != nil {
			//	t.Error(err)
			//}
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
