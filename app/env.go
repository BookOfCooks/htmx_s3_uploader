package app

import (
	"errors"
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var ENV struct {
	HOST string `env:"HOST,required"`
	PORT int    `env:"PORT,required"`

	GOOSE_DRIVER        string `env:"GOOSE_DRIVER,required"`
	GOOSE_DBSTRING      string `env:"GOOSE_DBSTRING,required"`
	GOOSE_MIGRATION_DIR string `env:"GOOSE_MIGRATION_DIR,required"`
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
