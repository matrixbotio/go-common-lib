package main

import (
	"os"

	cmnLogger "github.com/matrixbotio/go-common-lib/internal/logger"
)

func main() {
	logger, err := cmnLogger.InitESLogger("test-common-lib", os.Getenv("LOG_LEVEL"), os.Getenv("ES_PROTO"),
		os.Getenv("ES_HOST"), os.Getenv("ES_PORT"), os.Getenv("ES_INDEX"))
	if err != nil {
		panic(err.Error())
	}
	logger.Log("Testing ES logging")
	logger.Verbose("Testing Verbose")
	logger.Log("Testing Log")
	logger.Warn("Testing Warn")
	logger.Error("Testing Error")
	logger.Critical("Testing Critical")
	cmnLogger.AwaitLoggers()
}
