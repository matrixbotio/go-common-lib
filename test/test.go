package main

import (
	goCommonLib "github.com/matrixbotio/go-common-lib/logger"
	"os"
)

func main() {
	logger, err := goCommonLib.InitESLogger("test-common-lib", os.Getenv("LOG_LEVEL"), os.Getenv("ES_PROTO"),
		os.Getenv("ES_HOST"), os.Getenv("ES_PORT"), os.Getenv("ES_INDEX"))
	if err != nil {
		panic(err.Error())
	}
	logger.Log("Testing ES logging")
	goCommonLib.AwaitLoggers()
}
