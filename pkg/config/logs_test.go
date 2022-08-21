package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/matrixbotio/go-common-lib/pkg/config"
)

func TestGetLogsConfig(t *testing.T) {
	tests := []struct {
		name        string
		envs        map[string]string
		expected    config.Logs
		wantedError error
	}{
		{
			name: "Not set",
			expected: config.Logs{
				LogToES:           true,
				AppLevel:          zapcore.InfoLevel,
				RMQWorkerLibLevel: zapcore.InfoLevel,
			},
		},
		{
			name: "LogToES=true",
			envs: map[string]string{
				"LOG_TO_ES": "true",
			},
			expected: config.Logs{
				LogToES:           true,
				AppLevel:          zapcore.InfoLevel,
				RMQWorkerLibLevel: zapcore.InfoLevel,
			},
		},
		{
			name: "LogToES=false",
			envs: map[string]string{
				"LOG_TO_ES": "false",
			},
			expected: config.Logs{
				LogToES:           false,
				AppLevel:          zapcore.InfoLevel,
				RMQWorkerLibLevel: zapcore.InfoLevel,
			},
		},
		{
			name: "LogToES=f",
			envs: map[string]string{
				"LOG_TO_ES": "f",
			},
			expected: config.Logs{
				LogToES:           false,
				AppLevel:          zapcore.InfoLevel,
				RMQWorkerLibLevel: zapcore.InfoLevel,
			},
		},
		{
			name: "LogToES=0",
			envs: map[string]string{
				"LOG_TO_ES": "0",
			},
			expected: config.Logs{
				LogToES:           false,
				AppLevel:          zapcore.InfoLevel,
				RMQWorkerLibLevel: zapcore.InfoLevel,
			},
		},
		{
			name: "LogToES=t",
			envs: map[string]string{
				"LOG_TO_ES": "t",
			},
			expected: config.Logs{
				LogToES:           true,
				AppLevel:          zapcore.InfoLevel,
				RMQWorkerLibLevel: zapcore.InfoLevel,
			},
		},
		{
			name: "LogToES=1",
			envs: map[string]string{
				"LOG_TO_ES": "1",
			},
			expected: config.Logs{
				LogToES:           true,
				AppLevel:          zapcore.InfoLevel,
				RMQWorkerLibLevel: zapcore.InfoLevel,
			},
		},
		{
			name: "AppLevel=log;RMQWorkerLibLevel=log",
			envs: map[string]string{
				"LOG_LEVEL":              "log",
				"LOG_LEVEL_RMQWORKERLIB": "log",
			},
			expected: config.Logs{
				LogToES:           true,
				AppLevel:          zapcore.InfoLevel,
				RMQWorkerLibLevel: zapcore.InfoLevel,
			},
		},
		{
			name: "AppLevel=verbose;RMQWorkerLibLevel=verbose",
			envs: map[string]string{
				"LOG_LEVEL":              "verbose",
				"LOG_LEVEL_RMQWORKERLIB": "verbose",
			},
			expected: config.Logs{
				LogToES:           true,
				AppLevel:          zapcore.DebugLevel,
				RMQWorkerLibLevel: zapcore.DebugLevel,
			},
		},
		{
			name: "AppLevel=xxx",
			envs: map[string]string{
				"LOG_LEVEL": "xxx",
			},
			wantedError: config.LogLevelIncorrectError{
				Level: "xxx",
				Field: "AppLevel",
			},
		},
		{
			name: "RMQWorkerLibLevel=xxx",
			envs: map[string]string{
				"LOG_LEVEL_RMQWORKERLIB": "xxx",
			},
			wantedError: config.LogLevelIncorrectError{
				Level: "xxx",
				Field: "RMQWorkerLibLevel",
			},
		},
	}

	for _, tc := range tests {
		os.Clearenv()

		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				require.NoError(t, os.Setenv(k, v))
			}

			c, err := config.GetLogsConfig()

			if tc.wantedError != nil {
				assert.ErrorIs(t, err, tc.wantedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, c)
			}
		})
	}
}
