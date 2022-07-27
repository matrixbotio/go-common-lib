package logger

import (
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/matrixbotio/constants-lib"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var wg sync.WaitGroup
var logConfig = getLogs("https://config.matrixbot.io/public/log-levels.json")

func InitESLogger(sourceName string, logLevel string, esProto string, esHost string, esPort string,
	esIndex string) (*Logger, error) {
	esLogger := esLogger{
		proto: esProto,
		host:  esHost,
		port:  esPort,
		index: esIndex,
	}

	hostName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	logger, err := NewLogger(&esLogger, hostName, sourceName, logLevel), nil
	if err != nil {
		return nil, err
	}

	err = esLogger.initEs(logger)
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func NewLogger(dev interface{}, host string, source string, lowestLevelName ...string) *Logger {
	format, formatLen := getSuitableDatetimeFormat(logConfig["datetime_format"].(string))
	logLevels := make(map[string]*logLevelDesc)
	lowestLevel := 2
	if levelsSection, ok := logConfig["levels"].(map[string]interface{}); ok {
		for strLevel, element := range levelsSection {
			if level, err := strconv.Atoi(strLevel); err == nil {
				if elMap, ok := element.(map[string]interface{}); ok {
					logLevel := &logLevelDesc{
						Level:  level,
						Stderr: false,
					}
					if stderr, exists := elMap["stderr_format"]; exists {
						logLevel.Stderr = true
						logLevel.Format = stderr.(string)
					} else if stdout, exists := elMap["stdout_format"]; exists {
						logLevel.Format = stdout.(string)
					}
					levelName := elMap["name"].(string)
					logLevels[levelName] = logLevel
					if len(lowestLevelName) > 0 && lowestLevelName[0] == levelName {
						lowestLevel = level
					}
				}
			}
		}
	}
	return &Logger{
		Dev:         dev.(logDevice),
		Host:        host,
		Source:      source,
		DTFormat:    format,
		DTFormatLen: formatLen,
		LogLevels:   logLevels,
		LowestLevel: lowestLevel,
	}
}

func (l *esLogger) Send(data string) {
	if l.client == nil {
		return
	}

	_, err := l.client.Index(
		l.index,
		strings.NewReader(data),
		l.client.Index.WithRefresh("true"),
	)
	if err != nil {
		log.Println("failed to send log to ES: " + err.Error())
	}
}

func (l *esLogger) initEs(logger *Logger) error {
	if l.proto == "" {
		logger.Log("ElasticSearch protocol is not passed, initialising logger without ElasticSearch")
		return nil
	}

	var err error
	esConfig := elasticsearch.Config{
		Addresses: []string{
			l.proto + "://" + l.host + ":" + l.port,
		},
	}
	l.client, err = elasticsearch.NewClient(esConfig)
	if err != nil {
		return err
	}
	return nil
}

func getSuitableDatetimeFormat(format string) (string, int) {
	return strings.NewReplacer("YYYY", "2006", "MM", "01", "dd", "02", "HH", "15", "mm", "04", "ss", "05", "SSS", "999").Replace(format), utf8.RuneCountInString(format)
}

func getLogs(url string) map[string]interface{} {
	storage := make(map[string]interface{})
	getJSON(url, &storage)
	return storage
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

	formattedTime := now.Format(l.DTFormat)
	formattedTime += strings.Repeat("0", l.DTFormatLen-utf8.RuneCountInString(formattedTime))
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

func AwaitLoggers() {
	wg.Wait()
}

// Verbose Very detailed logs
func (l *Logger) Verbose(message interface{}) {
	logLevel := l.LogLevels["verbose"]
	if l.LowestLevel > logLevel.Level {
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
	logLevel := l.LogLevels["log"]
	if l.LowestLevel > logLevel.Level {
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
	logLevel := l.LogLevels["warn"]
	if l.LowestLevel > logLevel.Level {
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
	logLevel := l.LogLevels["error"]
	if l.LowestLevel > logLevel.Level {
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
	logLevel := l.LogLevels["critical"]
	if l.LowestLevel > logLevel.Level {
		return
	}
	output := os.Stdout
	if logLevel.Stderr {
		output = os.Stderr
	}
	wg.Add(1)
	go l.baseWriter(message, output, logLevel.Format, logLevel.Level)
}
