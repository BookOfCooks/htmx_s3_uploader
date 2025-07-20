package app

import (
	"context"
	"embed"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

var DB *sqlx.DB

//go:embed migrations/*.sql
var embedMigrations embed.FS

func initDB(ctx context.Context) error {
	if _sqlx, err := sqlx.ConnectContext(ctx, "mysql", ENV.GOOSE_DBSTRING); err != nil {
		return fmt.Errorf("sqlx.ConnectContext: %v", err)
	} else {
		DB = _sqlx
	}

	if err := runMigrations(); err != nil {
		return fmt.Errorf("runMigrations: %w", err)
	}

	return nil
}

func runMigrations() error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(ENV.GOOSE_DRIVER); err != nil {
		return fmt.Errorf("goose.SetDialect: %w", err)
	}

	if err := goose.Up(DB.DB, "migrations"); err != nil {
		return fmt.Errorf("goose.Up: %w", err)
	}

	return nil
}
