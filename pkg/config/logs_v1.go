package config

import (
	"fmt"
	"reflect"

	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap/zapcore"
)

type Logs struct {
	LogToES            bool
	SamplingInitial    int
	SamplingThereafter int

	AppLevel          zapcore.Level
	RMQWorkerLibLevel zapcore.Level
}

// Logs config v1
func GetLogsConfig() (Logs, error) {
	var Parsed struct {
		LogToES            bool `envconfig:"LOG_TO_ES" default:"true"`
		SamplingInitial    int  `envconfig:"LOG_SAMPLING_INITIAL" default:"300"`
		SamplingThereafter int  `envconfig:"LOG_SAMPLING_THEREAFTER" default:"100"`

		AppLevel          string `envconfig:"LOG_LEVEL" default:"log"`
		RMQWorkerLibLevel string `envconfig:"LOG_LEVEL_RMQWORKERLIB" default:"log"`
	}

	if err := envconfig.Process("", &Parsed); err != nil {
		return Logs{}, fmt.Errorf("parse envs: %w", err)
	}

	cfg := Logs{
		LogToES:            Parsed.LogToES,
		SamplingInitial:    Parsed.SamplingInitial,
		SamplingThereafter: Parsed.SamplingThereafter,
	}

	for _, field := range []string{"AppLevel", "RMQWorkerLibLevel"} {
		val := reflect.ValueOf(Parsed).FieldByName(field).String()

		var zapLevel zapcore.Level
		switch val {
		case "log":
			zapLevel = zapcore.InfoLevel
		case "verbose":
			zapLevel = zapcore.DebugLevel
		case "warn":
			zapLevel = zapcore.WarnLevel
		case "error":
			zapLevel = zapcore.ErrorLevel
		case "critical":
			zapLevel = zapcore.FatalLevel
		default:
			return Logs{}, LogLevelIncorrectError{
				Level: val,
				Field: field,
			}
		}
		reflect.ValueOf(&cfg).Elem().FieldByName(field).Set(reflect.ValueOf(zapLevel))
	}

	return cfg, nil
}
