package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/matrixbotio/go-common-lib/cmd/zes/huge-example/pgk/service1and2"
	"github.com/matrixbotio/go-common-lib/cmd/zes/huge-example/pgk/service3"
	"github.com/matrixbotio/go-common-lib/pkg/zes"
)

const writeToES = false

func main() {
	// your flags that you need to set up what levels to what loggers
	applicationIsDebug := os.Getenv("APP_IS_DEBUG")
	service1IsDebug := os.Getenv("SERVICE1_IS_DEBUG")
	service2IsDebug := os.Getenv("SERVICE2_IS_DEBUG")
	service3IsDebug := os.Getenv("SERVICE3_IS_DEBUG")

	logger, err := zes.Init(
		writeToES,
		zap.String("env", "prod"),
		zap.String("stage", "stage"),
	)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	var appLgr *zap.Logger
	if applicationIsDebug == "1" {
		appLgr = logger.New(zap.DebugLevel)
	} else {
		appLgr = logger.New(zap.InfoLevel)
	}

	var service1Lgr *zap.Logger
	if service1IsDebug == "1" {
		service1Lgr = logger.New(zap.DebugLevel)
	} else {
		service1Lgr = logger.New(zap.InfoLevel)
	}

	var service2Lgr *zap.Logger
	if service2IsDebug == "1" {
		service2Lgr = logger.New(zap.DebugLevel)
	} else {
		service2Lgr = logger.New(zap.InfoLevel)
	}

	var service3Lgr *zap.Logger
	if service3IsDebug == "1" {
		service3Lgr = logger.New(zap.DebugLevel)
	} else {
		service3Lgr = logger.New(zap.InfoLevel)
	}

	// every service works with needed logs level as you wish
	src1 := service1and2.NewService1(service1Lgr)
	src2 := service1and2.NewService2(service2Lgr)
	src3 := service3.NewService(service3Lgr)

	appLgr.Info("start")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	src1.SuperPowerFunc()
	src2.SuperPowerFunc()
	src3.SuperPowerFunc()

	<-sigs
	appLgr.Info("stop")
}
