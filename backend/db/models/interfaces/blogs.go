package interfaces

import (
	"blog/entities"
	"context"
	"database/sql"
)

type BlogsModel interface {
	Create(ctx context.Context, tx *sql.Tx, blog entities.InBlog) (*entities.Blog, error)
	Update(ctx context.Context, tx *sql.Tx, blog entities.InBlog) (*entities.Blog, error)
	Get(ctx context.Context, db *sql.DB, id int) (*entities.Blog, error)
	List(ctx context.Context, db *sql.DB) ([]entities.Blog, error)
	ListByTopicID(ctx context.Context, db *sql.DB, topicID int) ([]entities.Blog, error)
	ListByTopicAndTagIDs(ctx context.Context, db *sql.DB, topicID, tagID []int) ([]entities.Blog, error)
	AdminGet(ctx context.Context, db *sql.DB, id int) (*entities.Blog, error)
	AdminList(ctx context.Context, db *sql.DB) ([]entities.Blog, error)
	AdminListByTopicID(ctx context.Context, db *sql.DB, topicID int) ([]entities.Blog, error)
	AdminListByTopicAndTagIDs(ctx context.Context, db *sql.DB, topicID, tagID []int) ([]entities.Blog, error)
	SoftDelete(ctx context.Context, tx *sql.Tx, id int) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	RestoreDeleted(ctx context.Context, tx *sql.Tx, id int) (*entities.Blog, error)
}
