package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/gosimple/slug"
)

// xxx_at are all in ISO 8601.
type Tag struct {
	ID          int    `json:"id"`
	Created_at  string `json:"created_at"`
	Updated_at  string `json:"updated_at"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
}

func (t *Tag) GenSlug() {
	t.Slug = slug.Make(t.Name)
}

func NewTag(name, description string) *Tag {
	tag := &Tag{
		Name:        name,
		Description: description,
	}
	tag.GenSlug()
	return tag
}

func (m *Models) CreateTag(ctx context.Context, tag Tag) (Tag, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(m.config.Timeout)*time.Second)
	defer cancel()

	stmt := `
	INSERT INTO tags
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
		return Tag{}, fmt.Errorf("CreateTag: begin transaction error: %w", err)
	}

	row := tx.QueryRowContext(
		ctxTimeout,
		stmt,
		tag.Name,
		tag.Description,
		tag.Slug,
	)
	if err := row.Err(); err != nil {
		if err := tx.Rollback(); err != nil {
			return Tag{}, fmt.Errorf("CreateTag: rollback error: %w", err)
		}
		return Tag{}, fmt.Errorf("CreateTag: insert tag failed: %w", err)
	}

	newTag := &Tag{}
	scanErr := row.Scan(
		&newTag.ID,
		&newTag.Created_at,
		&newTag.Updated_at,
		&newTag.Name,
		&newTag.Description,
		&newTag.Slug,
	)
	if scanErr != nil {
		return Tag{}, fmt.Errorf("CreateTag: scan error: %w", scanErr)
	}

	if err := tx.Commit(); err != nil {
		return Tag{}, fmt.Errorf("CreateTag: commit error: %w", err)
	}

	return *newTag, nil
}
