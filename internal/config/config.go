package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

const (
	envFilename = ".env"
)

type Config struct {
	Ip                     string `env:"WOW_IP" envDefault:"0.0.0.0"`
	Port                   string `env:"WOW_PORT" envDefault:"8801"`
	LogLevel               string `env:"WOW_LOGLEVEL" envDefault:"warn"` // debug, info, warn, error, fatal, ""
	ZeroNumber             int    `env:"WOW_ZERO_NUMBER" envDefault:"4"`
	MaxHashCountIterations int    `env:"WOW_MAX_IT" envDefault:"10000000"`
	ClientRequestNumber    int    `env:"WOW_CLIENT_REQUEST_NUMBER" envDefault:"100"`
	ClientSendIntervalMs   int    `env:"WOW_CLIENT_SEND_INTERVAL_MS" envDefault:"1000"`
}

func New(path string) (*Config, error) {
	cfg := &Config{}

	if err := loadEnv(cfg, path); err != nil {
		return nil, err
	}

	return cfg, nil
}
func loadEnv(config interface{}, path string) error {
	if path == "" {
		path = envFilename
	}

	if fileExists(path) {
		err := godotenv.Load(path)
		if err != nil {
			return fmt.Errorf("error while loading existing %s file: %v", envFilename, err)
		}
	}

	if err := env.Parse(config); err != nil {
		return err
	}

	return nil
}
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
