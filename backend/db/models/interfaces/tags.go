package interfaces

import (
	"blog/entities"
	"context"
	"database/sql"
)

type TagsModel interface {
	Create(ctx context.Context, tx *sql.Tx, tag entities.Tag) (*entities.Tag, error)
	GetByBlogID(ctx context.Context, db *sql.DB, blog_id int) ([]entities.Tag, error)
	List(ctx context.Context, db *sql.DB) ([]entities.Tag, error)
	Get(ctx context.Context, db *sql.DB, id int) (*entities.Tag, error)
	Update(ctx context.Context, tx *sql.Tx, tag entities.Tag) (*entities.Tag, error)
	Delete(ctx context.Context, tx *sql.Tx, id int) (int, error)
}
