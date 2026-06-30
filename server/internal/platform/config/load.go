package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func Load() (*Config, error) {
	_ = godotenv.Load()

	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		return nil, errors.New("APP_ENV is required")
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		return nil, errors.New("HTTP_PORT is required")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		return nil, errors.New("NATS_URL is required")
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		return nil, errors.New("LOG_LEVEL is required")
	}

	return &Config{
		App: AppConfig{
			Env: appEnv,
		},
		HTTP: HTTPConfig{
			Port: httpPort,
		},
		Database: DatabaseConfig{
			URL: databaseURL,
		},
		NATS: NATSConfig{
			URL: natsURL,
		},
		Logger: LoggerConfig{
			Level: logLevel,
		},
	}, nil
}
