package main

import (
	"errors"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/matrixbotio/go-common-lib/pkg/zes"
)

const writeToES = false

func main() {
	logger, err := zes.Init(
		writeToES,
		zap.String("env", "prod"),
		zap.String("stage", "stage"),
	)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	debugLgr := logger.New(zap.DebugLevel)
	debugLgr.Debug("debugLgr - DEBUG MESSAGE")
	debugLgr.Info("debugLgr - INFO MESSAGE")
	debugLgr.Warn("debugLgr - WARN MESSAGE")

	infoLgr := logger.New(zap.InfoLevel, zap.Bool("bool-additional", true))
	infoLgr.Debug("infoLgr - DEBUG MESSAGE")
	infoLgr.Info("infoLgr - INFO MESSAGE")
	infoLgr.Warn("infoLgr - WARN MESSAGE")

	warnLgr := logger.New(zap.WarnLevel)
	warnLgr.Debug("warnLgr - DEBUG MESSAGE")
	warnLgr.Info("warnLgr - INFO MESSAGE")
	warnLgr.Warn("warnLgr - WARN MESSAGE")

	// you can  SET and USE Global Variable everywhere
	zap.ReplaceGlobals(logger.New(zap.InfoLevel, zap.String("GLOBAL-add-field", "true")))

	go func() {
		startTime := time.Now()
		for {
			time.Sleep(1 * time.Second)

			zap.L().Info("test", zap.Int("number", rand.Int()), zap.Duration("duration-field", time.Now().Sub(startTime)))

			zap.L().Debug(
				"you will not see because logger is Info-level",
				zap.Int("number", rand.Int()),
				zap.Duration("duration-field", time.Now().Sub(startTime)),
			)
		}
	}()

	<-sigs

	writeError(debugLgr)

	debugLgr.Info("Application stopped")
}

func writeError(logger *zap.Logger) {
	logger.Error("OOOOPS", zap.Error(errors.New("this is an error")))
}
