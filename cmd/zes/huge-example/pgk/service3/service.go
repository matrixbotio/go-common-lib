package service3

import "go.uber.org/zap"

type Service struct {
	logger *zap.Logger
}

func NewService(logger *zap.Logger) *Service {
	return &Service{logger: logger}
}

func (s *Service) SuperPowerFunc() {
	// doing something
	s.logger.Debug("something here for local development")

	// doing something
	s.logger.Info("usual log")
}
