package service1and2

import "go.uber.org/zap"

type Service2 struct {
	logger *zap.Logger
}

func NewService2(logger *zap.Logger) *Service2 {
	return &Service2{logger: logger}
}

func (s *Service2) SuperPowerFunc() {
	// doing something
	s.logger.Debug("something here for local development")

	// doing something
	s.logger.Info("usual log")
}
