package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type BlogTag struct {
	BlogID int `json:"blog_id"`
	TagID  int `json:"tag_id"`
}

func NewBlogTag(blogID, tagID int) *BlogTag {
	blogTag := &BlogTag{
		BlogID: blogID,
		TagID:  tagID,
	}
	return blogTag
}

func (m *Models) CreateBlogTags(ctx context.Context, blog Blog) (Blog, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(m.config.Timeout)*time.Second)
	defer cancel()

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

	tx, err := m.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return Blog{}, fmt.Errorf("CreateBlog: begin transaction error: %w", err)
	}

	row := tx.QueryRowContext(
		ctxTimeout,
		stmt,
		blog.Title,
		blog.Content,
		blog.Description,
		blog.Slug,
		blog.Pined,
		blog.Visible,
	)
	if err := row.Err(); err != nil {
		if err := tx.Rollback(); err != nil {
			return Blog{}, fmt.Errorf("CreateBlog: rollback error: %w", err)
		}
		return Blog{}, fmt.Errorf("CreateBlog: insert blog failed: %w", err)
	}

	newBlog := &Blog{}
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
		return Blog{}, fmt.Errorf("CreateBlog: scan error: %w", scanErr)
	}

	if err := tx.Commit(); err != nil {
		return Blog{}, fmt.Errorf("CreateBlog: commit error: %w", err)
	}

	return *newBlog, nil
}
