package db

import "embed"

//go:embed migrations/sqlite/*.sql
var EmbedMigrationsSQLite embed.FS
