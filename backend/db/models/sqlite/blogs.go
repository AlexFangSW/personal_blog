package sqlite

import (
	"blog/entities"
	"blog/util"
	"context"
	"database/sql"
	"fmt"
)

type Blogs struct{}

func NewBlogs() *Blogs {
	return &Blogs{}
}

func (b *Blogs) Create(ctx context.Context, tx *sql.Tx, blog entities.InBlog) (*entities.Blog, error) {

	stmt := `
	INSERT INTO blogs
	(
		title,
		content,
		description,
		slug,
		pined,
		visible
	)
	VALUES
	( ?, ?, ?, ?, ?, ? )
	RETURNING *;
	`

	util.LogQuery(ctx, "CreateBlog:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		blog.Title,
		blog.Content,
		blog.Description,
		blog.Slug,
		blog.Pined,
		blog.Visible,
	)
	if err := row.Err(); err != nil {
		return &entities.Blog{}, fmt.Errorf("Create: insert blog failed: %w", err)
	}

	newBlog := &entities.Blog{}
	scanErr := row.Scan(
		&newBlog.ID,
		&newBlog.Created_at,
		&newBlog.Updated_at,
		&newBlog.Deleted_at,
		&newBlog.Title,
		&newBlog.Content,
		&newBlog.Description,
		&newBlog.Slug,
		&newBlog.Pined,
		&newBlog.Visible,
	)
	if scanErr != nil {
		return &entities.Blog{}, fmt.Errorf("Create: scan error: %w", scanErr)
	}

	return newBlog, nil
}

func (b *Blogs) Update(ctx context.Context, tx *sql.Tx, blog entities.InBlog) (*entities.Blog, error) {
	return &entities.Blog{}, nil
}
func (b *Blogs) Get(ctx context.Context, db *sql.DB, id int) (*entities.Blog, error) {
	return &entities.Blog{}, nil
}

func (b *Blogs) List(ctx context.Context, db *sql.DB) ([]entities.Blog, error) {
	return []entities.Blog{}, nil
}

func (b *Blogs) ListByTopicID(ctx context.Context, db *sql.DB, topicID int) ([]entities.Blog, error) {
	return []entities.Blog{}, nil
}

func (b *Blogs) ListByTopicAndTagIDs(ctx context.Context, db *sql.DB, topicID, tagID []int) ([]entities.Blog, error) {
	return []entities.Blog{}, nil
}

func (b *Blogs) AdminGet(ctx context.Context, db *sql.DB, id int) (*entities.Blog, error) {
	return &entities.Blog{}, nil
}

func (b *Blogs) AdminList(ctx context.Context, db *sql.DB) ([]entities.Blog, error) {
	return []entities.Blog{}, nil
}
func (b *Blogs) AdminListByTopicID(ctx context.Context, db *sql.DB, topicID int) ([]entities.Blog, error) {
	return []entities.Blog{}, nil
}

func (b *Blogs) AdminListByTopicAndTagIDs(ctx context.Context, db *sql.DB, topicID, tagID []int) ([]entities.Blog, error) {
	return []entities.Blog{}, nil
}
func (b *Blogs) SoftDelete(ctx context.Context, tx *sql.Tx, id int) (int, error) { return 0, nil }
func (b *Blogs) Delete(ctx context.Context, tx *sql.Tx, id int) (int, error)     { return 0, nil }
func (b *Blogs) RestoreDeleted(ctx context.Context, tx *sql.Tx, id int) (*entities.Blog, error) {
	return &entities.Blog{}, nil
}
