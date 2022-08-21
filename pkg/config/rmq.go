package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type RMQ struct {
	UseTLS   bool   `envconfig:"AMQP_USE_TLS" default:"false"`
	Host     string `envconfig:"AMQP_HOST" default:"localhost"`
	Port     int    `envconfig:"AMQP_PORT" default:"5672"`
	User     string `envconfig:"AMQP_USER"`
	Password string `envconfig:"AMQP_PASSWORD"`
}

func GetRMQConfig() (RMQ, error) {
	var c RMQ
	if err := envconfig.Process("", &c); err != nil {
		return RMQ{}, fmt.Errorf("parse envs: %w", err)
	}

	if c.User == "" {
		return RMQ{}, EmptyFieldError{"user"}
	}

	if c.Password == "" {
		return RMQ{}, EmptyFieldError{"password"}
	}

	return c, nil
}
