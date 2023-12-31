package handlers

import (
	"github.com/Genry72/collecting-metrics/internal/usecases/server"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandler_RunServer(t *testing.T) {
	type fields struct {
		useCases *server.Server
	}
	type args struct {
		port string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		//{
		//	name: "negative",
		//	fields: fields{
		//		useCases: nil,
		//	},
		//	args: args{
		//		port: "ab",
		//	},
		//	wantErr: true,
		//},
		{
			name: "positive",
			fields: fields{
				useCases: nil,
			},
			args: args{
				port: "8080",
			},
			wantErr: false,
		},
	}
	key := "superKey"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				useCases: tt.fields.useCases,
			}

			if tt.wantErr {
				err := h.RunServer(&tt.args.port, &key, nil, "")
				require.Error(t, err)
			}

		})
	}

}
