package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
)

type Config struct {
	DB
	Server
}

type DB struct {
	Username string `env:"DB_USERNAME" env-required:"true"`
	Host     string `env:"DB_HOST" env-required:"true"`
	Port     string `env:"DB_PORT" env-required:"true"`
	DBName   string `env:"DB_NAME" env-required:"true"`
	SSLMode  string `env:"SSL_MODE" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
}

type Server struct {
	Port      string        `env:"SERV_PORT" env-required:"true"`
	ReadTime  time.Duration `env:"READ_TIME" env-required:"true"`
	WriteTime time.Duration `env:"WRITE_TIME" env-required:"true"`
}

func LoadConfig() (*Config, error) {
	godotenv.Load() //don't handle errors because we can upload via docker

	cfg := &Config{}
	if err := env.Parse(&cfg.DB); err != nil {
		return nil, fmt.Errorf("configuration reading error DB: %w", err)
	}

	if err := env.Parse(&cfg.Server); err != nil {
		return nil, fmt.Errorf("configuration reading error Server: %w", err)
	}

	return cfg, nil
}
