package gorm

import (
	"runtime"
	"strconv"

	"go.uber.org/zap"

	"github.com/matrixbotio/go-common-lib/pkg/zes"
)

type Writer struct {
	logger *zap.SugaredLogger
}

func newWriter(loggersFactory *zes.Logger) *Writer {
	l := loggersFactory.New(
		zap.DebugLevel,
		zap.String("service", "gormLogger"),
	)
	return &Writer{
		logger: l.Sugar(),
	}
}

func (w *Writer) Printf(message string, fields ...interface{}) {
	if len(fields) == 0 {
		return
	}
	fields[0] = fileWithLineNum()
	w.logger.Debugf(message, fields...)
}

func fileWithLineNum() string {
	if _, file, line, ok := runtime.Caller(5); ok {
		return file + ":" + strconv.FormatInt(int64(line), 10)
	}
	return "Unknown file & line"
}
