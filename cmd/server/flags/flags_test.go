package flags_test

import (
	"flag"
	"github.com/Genry72/collecting-metrics/cmd/server/flags"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_parseFlag(t *testing.T) {
	port999 := ":9999"
	port8081 := ":8081"
	defPort8081 := ":8080"
	confPath := "./test_config.json"
	defTrue := true
	False := false
	defStoreInterval := 300
	confStoreInterval := 1
	defStoreFile := "/tmp/metrics-db.json"
	subnet := "127.0.0.1"

	type args struct {
		f func()
	}
	tests := []struct {
		name string
		want *flags.RunParameters
		args args
	}{
		{
			name: "Флаг",
			want: &flags.RunParameters{
				Address:       &port999,
				Restore:       &defTrue, // default
				StoreInterval: &defStoreInterval,
				StoreFile:     &defStoreFile,
				DatabaseDsn:   nil,
				CryptoKey:     nil,
			},
			args: args{func() {
				os.Args = append(os.Args, `-a=:9999`)
			}},
		},
		{
			name: "Флаг TrustedSubnet",
			want: &flags.RunParameters{
				Address:       &defPort8081,
				Restore:       &defTrue, // default
				StoreInterval: &defStoreInterval,
				StoreFile:     &defStoreFile,
				DatabaseDsn:   nil,
				CryptoKey:     nil,
				TrustedSubnet: &subnet,
			},
			args: args{func() {
				os.Args = append(os.Args, "-t="+subnet)
			}},
		},
		{
			name: "Переменная окружения",
			want: &flags.RunParameters{
				Address:       &port999,
				Restore:       &defTrue,          // default
				StoreInterval: &defStoreInterval, // default
				StoreFile:     &defStoreFile,     // default
				DatabaseDsn:   nil,
				CryptoKey:     nil,
			},
			args: args{func() {
				os.Setenv("ADDRESS", port999)
			}},
		},
		{
			name: "Переменная окружения и флаг",
			want: &flags.RunParameters{
				Address:       &port999,
				Restore:       &defTrue,          // default
				StoreInterval: &defStoreInterval, // default
				StoreFile:     &defStoreFile,     // default
				DatabaseDsn:   nil,
				CryptoKey:     nil,
			},
			args: args{func() {
				os.Setenv("ADDRESS", port999)
				os.Args = append(os.Args, `-a=:9998`)
			}},
		},
		{
			name: "Конфиг",
			want: &flags.RunParameters{
				Address:       &port8081,
				Restore:       &False,             // default
				StoreInterval: &confStoreInterval, // default
				StoreFile:     &defStoreFile,      // default
				DatabaseDsn:   nil,
				CryptoKey:     nil,
			},
			args: args{func() {
				os.Setenv("CONFIG", confPath)
			}},
		},
		{
			name: "Конфиг через флаг",
			want: &flags.RunParameters{
				Address:       &port8081,
				Restore:       &False,             // default
				StoreInterval: &confStoreInterval, // conf
				StoreFile:     &defStoreFile,      // default
				DatabaseDsn:   nil,
				CryptoKey:     nil,
			},
			args: args{func() {
				os.Args = append(os.Args, `-config=`+confPath)
			}},
		},
		{
			name: "Конфиг и флаг",
			want: &flags.RunParameters{
				Address:       &port999,
				Restore:       &False,             // default
				StoreInterval: &confStoreInterval, // conf
				StoreFile:     &defStoreFile,      // default
				DatabaseDsn:   nil,
				CryptoKey:     nil,
			},
			args: args{func() {
				os.Args = append(os.Args, `-config=`+confPath)
				os.Args = append(os.Args, `-a=:9999`)
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = os.Args[:1]
			tt.args.f()
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			res, err := flags.ParseFlag()
			assert.NoError(t, err)

			assert.Equal(t, tt.want, res)
			os.Clearenv()

		})

	}
}
