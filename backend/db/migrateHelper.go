package db

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/pressly/goose/v3"
)

func Up(dbConn *sql.DB, migrations fs.FS, dialect string, path string) error {
	goose.SetLogger(goose.NopLogger())
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("Up: set dialect failed: %w", err)
	}

	if err := goose.Up(dbConn, path); err != nil {
		return fmt.Errorf("Up: up failed: %w", err)
	}

	return nil
}

func Down(dbConn *sql.DB, migrations fs.FS, dialect string, path string) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect(dialect); err != nil {
		return fmt.Errorf("Down: set dialect failed: %w", err)
	}

	if err := goose.Down(dbConn, path); err != nil {
		return fmt.Errorf("Down: down failed: %w", err)
	}

	return nil
}
