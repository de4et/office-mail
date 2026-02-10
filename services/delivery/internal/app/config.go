package app

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	KAFKA_MAIL_TOPIC     string `env:"KAFKA_MAIL_TOPIC"`
	KAFKA_ADDRESSES      string `env:"KAFKA_ADDRESSES"`
	KAFKA_CONSUMER_GROUP string `env:"KAFKA_CONSUMER_GROUP"`
	REDIS_HOST           string `env:"REDIS_HOST"`
	REDIS_PORT           string `env:"REDIS_PORT"`
	REDIS_PASSWORD       string `env:"REDIS_PASSWORD"`
}

func LoadConfigFromENV() Config {
	var config Config
	if err := env.Parse(&config); err != nil {
		panic("Error reading the environment variables")
	}

	return config
}
