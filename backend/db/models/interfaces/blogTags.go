package interfaces

import (
	"context"
	"database/sql"
)

type BlogTagsModel interface {
	Upsert(ctx context.Context, tx *sql.Tx, blogID int, tagIDs []int) error
	Delete(ctx context.Context, tx *sql.Tx, blogID int, tagIDs []int) error
	InverseDelete(ctx context.Context, tx *sql.Tx, blogID int, tagIDs []int) error
}
