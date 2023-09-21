package agent

import (
	"github.com/Genry72/collecting-metrics/internal/logger"
	"github.com/Genry72/collecting-metrics/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAgent_send(t *testing.T) {
	type args struct {
		metric       *models.Metric
		url          string
		responseCode int
	}
	value := float64(5)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive#1",
			args: args{
				url:          "/update/gauge/testmetric/5",
				responseCode: http.StatusOK,
				metric: &models.Metric{
					ID:        "testmetric",
					MType:     "gauge",
					Delta:     nil,
					Value:     &value,
					ValueText: "",
				},
			},

			wantErr: false,
		},
		{
			name: "negative#1",
			args: args{
				url:          "/some/path",
				responseCode: http.StatusBadRequest,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				// Test request parameters

				// Send response to be tested
				rw.WriteHeader(tt.args.responseCode)
				if _, err := rw.Write([]byte(`OK`)); err != nil {
					t.Error(err)
				}

			}))

			defer server.Close()

			a := &Agent{
				httpClient:    resty.New(),
				hostPort:      server.URL,
				log:           logger.NewZapLogger("info"),
				ratelimitChan: make(chan struct{}, 1),
			}

			if err := a.sendByJSONBatch([]*models.Metric{tt.args.metric}); (err != nil) != tt.wantErr {
				t.Errorf("sendByURL() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestNewAgent(t *testing.T) {
	type args struct {
		hostPort string
	}
	zapLogger := logger.NewZapLogger("info")
	tests := []struct {
		name string
		args args
		want *Agent
	}{
		{
			name: "",
			args: args{
				hostPort: "",
			},
			want: nil,
		},
	}
	key := "superKey"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAgent(tt.args.hostPort, zapLogger, &key, 1); !reflect.DeepEqual(got, tt.want) {
				require.IsType(t, &Agent{}, got)
			}
		})
	}
}
