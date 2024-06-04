package interfaces

import (
	"context"
	"database/sql"
)

type BlogTopicsModel interface {
	Upsert(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error
	Delete(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error
	InverseDelete(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error
}
