package interfaces

import (
	"blog/entities"
	"context"
	"database/sql"
)

// Concrete implementations are at db/models/<db name>/
type TagsModel interface {
	Create(ctx context.Context, tx *sql.Tx, tag entities.Tag) (*entities.Tag, error)
	ListByBlogID(ctx context.Context, db *sql.DB, blogID int) ([]entities.Tag, error)
	ListSlugByBlogID(ctx context.Context, db *sql.DB, blogID int) ([]string, error)
	List(ctx context.Context, db *sql.DB) ([]entities.Tag, error)
	ListByTopicID(ctx context.Context, db *sql.DB, topicID int) ([]entities.Tag, error)
	Get(ctx context.Context, db *sql.DB, id int) (*entities.Tag, error)
	Update(ctx context.Context, tx *sql.Tx, tag entities.Tag, id int) (*entities.Tag, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) (int, error)
}
