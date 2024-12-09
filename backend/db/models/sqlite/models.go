package sqlite

import (
	"blog/config"
	"blog/db"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type Models struct {
	db     *sql.DB
	config config.DBSetting
}

func New(db *sql.DB, config config.DBSetting) *Models {
	return &Models{
		db:     db,
		config: config,
	}
}

func (m *Models) Prepare(ctx context.Context, migrate bool) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(m.config.Timeout)*time.Second)
	defer cancel()

	// enable sqlite foreign key
	slog.Info("activate foreign keys")
	_, err := m.db.ExecContext(ctxTimeout, "PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("PrepareSqlite: enable foreign key failed: %w", err)
	}

	// migrate db
	if migrate {
		slog.Info("perform db migration")
		if err := db.Up(m.db, db.EmbedMigrationsSQLite, "sqlite3", "migrations/sqlite"); err != nil {
			return fmt.Errorf("PrepareSqlite: migrate up failed: %s", err)
		}
	}

	// start vacume

	return nil
}
