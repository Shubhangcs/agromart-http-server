package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	var (
		databaseHost     = os.Getenv("DATABASE_HOST")
		databasePort     = os.Getenv("DATABASE_PORT")
		databaseUser     = os.Getenv("DATABASE_USER")
		databaseName     = os.Getenv("DATABASE_NAME")
		databasePassword = os.Getenv("DATABASE_PASSWORD")
		databaseSSLMode  = os.Getenv("DATABASE_SSL_MODE")
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
