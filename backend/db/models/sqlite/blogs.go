package sqlite

import (
	"blog/entities"
	"blog/util"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (m *Models) CreateBlog(ctx context.Context, blog entities.InBlog) (entities.OutBlog, error) {
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

	util.LogQuery(ctxTimeout, "CreateBlog:", stmt)

	tx, err := m.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return entities.OutBlog{}, fmt.Errorf("CreateBlog: begin transaction error: %w", err)
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
			return entities.OutBlog{}, fmt.Errorf("CreateBlog: query rollback error: %w", err)
		}
		return entities.OutBlog{}, fmt.Errorf("CreateBlog: insert blog failed: %w", err)
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
		if err := tx.Rollback(); err != nil {
			return entities.OutBlog{}, fmt.Errorf("CreateBlog: scan rollback error: %w", err)
		}
		return entities.OutBlog{}, fmt.Errorf("CreateBlog: scan error: %w", scanErr)
	}

	if err := m.createBlogTags(ctxTimeout, tx, newBlog.ID, blog.Tags); err != nil {
		if err := tx.Rollback(); err != nil {
			return entities.OutBlog{}, fmt.Errorf("CreateBlog: insert blog_tags rollback error: %w", err)
		}
		return entities.OutBlog{}, fmt.Errorf("CreateBlog: insert blog_tags error: %w", err)
	}

	if err := m.createBlogTopics(ctxTimeout, tx, newBlog.ID, blog.Topics); err != nil {
		if err := tx.Rollback(); err != nil {
			return entities.OutBlog{}, fmt.Errorf("CreateBlog: insert blog_topics rollback error: %w", err)
		}
		return entities.OutBlog{}, fmt.Errorf("CreateBlog: insert blog_topics error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return entities.OutBlog{}, fmt.Errorf("CreateBlog: commit error: %w", err)
	}

	tags, err := m.GetTagsByBlogID(ctxTimeout, newBlog.ID)
	if err != nil {
		return entities.OutBlog{}, fmt.Errorf("CreateBlog: get tags by blog id error: %w", err)
	}
	topics, err := m.GetTopicsByBlogID(ctxTimeout, newBlog.ID)
	if err != nil {
		return entities.OutBlog{}, fmt.Errorf("CreateBlog: get topics by blog id error: %w", err)
	}

	outBlog := entities.NewOutBlog(*newBlog, tags, topics)
	return *outBlog, nil
}
