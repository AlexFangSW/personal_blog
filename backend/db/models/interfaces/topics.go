package interfaces

import (
	"blog/entities"
	"context"
	"database/sql"
)

type TopicsModel interface {
	Create(ctx context.Context, tx *sql.Tx, topic entities.Topic) (*entities.Topic, error)
	GetByBlogID(ctx context.Context, db *sql.DB, blog_id int) ([]entities.Topic, error)
}
