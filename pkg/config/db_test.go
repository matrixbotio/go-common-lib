package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/matrixbotio/go-common-lib/pkg/config"
)

func TestGetDBConfig(t *testing.T) {
	tests := []struct {
		name            string
		envs            map[string]string
		defaultName     string
		defaultUser     string
		defaultPassword string
		expected        config.DB
		wantedError     error
	}{
		{
			name:            "Not set",
			defaultName:     "nm",
			defaultUser:     "usr",
			defaultPassword: "psw",
			expected: config.DB{
				Host:          "localhost",
				Port:          3306,
				Name:          "nm",
				User:          "usr",
				Password:      "psw",
				ConnTimeoutMS: 5000,
				MaxOpenConns:  10,
			},
		},
		{
			name: "Set",
			envs: map[string]string{
				"DB_HOST":           "hst1",
				"DB_PORT":           "567",
				"DB_NAME":           "nm2",
				"DB_USER":           "usr3",
				"DB_PASSWORD":       "psw4",
				"DB_CONN_TIMEOUT":   "678",
				"DB_MAX_OPEN_CONNS": "5",
			},
			expected: config.DB{
				Host:          "hst1",
				Port:          567,
				Name:          "nm2",
				User:          "usr3",
				Password:      "psw4",
				ConnTimeoutMS: 678,
				MaxOpenConns:  5,
			},
		},
		{
			name: "Host set to empty",
			envs: map[string]string{
				"DB_HOST": "",
			},
			wantedError: config.EmptyFieldError{
				Field: "host",
			},
		},
	}

	for _, tc := range tests {
		os.Clearenv()

		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.envs {
				require.NoError(t, os.Setenv(k, v))
			}

			c, err := config.GetDBConfig(tc.defaultName, tc.defaultUser, tc.defaultPassword)

			if tc.wantedError != nil {
				assert.ErrorIs(t, err, tc.wantedError)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, c)
			}
		})
	}
}

func TestDB_GetDSN(t *testing.T) {
	c := config.DB{
		Host:          "hst",
		Port:          45,
		Name:          "nm",
		User:          "us",
		Password:      "ps",
		ConnTimeoutMS: 67,
	}

	dsn := c.GetDSN()
	assert.Equal(t, "us:ps@tcp(hst:45)/nm?timeout=67ms", dsn)
}
