package sqlite

import (
	"blog/entities"
	"blog/util"
	"context"
	"database/sql"
	"fmt"
)

type Tags struct{}

func NewTags() *Tags {
	return &Tags{}
}

func (t *Tags) Create(ctx context.Context, tx *sql.Tx, tag entities.Tag) (*entities.Tag, error) {

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

	util.LogQuery(ctx, "CreateTag:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		tag.Name,
		tag.Description,
		tag.Slug,
	)
	if err := row.Err(); err != nil {
		return &entities.Tag{}, fmt.Errorf("Create: insert tag failed: %w", err)
	}

	newTag := &entities.Tag{}
	scanErr := row.Scan(
		&newTag.ID,
		&newTag.Created_at,
		&newTag.Updated_at,
		&newTag.Name,
		&newTag.Description,
		&newTag.Slug,
	)
	if scanErr != nil {
		return &entities.Tag{}, fmt.Errorf("Create: scan error: %w", scanErr)
	}

	return newTag, nil
}

func (t *Tags) GetByBlogID(ctx context.Context, db *sql.DB, blog_id int) ([]entities.Tag, error) {
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

	util.LogQuery(ctx, "GetTagsByBlogID:", stmt)

	rows, err := db.QueryContext(
		ctx,
		stmt,
		blog_id,
	)
	if err != nil {
		return []entities.Tag{}, fmt.Errorf("GetByBlogID: query context failed: %w", err)
	}

	result := []entities.Tag{}
	for {
		tag := entities.Tag{}
		if !rows.Next() {
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
				return []entities.Tag{}, fmt.Errorf("GetByBlogID: close rows failed: %w", err)
			}
			return []entities.Tag{}, fmt.Errorf("GetByBlogID: scan failed: %w", err)
		}
		result = append(result, tag)
	}

	if err := rows.Err(); err != nil {
		return []entities.Tag{}, fmt.Errorf("GetByBlogID: rows iteration error: %w", err)
	}

	return result, nil
}

func (t *Tags) List(ctx context.Context, db *sql.DB) ([]entities.Tag, error) {

	stmt := `SELECT * FROM tags;`
	util.LogQuery(ctx, "ListTags:", stmt)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return []entities.Tag{}, fmt.Errorf("List: query failed: %w", err)
	}

	result := []entities.Tag{}
	for {
		if !rows.Next() {
			break
		}
		tag := entities.Tag{}
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
				return []entities.Tag{}, fmt.Errorf("List: close rows failed: %w", err)
			}
			return []entities.Tag{}, fmt.Errorf("List: scan failed: %w", err)
		}
		result = append(result, tag)
	}

	if err := rows.Err(); err != nil {
		return []entities.Tag{}, fmt.Errorf("List: rows iteration error: %w", err)
	}

	return result, nil
}

func (t *Tags) Get(ctx context.Context, db *sql.DB, id int) (*entities.Tag, error) {
	stmt := `SELECT * FROM tags WHERE id = ?;`
	util.LogQuery(ctx, "GetTag:", stmt)

	row := db.QueryRowContext(ctx, stmt, id)
	if err := row.Err(); err != nil {
		return &entities.Tag{}, fmt.Errorf("Get: query failed: %w", err)
	}

	tag := entities.Tag{}
	err := row.Scan(
		&tag.ID,
		&tag.Created_at,
		&tag.Updated_at,
		&tag.Name,
		&tag.Description,
		&tag.Slug,
	)
	if err != nil {
		return &entities.Tag{}, fmt.Errorf("Get: row scan failed: %w", err)
	}

	return &tag, nil
}

func (t *Tags) Update(ctx context.Context, db *sql.DB, tag entities.Tag) (*entities.Tag, error) {
	return &entities.Tag{}, nil
}
func (t *Tags) Delete(ctx context.Context, db *sql.DB, id int) error { return nil }
