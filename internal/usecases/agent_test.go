package usecases

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestAgent_send(t *testing.T) {
	type args struct {
		url          string
		responseCode int
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "positive#1",
			args: args{
				url:          "/some/path",
				responseCode: http.StatusOK,
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
				require.Equal(t, req.URL.String(), tt.args.url)

				// Send response to be tested
				rw.WriteHeader(tt.args.responseCode)
				if _, err := rw.Write([]byte(`OK`)); err != nil {
					t.Error(err)
				}

			}))

			defer server.Close()

			a := &Agent{
				httpClient: &http.Client{},
				hostPort:   server.URL,
			}

			if err := a.send(tt.args.url); (err != nil) != tt.wantErr {
				t.Errorf("send() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}

func TestNewAgent(t *testing.T) {
	type args struct {
		hostPort string
	}
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAgent(tt.args.hostPort); !reflect.DeepEqual(got, tt.want) {
				require.IsType(t, &Agent{}, got)
			}
		})
	}
}
