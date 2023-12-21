package flags

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
)

/*
RunParameters параметры запуска сервера.
Приоритеты:
1. Переменная окружения
2. Флаг
3. Файл конфигурации
4. Дефолтное значение (указывается для флага)
*/

const (
	// Ключ для шифрования.
	envKeyHash  = "KEY"
	envConfPath = "CONFIG"
)

type RunParameters struct {
	// Адрес и порт для запуска сервера
	Address *string `json:"address" env:"ADDRESS" flag:"a" default:":8080" comment:"Адрес и порт для запуска сервера"`
	// булево значение (true/false), определяющее, загружать или нет ранее сохранённые значения из указанного файла
	// при старте сервера (по умолчанию true)
	Restore *bool `json:"restore" env:"RESTORE" flag:"r" default:"true" comment:"Загрузка из файла при старте"`
	// интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск
	// (по умолчанию 300 секунд, значение 0 делает запись синхронной)
	StoreInterval *int `json:"store_interval" env:"STORE_INTERVAL" flag:"i" default:"300" comment:"Периодичность сохранения на диск"`
	// полное имя файла, куда сохраняются текущие значения (по умолчанию /tmp/metrics-db.json, пустое значение
	// отключает функцию записи на диск)
	StoreFile *string `json:"store_file" env:"FILE_STORAGE_PATH" flag:"f" default:"/tmp/metrics-db.json" comment:"Путь до файля для сохрания метрик"`
	// Подключение к базе данных
	DatabaseDsn *string `json:"database_dsn" env:"DATABASE_DSN" flag:"d" default:"" comment:"Подключение к базе данных"`
	// Путь до файла с приватным ключом
	CryptoKey *string `json:"crypto_key" env:"CRYPTO_KEY" flag:"crypto-key" default:"" comment:"Путь до файла с приватным ключом"`
	// Кдюч шифрования
	KeyHash *string `json:"-"`
}

func ParseFlag() (*RunParameters, error) {
	params := RunParameters{}

	t := reflect.TypeOf(params)

	v := reflect.ValueOf(&params)

	funcs := make([]func(configValue reflect.Value) (reflect.Value, error), 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := reflect.New(field.Type.Elem()).Elem()
		flagName := field.Tag.Get("flag")
		flagComment := field.Tag.Get("comment")
		flagDefault := field.Tag.Get("default")
		tagEnv := field.Tag.Get("env")

		var f func(configValue reflect.Value) (reflect.Value, error)

		switch fieldValue.Kind() {
		case reflect.String:
			var flagValue string

			flag.StringVar(&flagValue, flagName, flagDefault, flagComment)

			f = func(configValue reflect.Value) (reflect.Value, error) {
				// Переменная окружения
				if envVal, ok := os.LookupEnv(tagEnv); ok {
					return reflect.ValueOf(&envVal), nil
				}

				// Значение флага
				if isFlagPassed(flagName) {
					return reflect.ValueOf(&flagValue), nil
				}

				// Значение из файла конфигурации
				if !configValue.IsZero() && configValue.Elem().String() != "" {
					return configValue, nil
				}

				// Значение по умолчанию
				if flagDefault == "" {
					var c *string
					return reflect.ValueOf(c), nil
				}

				return reflect.ValueOf(&flagDefault), nil
			}

		case reflect.Int:
			var flagValue int

			def, err := strconv.ParseInt(flagDefault, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("strconv.ParseInt: %w", err)
			}

			flag.IntVar(&flagValue, flagName, int(def), flagComment)

			f = func(configValue reflect.Value) (reflect.Value, error) {
				// Переменная окружения
				if envVal, ok := os.LookupEnv(tagEnv); ok {
					def, err := strconv.ParseInt(envVal, 10, 64)
					if err != nil {
						return reflect.Value{}, fmt.Errorf("strconv.ParseInt: %w", err)
					}

					z := int(def)

					return reflect.ValueOf(&z), nil
				}

				// Значение флага
				if isFlagPassed(flagName) {
					return reflect.ValueOf(&flagValue), nil
				}

				var c *int
				// Значение из файла конфигурации
				if !configValue.IsZero() && configValue.Elem().String() != "" {
					return configValue, nil
				}

				// Значение по умолчанию
				if flagDefault == "" {
					return reflect.ValueOf(c), nil
				}

				z := int(def)
				return reflect.ValueOf(&z), nil

			}

		case reflect.Bool:
			var flagValue bool

			def, err := strconv.ParseBool(flagDefault)
			if err != nil {
				return nil, fmt.Errorf("strconv.ParseInt: %w", err)
			}

			flag.BoolVar(&flagValue, flagName, def, flagComment)

			f = func(configValue reflect.Value) (reflect.Value, error) {
				df := reflect.TypeOf(configValue).Kind()
				_ = df
				// Переменная окружения
				if envVal, ok := os.LookupEnv(tagEnv); ok {
					def, err := strconv.ParseBool(envVal)
					if err != nil {
						return reflect.Value{}, fmt.Errorf("env.ParseBool: %w", err)
					}

					return reflect.ValueOf(&def), nil
				}

				// Значение флага
				if isFlagPassed(flagName) {
					return reflect.ValueOf(&flagValue), nil
				}

				// Значение из файла конфигурации
				if !configValue.IsZero() && configValue.Elem().String() != "" {
					return configValue, nil
				}

				// Значение по умолчанию
				if flagDefault == "" {
					var c *bool
					return reflect.ValueOf(c), nil
				}

				return reflect.ValueOf(&def), nil
			}

		default:
			return nil, fmt.Errorf("unknown type: %s", fieldValue.Kind())
		}

		funcs = append(funcs, f)
	}

	var (
		flagConfig  string
		flagKeyHash string
	)

	flag.StringVar(&flagConfig, "с", "", "Путь до файла конфигурации")
	flag.StringVar(&flagConfig, "config", "", "Путь до файла конфигурации")
	flag.StringVar(&flagKeyHash, "k", "", "Кдюч шифрования")

	flag.Parse()

	if env := os.Getenv(envConfPath); env != "" {
		flagConfig = env
	}

	if env := os.Getenv(envKeyHash); env != "" {
		flagKeyHash = env
	}

	var err error

	config := RunParameters{}

	if flagConfig != "" {
		config, err = parseConf(flagConfig)
		if err != nil {
			return nil, fmt.Errorf("parseConf: %w", err)
		}
	}

	vconf := reflect.ValueOf(config)

	for i := range funcs {
		val, err := funcs[i](vconf.Field(i))
		if err != nil {
			return nil, fmt.Errorf("funcs[i]: %w", err)
		}

		v.Elem().Field(i).Set(val)
	}

	return &params, nil
}

func parseConf(flagConfig string) (RunParameters, error) {
	conf := RunParameters{}
	if flagConfig == "" {
		return conf, nil
	}

	f, err := os.ReadFile(flagConfig)

	if err != nil {
		return conf, fmt.Errorf("os.ReadFile: %w", err)
	}

	if err := json.Unmarshal(f, &conf); err != nil {
		return conf, fmt.Errorf("json.Unmarshal: %w", err)
	}

	return conf, nil
}

func isFlagPassed(name string) bool {
	found := false

	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})

	return found
}
