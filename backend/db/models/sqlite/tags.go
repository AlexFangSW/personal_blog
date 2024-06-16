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

func (t *Tags) ListByBlogID(ctx context.Context, db *sql.DB, blogID int) ([]entities.Tag, error) {
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
		blogID,
	)
	if err != nil {
		return []entities.Tag{}, fmt.Errorf("ListByBlogID: query context failed: %w", err)
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
				return []entities.Tag{}, fmt.Errorf("ListByBlogID: close rows failed: %w", err)
			}
			return []entities.Tag{}, fmt.Errorf("ListByBlogID: scan failed: %w", err)
		}
		result = append(result, tag)
	}

	if err := rows.Err(); err != nil {
		return []entities.Tag{}, fmt.Errorf("ListByBlogID: rows iteration error: %w", err)
	}

	return result, nil
}

func (t *Tags) ListSlugByBlogID(ctx context.Context, db *sql.DB, blogID int) ([]string, error) {
	stmt := `
	SELECT 
		tags.slug 
	FROM tags INNER JOIN blog_tags
	WHERE 
		(blog_tags.blog_id = ?) AND (blog_tags.tag_id = tags.id);
	`

	util.LogQuery(ctx, "ListSlugByBlogID:", stmt)

	rows, err := db.QueryContext(
		ctx,
		stmt,
		blogID,
	)
	if err != nil {
		return []string{}, fmt.Errorf("ListSlugByBlogID: query context failed: %w", err)
	}

	result := []string{}
	for {
		slug := ""
		if !rows.Next() {
			break
		}
		err := rows.Scan(
			&slug,
		)
		if err != nil {
			if err := rows.Close(); err != nil {
				return []string{}, fmt.Errorf("ListSlugByBlogID: close rows failed: %w", err)
			}
			return []string{}, fmt.Errorf("ListSlugByBlogID: scan failed: %w", err)
		}
		result = append(result, slug)
	}

	if err := rows.Err(); err != nil {
		return []string{}, fmt.Errorf("ListSlugByBlogID: rows iteration error: %w", err)
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

func (t *Tags) ListByTopicID(ctx context.Context, db *sql.DB, topicID int) ([]entities.Tag, error) {
	stmt := `
	SELECT * FROM tags
	WHERE id IN (
		SELECT blog_tags.tag_id FROM blog_tags JOIN blog_topics 
		WHERE blog_tags.blog_id = blog_topics.blog_id AND blog_topics.topic_id = ?
		GROUP BY blog_tags.tag_id
	);`
	util.LogQuery(ctx, "ListByTopicID:", stmt)

	rows, err := db.QueryContext(ctx, stmt, topicID)
	if err != nil {
		return []entities.Tag{}, fmt.Errorf("ListByTopicID: query failed: %w", err)
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
				return []entities.Tag{}, fmt.Errorf("ListByTopicID: close rows failed: %w", err)
			}
			return []entities.Tag{}, fmt.Errorf("ListByTopicID: scan failed: %w", err)
		}
		result = append(result, tag)
	}

	if err := rows.Err(); err != nil {
		return []entities.Tag{}, fmt.Errorf("ListByTopicID: rows iteration error: %w", err)
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

func (t *Tags) Update(ctx context.Context, tx *sql.Tx, tag entities.Tag, id int) (*entities.Tag, error) {
	stmt := `
	UPDATE tags
	SET
		name = ?,
		description = ?,
		slug = ?
	WHERE 
		id = ?
	RETURNING *;
	`
	util.LogQuery(ctx, "UpdateTag:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		tag.Name,
		tag.Description,
		tag.Slug,
		id,
	)
	if err := row.Err(); err != nil {
		return &entities.Tag{}, fmt.Errorf("Update: update query failed: %w", err)
	}

	newTag := entities.Tag{}
	scanErr := row.Scan(
		&newTag.ID,
		&newTag.Created_at,
		&newTag.Updated_at,
		&newTag.Name,
		&newTag.Description,
		&newTag.Slug,
	)
	if scanErr != nil {
		return &entities.Tag{}, fmt.Errorf("Update: scan error: %w", scanErr)
	}

	return &newTag, nil
}

func (t *Tags) Delete(ctx context.Context, tx *sql.Tx, id int) (int, error) {
	stmt := `
	DELETE FROM tags WHERE id = ?;
	`
	util.LogQuery(ctx, "DeleteTags:", stmt)

	res, err := tx.ExecContext(ctx, stmt, id)
	if err != nil {
		return 0, fmt.Errorf("Delete: delete error: %w", err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("Delete: delete error: %w", err)
	}

	return int(affectedRows), nil
}
