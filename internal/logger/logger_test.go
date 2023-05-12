package logger

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatLogTime(t *testing.T) {
	// given
	cfgJSON := `{
		"datetime_format": "YYYY-MM-dd HH:mm:ss.999999999",
		"levels": {
		  "1": {
			"name": "verbose",
			"description": "Very detailed logs",
			"stdout_format": "\u001b[37;1m[%datetime%] \u001b[36;1mVerbose:\u001b[0m %message%"
		  },
		  "2": {
			"name": "log",
			"description": "Important logs",
			"stdout_format": "\u001b[37;1m[%datetime%] \u001b[32;1mInfo:\u001b[0m %message%"
		  },
		  "3": {
			"name": "warn",
			"description": "Something may go wrong",
			"stderr_format": "\u001b[37;1m[%datetime%] \u001b[33;1mWarning:\u001b[0m %message%"
		  },
		  "4": {
			"name": "error",
			"description": "Failed to do something. This may cause problems!",
			"stderr_format": "\u001b[37;1m[%datetime%] \u001b[31;1mERROR:\u001b[0m %message%"
		  },
		  "5": {
			"name": "critical",
			"description": "Critical error. Node's shutted down!",
			"stderr_format": "\u001b[37;1m[%datetime%] \u001b[35;1mCRITICAL:\u001b[0m %message%"
		  }
		}
	  }`

	timestamp := time.Unix(1683895354, 915012972)
	cfgRaw := map[string]interface{}{}
	require.NoError(t, json.Unmarshal([]byte(cfgJSON), &cfgRaw))

	cfg := parseLogConfig(cfgRaw)

	// when
	timeFormatted := timestamp.Format(cfg.DTFormat)

	// then
	assert.NotEmpty(t, timeFormatted)
}
