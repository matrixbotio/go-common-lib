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

	"github.com/matrixbotio/go-common-lib/internal/logger"
)

const (
	esBufSize       = 5 * 1024 // 5 kB
	esFlushInterval = 10 * time.Second
)

type Logger struct {
	readerError            chan error
	pipeReader             *io.PipeReader
	logger                 *zap.Logger
	commonAdditionalFields []zap.Field
}

func Init(writeToES bool, commonAdditionalFields ...zap.Field) (*Logger, error) {
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

	lgrCores := []zapcore.Core{
		zapcore.NewCore(
			newEncoder(func() zapcore.Encoder {
				return zapcore.NewConsoleEncoder(newConsoleConfig())
			}),
			os.Stdout,
			zapcore.DebugLevel,
		),
	}

	if writeToES {
		lgrCores = append(
			lgrCores,
			zapcore.NewCore(
				newEncoder(func() zapcore.Encoder {
					return zapcore.NewJSONEncoder(newJsonConfig())
				}),
				&ws,
				zapcore.DebugLevel,
			),
		)
	}

	lgr := zap.New(
		zapcore.NewTee(lgrCores...),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

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
		readerError:            readerError,
		pipeReader:             pipeReader,
		logger:                 lgr,
		commonAdditionalFields: commonAdditionalFields,
	}, nil
}

func (l *Logger) Close() error {
	l.logger.Sync()
	l.pipeReader.Close()

	return <-l.readerError
}

func (l *Logger) New(level zapcore.Level, additionalFields ...zap.Field) *zap.Logger {
	if l.logger == nil {
		panic("must be initialized at first")
	}

	wrap := func(c zapcore.Core) zapcore.Core {
		n, err := zapcore.NewIncreaseLevelCore(c, level)
		if err != nil {
			panic("should never happen")
		}
		return n
	}

	return l.logger.
		WithOptions(zap.WrapCore(wrap)).
		With(append(l.commonAdditionalFields, additionalFields...)...)
}
