package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

type TestConfig struct {
	App    TestAppConfig
	Logs   LogsConfig
	DB     DB
	RMQ    RMQ
	Redis  RedisConfig
	Sentry SentryConfig
}

type TestAppConfig struct{}

func TestProcessConfigEmptySuccess(t *testing.T) {
	// given
	var cfg TestConfig

	// when
	err := ProcessConfig(&cfg)

	// then
	require.NoError(t, err)
}

func TestProcessConfigPointerError(t *testing.T) {
	// given
	var cfg TestConfig

	// when
	err := ProcessConfig(cfg)

	// then
	require.ErrorContains(t, err, "non-nil pointer")
}

func TestProcessConfigPointerInvalidError(t *testing.T) {
	// given
	var cfg = ""

	// when
	err := ProcessConfig(&cfg)

	// then
	require.ErrorContains(t, err, "pointer to a struct")
}

func TestProcessConfigLogsSuccess(t *testing.T) {
	// given
	var cfg TestConfig

	defer os.Clearenv()
	os.Setenv("LOG_LEVEL", string(LogLevelWarn))
	os.Setenv("LOG_LEVEL_RMQWORKERLIB", string(LogLevelFatal))

	// when
	err := ProcessConfig(&cfg)

	// then
	require.NoError(t, err)
	assert.Equal(t, zapcore.Level(1), cfg.Logs.GetAppLevel())
	assert.Equal(t, zapcore.Level(5), cfg.Logs.GetRMQLogLevel())
}

func TestProcessConfigDBSuccess(t *testing.T) {
	// given
	var cfg TestConfig
	var expectedDBName = "mydb"
	var expectedConnTimeout int = 4500

	defer os.Clearenv()
	os.Setenv("DB_NAME", expectedDBName)
	os.Setenv("DB_CONN_TIMEOUT", strconv.Itoa(expectedConnTimeout))

	// when
	err := ProcessConfig(&cfg)

	// then
	require.NoError(t, err)
	assert.Equal(t, expectedDBName, cfg.DB.Name)
	assert.Equal(t, expectedConnTimeout, cfg.DB.ConnTimeoutMS)
}
