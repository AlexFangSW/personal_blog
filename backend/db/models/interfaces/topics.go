package interfaces

import (
	"blog/entities"
	"context"
	"database/sql"
)

// Concrete implementations are at db/models/<db name>/
type TopicsModel interface {
	Create(ctx context.Context, tx *sql.Tx, topic entities.Topic) (*entities.Topic, error)
	ListByBlogID(ctx context.Context, db *sql.DB, blog_id int) ([]entities.Topic, error)
	ListSlugByBlogID(ctx context.Context, db *sql.DB, blogID int) ([]string, error)
	List(ctx context.Context, db *sql.DB) ([]entities.Topic, error)
	Get(ctx context.Context, db *sql.DB, id int) (*entities.Topic, error)
	Update(ctx context.Context, tx *sql.Tx, topic entities.Topic, id int) (*entities.Topic, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) (int, error)
}
