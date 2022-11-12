package gorm

import (
	"gorm.io/gorm/logger"

	"github.com/matrixbotio/go-common-lib/pkg/zes"
)

func NewLogger(loggersFactory *zes.Logger, debugMode bool) logger.Interface {
	var lvl logger.LogLevel
	if debugMode {
		lvl = logger.Info
	} else {
		lvl = logger.Silent
	}
	writer := newWriter(loggersFactory)

	return logger.New(
		writer,
		logger.Config{
			LogLevel:                  lvl,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}
