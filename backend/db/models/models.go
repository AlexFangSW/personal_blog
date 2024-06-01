package models

import (
	"blog/structs"
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type Models struct {
	db     *sql.DB
	config structs.DBSetting
}

func New(db *sql.DB, config structs.DBSetting) *Models {
	return &Models{
		db:     db,
		config: config,
	}
}

func (m *Models) MigrateUp() error {
	driver, err := sqlite3.WithInstance(m.db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("MigrateUp: WithInstance error: %w", err)
	}
	migrate, err := migrate.NewWithDatabaseInstance("file://db/migrations", "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("MigrateUp: NewWithDatabaseInstance: %w", err)
	}
	if err := migrate.Up(); err != nil {
		return fmt.Errorf("MigrateUp: Up: %w", err)
	}
	return nil
}

func (m *Models) MigrateUpTo(target int) error {
	return nil
}

func (m *Models) MigrateDown() error {
	return nil
}
