package interfaces

import (
	"blog/entities"
	"context"
	"database/sql"
)

type TopicsModel interface {
	Create(ctx context.Context, tx *sql.Tx, topic entities.Topic) (*entities.Topic, error)
	GetByBlogID(ctx context.Context, db *sql.DB, blog_id int) ([]entities.Topic, error)
	List(ctx context.Context, db *sql.DB) ([]entities.Topic, error)
	Get(ctx context.Context, db *sql.DB, id int) (*entities.Topic, error)
	Update(ctx context.Context, tx *sql.Tx, topic entities.Topic) (*entities.Topic, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) (int, error)
}
