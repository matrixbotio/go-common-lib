package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type App struct {
	Env string `envconfig:"ENV" default:"none"`
}

func GetAppConfig() (App, error) {
	var c App
	if err := envconfig.Process("", &c); err != nil {
		return App{}, fmt.Errorf("parse envs: %w", err)
	}

	if c.Env == "" {
		return App{}, EmptyFieldError{"env"}
	}

	return c, nil
}
