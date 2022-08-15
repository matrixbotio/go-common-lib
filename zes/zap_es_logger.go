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
	esBufSize       = 5 * 1024 // 5 kB
	esFlushInterval = 10 * time.Second
)

type Logger struct {
	readerError chan error
	pipeReader  *io.PipeReader
}

func InitGlobalLogger(isDebug, writeToES bool) (*Logger, error) {
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

	lgrCores := make([]zapcore.Core, 0)
	if isDebug {
		lgrCores = append(
			lgrCores,
			zapcore.NewCore(
				&encoder{zapcore.NewConsoleEncoder(newConsoleConfig())},
				os.Stdout,
				zapcore.DebugLevel,
			),
		)

		if writeToES {
			lgrCores = append(
				lgrCores,
				zapcore.NewCore(
					&encoder{zapcore.NewJSONEncoder(newJsonConfig())},
					&ws,
					zapcore.DebugLevel,
				),
			)
		}
	} else {
		lgrCores = append(
			lgrCores,
			zapcore.NewCore(
				&encoder{zapcore.NewJSONEncoder(newJsonConfig())},
				os.Stdout,
				zapcore.InfoLevel,
			),
		)

		if writeToES {
			lgrCores = append(
				lgrCores,
				zapcore.NewCore(
					&encoder{zapcore.NewJSONEncoder(newJsonConfig())},
					&ws,
					zapcore.InfoLevel,
				),
			)
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
