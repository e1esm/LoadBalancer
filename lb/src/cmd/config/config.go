package config

import (
	"github.com/caarlos0/env/v11"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port  int `env:"PORT" envDefault:"8000"`
	Redis struct {
		Address  string `env:"ADDRESS"`
		Password string `env:"PASSWORD"`
		DB       int    `env:"DB"`
	} `envPrefix:"REDIS_"`
	MaxCapacity           int    `env:"MAX_PARALLEL_REQUESTS"`
	ResetInterval         string `env:"RESET_INTERVAL"`
	TargetServiceBaseName string `env:"BASE_SERVICE_NAME"`
	TargetServicePort     int    `env:"TARGET_SERVICE_PORT"`
}

func New() *Config {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		log.WithError(err).Fatal("cfg init error")
	}

	return &cfg
}
