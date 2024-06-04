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
		return &entities.Blog{}, fmt.Errorf("CreateBlog: insert blog failed: %w", err)
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
		return &entities.Blog{}, fmt.Errorf("CreateBlog: scan error: %w", scanErr)
	}

	return newBlog, nil
}
