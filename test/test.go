package main

import (
	"github.com/matrixbotio/constants-lib"
	goCommonLib "go-common-lib"
	"os"
)

func main() {
	logger, err := goCommonLib.InitLogger("test-common-lib", os.Getenv("LOG_LEVEL"), os.Getenv("ES_PROTO"),
		os.Getenv("ES_HOST"), os.Getenv("ES_PORT"), os.Getenv("ES_INDEX"))
	if err != nil {
		panic(err.Error())
	}
	logger.Log("Testing ES logging")
	constants.AwaitLoggers()
}
