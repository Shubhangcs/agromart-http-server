package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/shubhangcs/agromart-server/internal/env"
)

func Open() (*sql.DB, error) {
	var (
		databaseHost     = env.GetString("DATABASE_HOST", "localhost")
		databasePort     = env.GetString("DATABASE_PORT", "5432")
		databaseUser     = env.GetString("DATABASE_USER", "postgres")
		databaseName     = env.GetString("DATABASE_NAME", "postgres")
		databasePassword = env.GetString("DATABASE_PASSWORD", "postgres")
		databaseSSLMode  = env.GetString("DATABASE_SSL_MODE", "disable")
	)
	db, err := sql.Open(
		"pgx",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			databaseHost,
			databasePort,
			databaseUser,
			databaseName,
			databasePassword,
			databaseSSLMode,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}

	db.SetMaxOpenConns(env.GetInt("DATABASE_MAX_OPEN_CONNS", 25))
	db.SetMaxIdleConns(env.GetInt("DATABASE_MAX_IDLE_CONNS", 25))
	db.SetConnMaxLifetime(time.Duration(env.GetInt("DATABASE_CONN_MAX_LIFETIME_MIN", 15)) * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Open: %w", err)
	}
	return db, nil
}

func MigrateFS(db *sql.DB, migrationFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("Migrate: %w", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("Migrate: %w", err)
	}
	return nil
}
