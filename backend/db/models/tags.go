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

	util.LogQuery(ctxTimeout, "CreateTag:", stmt)

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

func (m *Models) GetTagsByBlogID(ctx context.Context, blog_id int) ([]Tag, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(m.config.Timeout)*time.Second)
	defer cancel()

	stmt := `
	SELECT 
		tags.id, 
		tags.created_at, 
		tags.updated_at, 
		tags.name, 
		tags.description, 
		tags.slug 
	FROM tags INNER JOIN blog_tags
	WHERE 
		(blog_tags.blog_id = ?) AND (blog_tags.tag_id = tags.id);
	`

	util.LogQuery(ctxTimeout, "GetTagsByBlogID:", stmt)

	rows, err := m.db.QueryContext(
		ctxTimeout,
		stmt,
		blog_id,
	)
	if err != nil {
		return []Tag{}, fmt.Errorf("GetTagsByBlogID: query context failed: %w", err)
	}

	result := []Tag{}
	for {
		tag := Tag{}
		if next := rows.Next(); next != true {
			break
		}
		err := rows.Scan(
			&tag.ID,
			&tag.Created_at,
			&tag.Updated_at,
			&tag.Name,
			&tag.Description,
			&tag.Slug,
		)
		if err != nil {
			if err := rows.Close(); err != nil {
				return []Tag{}, fmt.Errorf("GetTagsByBlogID: close rows error: %w", err)
			}
			return []Tag{}, fmt.Errorf("GetTagsByBlogID: scan error: %w", err)
		}
		result = append(result, tag)
	}

	return result, nil
}
