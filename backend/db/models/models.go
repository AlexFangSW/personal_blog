package models

import (
	"blog/structs"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Models struct {
	db     *sql.DB
	config structs.DBSetting
}

func NewModels(db *sql.DB, config structs.DBSetting) *Models {
	return &Models{
		db:     db,
		config: config,
	}
}

func (m *Models) PrepareSqlite(ctx context.Context, timeout int) error {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// enable sqlite foreign key
	_, err := m.db.ExecContext(ctxTimeout, "PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("PrepareSqlite: enable foreign key failed: %w", err)
	}

	return nil
}
