package interfaces

import (
	"context"
	"database/sql"
)

type BlogTopicsModel interface {
	Create(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error
}
