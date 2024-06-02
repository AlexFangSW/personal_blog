package models

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/gosimple/slug"
)

// xxx_at are all in ISO 8601.
type Blog struct {
	ID          int    `json:"id"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
	Deleted_at  string `json:"deleted_at"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	Pined       bool   `json:"pined"`
	Visible     bool   `json:"visible"`
}

func (b *Blog) GenSlug() {
	b.Slug = slug.Make(b.Title)
}

func NewBlog(title, content, description string, pined, visible bool) *Blog {
	blog := &Blog{
		Title:       title,
		Content:     content,
		Description: description,
		Pined:       pined,
		Visible:     visible,
	}
	blog.GenSlug()
	return blog
}

type InBlog struct {
	Blog
	Tags   []int `json:"tags"`
	Topics []int `json:"topics"`
}

func NewInBlog(blog Blog, tags, topics []int) *InBlog {
	return &InBlog{
		Blog:   blog,
		Tags:   tags,
		Topics: topics,
	}
}

type OutBlog struct {
	Blog
	Tags   []Tag   `json:"tags"`
	Topics []Topic `json:"topics"`
}

func NewOutBlog(blog Blog, tags []Tag, topics []Topic) *OutBlog {
	return &OutBlog{
		Blog:   blog,
		Tags:   tags,
		Topics: topics,
	}
}

func (m *Models) CreateBlog(ctx context.Context, blog InBlog) (OutBlog, error) {
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

	if debug := slog.Default().Enabled(ctxTimeout, slog.LevelDebug); debug {
		fmt.Println("CreateBlog:", stmt)
	}

	tx, err := m.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return OutBlog{}, fmt.Errorf("CreateBlog: begin transaction error: %w", err)
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
			return OutBlog{}, fmt.Errorf("CreateBlog: query rollback error: %w", err)
		}
		return OutBlog{}, fmt.Errorf("CreateBlog: insert blog failed: %w", err)
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
		if err := tx.Rollback(); err != nil {
			return OutBlog{}, fmt.Errorf("CreateBlog: scan rollback error: %w", err)
		}
		return OutBlog{}, fmt.Errorf("CreateBlog: scan error: %w", scanErr)
	}

	if err := m.createBlogTags(ctxTimeout, tx, newBlog.ID, blog.Tags); err != nil {
		if err := tx.Rollback(); err != nil {
			return OutBlog{}, fmt.Errorf("CreateBlog: insert blog_tags rollback error: %w", err)
		}
		return OutBlog{}, fmt.Errorf("CreateBlog: insert blog_tags error: %w", err)
	}

	if err := m.createBlogTopics(ctxTimeout, tx, newBlog.ID, blog.Topics); err != nil {
		if err := tx.Rollback(); err != nil {
			return OutBlog{}, fmt.Errorf("CreateBlog: insert blog_topics rollback error: %w", err)
		}
		return OutBlog{}, fmt.Errorf("CreateBlog: insert blog_topics error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return OutBlog{}, fmt.Errorf("CreateBlog: commit error: %w", err)
	}

	tags, err := m.GetTagsByBlogID(ctxTimeout, newBlog.ID)
	if err != nil {
		return OutBlog{}, fmt.Errorf("CreateBlog: get tags by blog id error: %w", err)
	}
	topics, err := m.GetTopicsByBlogID(ctxTimeout, newBlog.ID)
	if err != nil {
		return OutBlog{}, fmt.Errorf("CreateBlog: get topics by blog id error: %w", err)
	}

	outBlog := NewOutBlog(*newBlog, tags, topics)
	return *outBlog, nil
}
