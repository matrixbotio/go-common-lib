package zes

import (
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const (
	levelKey = "level"
)

type encoder struct {
	zapcore.Encoder
}

func (e *encoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	var level int

	switch entry.Level {
	case zapcore.DebugLevel:
		level = 1
	case zapcore.InfoLevel:
		level = 2
	case zapcore.WarnLevel:
		level = 3
	case zapcore.ErrorLevel:
		level = 4
	default:
		level = 5
	}

	return e.Encoder.EncodeEntry(entry, append(fields, zap.Int(levelKey, level)))
}
