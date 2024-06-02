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
type Topic struct {
	ID          int    `json:"id"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

func (t *Topic) GenSlug() {
	t.Slug = slug.Make(t.Name)
}

func NewTopic(name, description string) *Topic {
	topic := &Topic{
		Name:        name,
		Description: description,
	}
	topic.GenSlug()
	return topic
}

func (m *Models) CreateTopic(ctx context.Context, topic Topic) (Topic, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(m.config.Timeout)*time.Second)
	defer cancel()

	stmt := `
	INSERT INTO topics
	(
		name,
		description,
		slug
	)
	VALUES
	( ?, ?, ? )
	RETURNING *;
	`

	if debug := slog.Default().Enabled(ctxTimeout, slog.LevelDebug); debug {
		fmt.Println("CreateTopic:", stmt)
	}

	tx, err := m.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return Topic{}, fmt.Errorf("CreateTopic: begin transaction error: %w", err)
	}

	row := tx.QueryRowContext(
		ctxTimeout,
		stmt,
		topic.Name,
		topic.Description,
		topic.Slug,
	)
	if err := row.Err(); err != nil {
		if err := tx.Rollback(); err != nil {
			return Topic{}, fmt.Errorf("CreateTopic: rollback error: %w", err)
		}
		return Topic{}, fmt.Errorf("CreateTopic: insert topic failed: %w", err)
	}

	newTopic := &Topic{}
	scanErr := row.Scan(
		&newTopic.ID,
		&newTopic.Created_at,
		&newTopic.Updated_at,
		&newTopic.Name,
		&newTopic.Description,
		&newTopic.Slug,
	)
	if scanErr != nil {
		return Topic{}, fmt.Errorf("CreateTopic: scan error: %w", scanErr)
	}

	if err := tx.Commit(); err != nil {
		return Topic{}, fmt.Errorf("CreateTopic: commit error: %w", err)
	}

	return *newTopic, nil
}

func (m *Models) GetTopicsByBlogID(ctx context.Context, blog_id int) ([]Topic, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(m.config.Timeout)*time.Second)
	defer cancel()

	stmt := `
	SELECT 
		topics.id, 
		topics.created_at, 
		topics.updated_at, 
		topics.name, 
		topics.description, 
		topics.slug 
	FROM topics INNER JOIN blog_topics
	WHERE 
		(blog_topics.blog_id = ?) AND (blog_topics.topic_id = topics.id);
	`

	if debug := slog.Default().Enabled(ctxTimeout, slog.LevelDebug); debug {
		fmt.Println("GetTopicsByBlogID:", stmt)
	}

	rows, err := m.db.QueryContext(
		ctxTimeout,
		stmt,
		blog_id,
	)
	if err != nil {
		return []Topic{}, fmt.Errorf("GetTopicsByBlogID: query context failed: %w", err)
	}

	result := []Topic{}
	for {
		topic := Topic{}
		if next := rows.Next(); next != true {
			break
		}
		err := rows.Scan(
			&topic.ID,
			&topic.Created_at,
			&topic.Updated_at,
			&topic.Name,
			&topic.Description,
			&topic.Slug,
		)
		if err != nil {
			if err := rows.Close(); err != nil {
				return []Topic{}, fmt.Errorf("GetTopicsByBlogID: close rows error: %w", err)
			}
			return []Topic{}, fmt.Errorf("GetTopicsByBlogID: scan error: %w", err)
		}
		result = append(result, topic)
	}

	return result, nil
}
