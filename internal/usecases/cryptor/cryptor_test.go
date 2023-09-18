package cryptor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	key                = "superKey"
	cryptedString      = "b8ed6282843b48484b146da1e7ad17f016305ad4c7c596ad7b077c5a"
	cryptedEmptyString = "6af590b28f547e58a7a4b28065c9ddd6"
)

func TestEncrypt(t *testing.T) {
	type args struct {
		value    []byte
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				value:    []byte{},
				password: key,
			},
			want:    cryptedEmptyString,
			wantErr: false,
		},
		{
			name: "#2",
			args: args{
				value:    []byte("Привет"),
				password: key,
			},
			want:    cryptedString,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.value, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		value    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				value:    cryptedEmptyString,
				password: key,
			},
			want:    []byte(nil),
			wantErr: false,
		},
		{
			name: "#2",
			args: args{
				value:    cryptedString,
				password: key,
			},
			want:    []byte("Привет"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.value, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Decrypt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
