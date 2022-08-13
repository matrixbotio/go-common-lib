package main

import (
	"errors"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/matrixbotio/go-common-lib/zes"
)

const isDebug = false

func main() {
	logger, err := zes.InitGlobalLogger(isDebug)
	if err != nil {
		panic(err)
	}
	defer logger.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	zap.L().Info("Application started")
	zap.L().Warn("WARN MESSAGE")
	zap.L().Debug("DEBUG MESSAGE - visible only if IsDebug=true")

	startTime := time.Now()

	go func() {
		for {
			time.Sleep(1 * time.Second)
			zap.L().Info("test", zap.Int("number", rand.Int()), zap.Duration("duration-field", time.Now().Sub(startTime)))
		}
	}()

	<-sigs

	writeError()

	zap.L().Info("Application stopped")
}

func writeError() {
	zap.L().Error("OOOOPS", zap.Error(errors.New("this is an error")))
}
