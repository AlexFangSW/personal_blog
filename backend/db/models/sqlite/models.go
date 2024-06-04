package sqlite

import (
	"blog/config"
	"context"
	"database/sql"
	"fmt"
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

func (m *Models) Prepare(ctx context.Context) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(m.config.Timeout)*time.Second)
	defer cancel()

	// enable sqlite foreign key
	_, err := m.db.ExecContext(ctxTimeout, "PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("PrepareSqlite: enable foreign key failed: %w", err)
	}

	return nil
}
