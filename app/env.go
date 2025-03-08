package app

import (
	"errors"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var ENV struct {
	DB_HOST string `env:"DB_HOST,required"`
	DB_USER string `env:"DB_USER,required"`
	DB_PASS string `env:"DB_PASS,required"`
	DB_NAME string `env:"DB_NAME,required"`
}

func loadEnv() error {
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("godotenv.Load: %w", err)
	}

	err := env.Parse(&ENV)
	if err != nil {
		var aggErr env.AggregateError
		if errors.As(err, &aggErr) {
			return errors.Join(aggErr.Errors...)
		} else {
			return err
		}
	}

	return nil
}
