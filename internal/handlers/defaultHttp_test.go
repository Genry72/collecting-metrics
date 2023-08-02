package handlers

import (
	"github.com/Genry72/collecting-metrics/internal/usecases"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandler_RunServer(t *testing.T) {
	type fields struct {
		useCases *usecases.ServerUc
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
		{
			name: "negative",
			fields: fields{
				useCases: nil,
			},
			args: args{
				port: "ab",
			},
			wantErr: true,
		},
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				useCases: tt.fields.useCases,
			}

			if tt.wantErr {
				err := h.RunServer(tt.args.port)
				require.Error(t, err)
			}

		})
	}

}
