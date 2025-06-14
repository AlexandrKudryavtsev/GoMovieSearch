package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP    HTTP    `yaml:"http"`
		Log     Log     `yaml:"logger"`
		Elastic Elastic `yaml:"elastic"`
	}

	HTTP struct {
		Port string `yaml:"port"`
	}

	Log struct {
		Level       string `yaml:"level"`
		Destination string `yaml:"destination" env:"LOG_DESTINATION"`
	}

	Elastic struct {
		Addresses    []string      `yaml:"addresses" env:"ES_ADDRESSES" env-separator:","`
		Username     string        `yaml:"username" env:"ES_USERNAME"`
		Password     string        `yaml:"password" env:"ES_PASSWORD"`
		ConnTimeout  time.Duration `yaml:"conn_timeout" env:"ES_CONN_TIMEOUT"`
		ConnAttempts int           `yaml:"conn_attempts" env:"ES_CONN_ATTEMPTS"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("can't read yml config: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't read env config: %w", err)
	}

	if err := validateDestination(cfg.Log.Destination); err != nil {
		return nil, err
	}

	return cfg, nil
}

func validateDestination(destination string) error {
	destination = strings.ToLower(destination)
	if destination != "file" && destination != "console" {
		return fmt.Errorf("invalid log destination: %s. Use 'file' or 'console'", destination)
	}
	return nil
}
