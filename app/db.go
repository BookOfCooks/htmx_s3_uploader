package app

import (
	"context"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func initDB(ctx context.Context) error {
	cfg := mysql.Config{
		Addr:                 ENV.DB_HOST,
		User:                 ENV.DB_USER,
		Passwd:               ENV.DB_PASS,
		DBName:               ENV.DB_NAME,
		Net:                  "tcp",
		AllowNativePasswords: true,
	}

	if _sqlx, err := sqlx.ConnectContext(ctx, "mysql", cfg.FormatDSN()); err != nil {
		return fmt.Errorf("sqlx.ConnectContext: %v", err)
	} else {
		DB = _sqlx
	}

	return nil
}
