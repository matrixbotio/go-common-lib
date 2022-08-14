package zes

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/matrixbotio/go-common-lib/logger"
)

const (
	timeKey       = "timestamp"
	messageKey    = "message"
	stacktraceKey = "stack"
	levelTextKey  = "level_text"
)

const (
	esBufSize       = 5 * 1024 // 5 kB
	esFlushInterval = 10 * time.Second
)

type Logger struct {
	readerError chan error
	pipeReader  *io.PipeReader
}

func InitGlobalLogger(isDebug bool) (*Logger, error) {
	elastic, err := logger.InitESLogger(
		"",
		"",
		os.Getenv("ES_PROTO"),
		os.Getenv("ES_HOST"),
		os.Getenv("ES_PORT"),
		os.Getenv("ES_INDEX"),
	)
	if err != nil {
		return nil, fmt.Errorf("init elastic: %w", err)
	}

	pipeReader, pipeWriter := io.Pipe()
	ws := zapcore.BufferedWriteSyncer{
		WS:            zapcore.AddSync(pipeWriter),
		Size:          esBufSize,
		FlushInterval: esFlushInterval,
	}

	var lgrCores []zapcore.Core
	if isDebug {
		lgrCfg := zap.NewDevelopmentEncoderConfig()
		lgrCfg.TimeKey = timeKey
		lgrCfg.MessageKey = messageKey
		lgrCfg.StacktraceKey = stacktraceKey
		lgrCfg.LevelKey = levelTextKey

		lgrCores = []zapcore.Core{
			zapcore.NewCore(&encoder{zapcore.NewConsoleEncoder(lgrCfg)}, os.Stdout, zapcore.DebugLevel),
		}
	} else {
		lgrCfg := zap.NewProductionEncoderConfig()
		lgrCfg.TimeKey = timeKey
		lgrCfg.MessageKey = messageKey
		lgrCfg.StacktraceKey = stacktraceKey
		lgrCfg.LevelKey = levelTextKey
		lgrCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		lgrCores = []zapcore.Core{
			zapcore.NewCore(&encoder{zapcore.NewJSONEncoder(lgrCfg)}, os.Stdout, zapcore.InfoLevel),
			zapcore.NewCore(&encoder{zapcore.NewJSONEncoder(lgrCfg)}, &ws, zapcore.InfoLevel),
		}
	}

	lgr := zap.New(
		zapcore.NewTee(lgrCores...),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	zap.ReplaceGlobals(lgr)

	readerError := make(chan error, 1)

	go func() {
		scanner := bufio.NewScanner(pipeReader)
		for scanner.Scan() {
			elastic.Dev.Send(scanner.Text())
		}

		if err := scanner.Err(); err != nil && !errors.Is(err, io.ErrClosedPipe) {
			readerError <- err
		}

		close(readerError)
	}()

	return &Logger{
		readerError: readerError,
		pipeReader:  pipeReader,
	}, nil
}

func (l *Logger) Close() error {
	zap.L().Sync()
	l.pipeReader.Close()

	return <-l.readerError
}
