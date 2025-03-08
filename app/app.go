package app

import (
	"context"
	"fmt"
)

func Init(ctx context.Context) error {
	if err := loadEnv(); err != nil {
		return fmt.Errorf("loadEnv: %w", err)
	}

	if err := initDB(ctx); err != nil {
		return fmt.Errorf("initDB: %w", err)
	}

	return nil
}
