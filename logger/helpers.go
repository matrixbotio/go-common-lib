package logger

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"unicode/utf8"
)

func getJSON(url string, storage interface{}) {
	resp, err := http.Get(url)
	if err != nil {
		log.Panicln("Unable to download errors JSON file: " + url)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Panicln("Exception while closing errors JSON body: " + err.Error())
		}
	}(resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln("Exception while reading errors JSON body: " + err.Error())
		return
	}
	err = json.Unmarshal(body, storage)
	if err != nil {
		log.Panicln("Exception while unmarshalling errors JSON body: " + err.Error())
		return
	}
}

func getCallerInfo(skip int) (string, string) {
	pc, _, _, ok := runtime.Caller(skip)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		name := details.Name()
		if name != "" {
			lastSlashIndex := strings.LastIndex(name, "/")
			if lastSlashIndex != -1 {
				name = name[lastSlashIndex:]
			}
			nameParts := strings.Split(name, ".")
			packageName := nameParts[0]
			functionName := nameParts[len(nameParts)-1]
			return packageName, functionName
		}
	}
	return "", ""
}

func toKey(packageName string, functionName string) string {
	key := ""
	if packageName == "" {
		key = functionName
	} else if functionName == "" {
		key = packageName
	} else {
		key = packageName + "_" + functionName
	}
	return strings.ToLower(key)
}

func getSuitableDatetimeFormat(format string) (string, int) {
	return strings.NewReplacer("YYYY", "2006", "MM", "01", "dd", "02", "HH", "15", "mm", "04", "ss", "05", "SSS", "999").Replace(format), utf8.RuneCountInString(format)
}

func getLogConfig(url string) logConfiguration {
	logConfig := make(map[string]interface{})
	getJSON(url, &logConfig)
	dtFormat, dtFormatLen := getSuitableDatetimeFormat(logConfig["datetime_format"].(string))
	logLevels := make(map[string]*logLevelDesc)
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
				}
			}
		}
	}
	return logConfiguration{
		LogLevels:   logLevels,
		DTFormat:    dtFormat,
		DTFormatLen: dtFormatLen,
	}
}

func GetPartLogLevels() map[string]int {
	partLogLevels := make(map[string]int)
	for _, entry := range os.Environ() {
		if strings.HasPrefix(entry, "LOG_LEVEL_") {
			entryParts := strings.Split(entry, "=")
			key := strings.ToLower(strings.TrimPrefix(entryParts[0], "LOG_LEVEL_"))
			value := strings.TrimSpace(entryParts[1])
			partLogLevels[key] = parseLogLevel(value)
		}
	}
	return partLogLevels
}

func parseLogLevel(level string) int {
	logLevelDesc, exists := logConfig.LogLevels[level]
	if exists {
		return logLevelDesc.Level
	}
	return defaultLogLevel
}
