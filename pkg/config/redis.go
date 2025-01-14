package config

type RedisConfig struct {
	Host     string `envconfig:"REDIS_HOST" default:"0.0.0.0"`
	Port     string `envconfig:"REDIS_PORT" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD" default:""`
}
