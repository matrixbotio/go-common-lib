package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type DB struct {
	Host          string `envconfig:"DB_HOST" default:"localhost"`
	Port          int    `envconfig:"DB_PORT" default:"3306"`
	Name          string `envconfig:"DB_NAME"`
	User          string `envconfig:"DB_USER"`
	Password      string `envconfig:"DB_PASSWORD"`
	ConnTimeoutMS int    `envconfig:"DB_CONN_TIMEOUT" default:"5000"`
}

func (db DB) GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?timeout=%dms",
		db.User, db.Password, db.Host, db.Port, db.Name, db.ConnTimeoutMS,
	)
}

func GetDBConfig(defaultName, defaultUser, defaultPassword string) (DB, error) {
	var c DB
	if err := envconfig.Process("", &c); err != nil {
		return DB{}, fmt.Errorf("parse envs: %w", err)
	}

	if c.Host == "" {
		return DB{}, EmptyFieldError{"host"}
	}

	if c.Name == "" {
		c.Name = defaultName
	}
	if c.User == "" {
		c.User = defaultUser
	}
	if c.Password == "" {
		c.Password = defaultPassword
	}

	return c, nil
}
