package logger

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/matrixbotio/constants-lib"
)

var wg sync.WaitGroup
var logConfig = getLogConfig(logConfigUrl)

func NewLogger(dev interface{}, host string, source string, lowestLevelName ...string) *Logger {
	lowestLevel := parseLogLevel(lowestLevelName[0])
	partLogLevels := GetPartLogLevels()
	return &Logger{
		Dev:           dev.(logDevice),
		Host:          host,
		Source:        source,
		LowestLevel:   lowestLevel,
		PartLogLevels: partLogLevels,
	}
}

func AwaitLoggers() {
	wg.Wait()
}

// Verbose Very detailed logs
func (l *Logger) Verbose(message interface{}) {
	logLevel := logConfig.LogLevels["verbose"]
	if !l.isCorrectLevel(*logLevel) {
		return
	}
	output := os.Stdout
	if logLevel.Stderr {
		output = os.Stderr
	}
	wg.Add(1)
	go l.baseWriter(message, output, logLevel.Format, logLevel.Level)
}

// Log Important logs
func (l *Logger) Log(message interface{}) {
	logLevel := logConfig.LogLevels["log"]
	if !l.isCorrectLevel(*logLevel) {
		return
	}
	output := os.Stdout
	if logLevel.Stderr {
		output = os.Stderr
	}
	wg.Add(1)
	go l.baseWriter(message, output, logLevel.Format, logLevel.Level)
}

// Warn Something may go wrong
func (l *Logger) Warn(message interface{}) {
	logLevel := logConfig.LogLevels["warn"]
	if !l.isCorrectLevel(*logLevel) {
		return
	}
	output := os.Stdout
	if logLevel.Stderr {
		output = os.Stderr
	}
	wg.Add(1)
	go l.baseWriter(message, output, logLevel.Format, logLevel.Level)
}

// Failed to do something. This may cause problems!
func (l *Logger) Error(message interface{}) {
	logLevel := logConfig.LogLevels["error"]
	if !l.isCorrectLevel(*logLevel) {
		return
	}
	output := os.Stdout
	if logLevel.Stderr {
		output = os.Stderr
	}
	wg.Add(1)
	go l.baseWriter(message, output, logLevel.Format, logLevel.Level)
}

// Critical error. Node's shut down!
func (l *Logger) Critical(message interface{}) {
	logLevel := logConfig.LogLevels["critical"]
	if !l.isCorrectLevel(*logLevel) {
		return
	}
	output := os.Stdout
	if logLevel.Stderr {
		output = os.Stderr
	}
	wg.Add(1)
	go l.baseWriter(message, output, logLevel.Format, logLevel.Level)
}

func (l *Logger) baseWriter(message interface{}, output *os.File, template string, level int) {
	defer wg.Done()
	now := time.Now()

	sendObj := &sendMessageFormat{
		Source:    l.Source,
		Host:      l.Host,
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Level:     level,
	}

	if message == nil {
		sendObj.Message = "Log nil message. Please, don't log nils"
	} else if msg, ok := message.(string); ok {
		sendObj.Message = msg
	} else if err, ok := message.(*constants.APIError); ok {
		sendObj.Message = err.Message
		sendObj.Stack = err.Stack
		sendObj.Code = err.Code
	} else {
		sendObj.Message = "Logger error: can't cast provided message to string or *APIError"
	}

	formattedTime := now.Format(logConfig.DTFormat)
	formattedMessage := sendObj.Message
	if sendObj.Stack != nil {
		formattedMessage += "\n" + sendObj.Stack.(string)
	}

	_, err := output.WriteString(strings.NewReplacer(
		"%datetime%", formattedTime,
		"%message%", formattedMessage,
	).Replace(template) + "\n")

	if err != nil {
		log.Fatalf("failed to write to log file: %s", err.Error())
	}
	r, _ := json.Marshal(sendObj)

	l.Dev.Send(string(r))
}

func (l *Logger) isCorrectLevel(logLevel logLevelDesc) bool {
	packageName, functionName := getCallerInfo(3)
	pkgKey := toKey(packageName, "")
	set, isCorrect := l.isCorrectLevelForKey(pkgKey, logLevel)
	if set {
		return isCorrect
	}
	pkgFuncKey := toKey(packageName, functionName)
	set, isCorrect = l.isCorrectLevelForKey(pkgFuncKey, logLevel)
	if set {
		return isCorrect
	}
	return logLevel.Level >= l.LowestLevel
}

// IsCorrectLevelForKey returns 1. level is set, 2. level is correct
func (l *Logger) isCorrectLevelForKey(key string, logLevel logLevelDesc) (bool, bool) {
	level, exists := l.PartLogLevels[key]
	if exists {
		return true, logLevel.Level >= level
	}
	return false, false
}
