package models

import (
	"blog/structs"
	"database/sql"
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
