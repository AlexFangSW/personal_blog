package interfaces

import (
	"context"
	"database/sql"
)

// Concrete implementations are at db/models/<db name>/
type BlogTopicsModel interface {
	Upsert(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error
	Delete(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error
	InverseDelete(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error
}
