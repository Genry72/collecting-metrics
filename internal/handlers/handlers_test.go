package handlers

import (
	"github.com/Genry72/collecting-metrics/internal/repositories"
	"github.com/Genry72/collecting-metrics/internal/usecases"
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
		useCases *usecases.ServerUc
	}

	type args struct {
		method string
		url    string
	}

	repo := repositories.NewMemStorage()

	uc := usecases.NewServerUc(repo)

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
				code:     http.StatusMethodNotAllowed,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Handler{
				useCases: tt.fields.useCases,
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			h.setMetrics(w, r)
			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)
			if err := res.Body.Close(); err != nil {
				t.Error(err)
			}
		})
	}

}

func TestNewServer(t *testing.T) {
	type args struct {
		uc *usecases.ServerUc
	}
	tests := []struct {
		name string
		args args
		want *Handler
	}{
		{
			name: "positive",
			args: args{
				uc: &usecases.ServerUc{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.IsTypef(t, tt.want, NewServer(tt.args.uc), "NewServer(%v)", tt.args.uc)
		})
	}
}
