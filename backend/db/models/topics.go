package models

import (
	"context"
	"database/sql"
	"fmt"
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
	tag := &Topic{
		Name:        name,
		Description: description,
	}
	tag.GenSlug()
	return tag
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
