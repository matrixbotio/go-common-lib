package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matrixbotio/go-common-lib/pkg/config"
)

func TestGetRMQConfig(t *testing.T) {
	tests := []struct {
		name        string
		envs        map[string]string
		expected    config.RMQ
		wantedError error
	}{
		{
			name: "User not set",
			envs: map[string]string{
				"AMQP_PASSWORD": "psw2",
			},
			wantedError: config.EmptyFieldError{"user"},
		},
		{
			name: "Password not set",
			envs: map[string]string{
				"AMQP_USER": "usr1",
			},
			wantedError: config.EmptyFieldError{"password"},
		},
		{
			name: "Set",
			envs: map[string]string{
				"AMQP_USE_TLS":  "t",
				"AMQP_HOST":     "hst",
				"AMQP_PORT":     "1234",
				"AMQP_USER":     "usr1",
				"AMQP_PASSWORD": "psw2",
			},
			expected: config.RMQ{
				UseTLS:   true,
				Host:     "hst",
				Port:     1234,
				User:     "usr1",
				Password: "psw2",
			},
		},
		{
			name: "Set with some defaults",
			envs: map[string]string{
				"AMQP_USER":     "usr1",
				"AMQP_PASSWORD": "psw2",
			},
			expected: config.RMQ{
				UseTLS:   false,
				Host:     "localhost",
				Port:     5672,
				User:     "usr1",
				Password: "psw2",
			},
		},
	}

	for _, tc := range tests {
		os.Clearenv()

		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				require.NoError(t, os.Setenv(k, v))
			}

			c, err := config.GetRMQConfig()

			if tc.wantedError != nil {
				assert.ErrorIs(t, err, tc.wantedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, c)
			}
		})
	}
}
