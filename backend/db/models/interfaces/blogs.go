package interfaces

import (
	"blog/entities"
	"context"
	"database/sql"
)

type BlogsModel interface {
	Create(ctx context.Context, tx *sql.Tx, blog entities.InBlog) (*entities.Blog, error)
}
