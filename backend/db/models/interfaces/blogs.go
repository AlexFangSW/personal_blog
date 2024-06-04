package interfaces

import (
	"blog/entities"
	"context"
	"database/sql"
)

type BlogsModel interface {
	Create(ctx context.Context, tx *sql.Tx, blog entities.InBlog) (*entities.Blog, error)
	Update(ctx context.Context, tx *sql.Tx, blog entities.InBlog) (*entities.Blog, error)
	Get(ctx context.Context, tx *sql.Tx, id int) (*entities.Blog, error)
	List(ctx context.Context, tx *sql.Tx) ([]entities.Blog, error)
	ListByTopicID(ctx context.Context, tx *sql.Tx, topicID int) ([]entities.Blog, error)
	AdminGet(ctx context.Context, tx *sql.Tx, id int) (*entities.Blog, error)
	AdminList(ctx context.Context, tx *sql.Tx) ([]entities.Blog, error)
	AdminListByTopicID(ctx context.Context, tx *sql.Tx, topicID int) ([]entities.Blog, error)
	SoftDelete(ctx context.Context, tx *sql.Tx, id int) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	RestoreDeleted(ctx context.Context, tx *sql.Tx, id int) (*entities.Blog, error)
}
