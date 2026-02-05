package app

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	PORT             int    `env:"PORT"`
	DB_HOST          string `env:"DB_HOST"`
	DB_PORT          int    `env:"DB_PORT"`
	DB_DATABASE      string `env:"DB_DATABASE"`
	DB_USERNAME      string `env:"DB_USERNAME"`
	DB_PASSWORD      string `env:"DB_PASSWORD"`
	KAFKA_MAIL_TOPIC string `env:"KAFKA_MAIL_TOPIC"`
	KAFKA_ADDRESSES  string `env:"KAFKA_ADDRESSES"`
}

func LoadConfigFromENV() Config {
	var config Config
	if err := env.Parse(&config); err != nil {
		panic("Error reading the environment variables")
	}

	return config
}
