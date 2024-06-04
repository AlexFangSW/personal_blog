package interfaces

import (
	"context"
	"database/sql"
)

type BlogTagsModel interface {
	Create(ctx context.Context, tx *sql.Tx, blogID int, tagIDs []int) error
}
