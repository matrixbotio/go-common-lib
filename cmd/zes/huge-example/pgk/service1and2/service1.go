package service1and2

import "go.uber.org/zap"

type Service1 struct {
	logger *zap.Logger
}

func NewService1(logger *zap.Logger) *Service1 {
	return &Service1{logger: logger}
}

func (s *Service1) SuperPowerFunc() {
	// doing something
	s.logger.Debug("something here for local development")

	// doing something
	s.logger.Info("usual log")
}
