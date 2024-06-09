package interfaces

import (
	"context"
	"database/sql"
)

// Concrete implementations are at db/models/<db name>/
type BlogTagsModel interface {
	Upsert(ctx context.Context, tx *sql.Tx, blogID int, tagIDs []int) error
	Delete(ctx context.Context, tx *sql.Tx, blogID int) error
	InverseDelete(ctx context.Context, tx *sql.Tx, blogID int, tagIDs []int) error
}
