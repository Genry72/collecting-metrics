package flags_test

import (
	"flag"
	"github.com/Genry72/collecting-metrics/cmd/agent/flags"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_parseFlag(t *testing.T) {
	port999 := ":9999"
	port8081 := ":8081"
	confPath := "./test_config.json"
	defReportInterval := 10
	defRateLimit := 1
	defPooltInterval := 2
	confReportInterval := 1
	confPooltInterval := 1
	confCryptoKey := "/path/to/key.pem"

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
				Address:        &port999,
				ReportInterval: &defReportInterval,
				PollInterval:   &defPooltInterval,
				CryptoKey:      nil,
				RateLimit:      &defRateLimit,
			},
			args: args{func() {
				os.Args = append(os.Args, `-a=:9999`)
			}},
		},
		{
			name: "Переменная окружения",
			want: &flags.RunParameters{
				Address:        &port999,
				ReportInterval: &defReportInterval,
				PollInterval:   &defPooltInterval,
				CryptoKey:      nil,
				RateLimit:      &defRateLimit,
			},
			args: args{func() {
				os.Setenv("ADDRESS", port999)
			}},
		},
		{
			name: "Переменная окружения и флаг",
			want: &flags.RunParameters{
				Address:        &port999,
				ReportInterval: &defReportInterval,
				PollInterval:   &defPooltInterval,
				CryptoKey:      nil,
				RateLimit:      &defRateLimit,
			},
			args: args{func() {
				os.Setenv("ADDRESS", port999)
				os.Args = append(os.Args, `-a=:9998`)
			}},
		},
		{
			name: "Конфиг",
			want: &flags.RunParameters{
				Address:        &port8081,
				ReportInterval: &confReportInterval,
				PollInterval:   &confPooltInterval,
				CryptoKey:      &confCryptoKey,
				RateLimit:      &defRateLimit,
			},
			args: args{func() {
				os.Setenv("CONFIG", confPath)
			}},
		},
		{
			name: "Конфиг через флаг",
			want: &flags.RunParameters{
				Address:        &port8081,
				ReportInterval: &confReportInterval,
				PollInterval:   &confPooltInterval,
				CryptoKey:      &confCryptoKey,
				RateLimit:      &defRateLimit,
			},
			args: args{func() {
				os.Args = append(os.Args, `-config=`+confPath)
			}},
		},
		{
			name: "Конфиг и флаг",
			want: &flags.RunParameters{
				Address:        &port999,
				ReportInterval: &confReportInterval,
				PollInterval:   &confPooltInterval,
				CryptoKey:      &confCryptoKey,
				RateLimit:      &defRateLimit,
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
