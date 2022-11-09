package gorm

import (
	"log"
	"os"

	"gorm.io/gorm/logger"
)

func NewLogger(debugMode bool) logger.Interface {
	var lvl logger.LogLevel
	if debugMode {
		lvl = logger.Info
	} else {
		lvl = logger.Silent
	}

	return logger.New(
		log.New(os.Stdout, logger.Green+"GORM "+logger.Reset, log.Ldate|log.Lmicroseconds|log.Lmsgprefix),
		logger.Config{
			LogLevel:                  lvl,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}
