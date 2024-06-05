package sqlite

import (
	"blog/entities"
	"blog/util"
	"context"
	"database/sql"
	"fmt"
)

type Topics struct{}

func NewTopics() *Topics {
	return &Topics{}
}

func (t *Topics) Create(ctx context.Context, tx *sql.Tx, topic entities.Topic) (*entities.Topic, error) {

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

	util.LogQuery(ctx, "CreateTopic:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		topic.Name,
		topic.Description,
		topic.Slug,
	)
	if err := row.Err(); err != nil {
		return &entities.Topic{}, fmt.Errorf("Create: insert topic failed: %w", err)
	}

	newTopic := entities.Topic{}
	scanErr := row.Scan(
		&newTopic.ID,
		&newTopic.Created_at,
		&newTopic.Updated_at,
		&newTopic.Name,
		&newTopic.Description,
		&newTopic.Slug,
	)
	if scanErr != nil {
		return &entities.Topic{}, fmt.Errorf("Create: scan error: %w", scanErr)
	}

	return &newTopic, nil
}

func (t *Topics) GetByBlogID(ctx context.Context, db *sql.DB, blog_id int) ([]entities.Topic, error) {

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

	util.LogQuery(ctx, "GetTopicsByBlogID:", stmt)

	rows, err := db.QueryContext(
		ctx,
		stmt,
		blog_id,
	)
	if err != nil {
		return []entities.Topic{}, fmt.Errorf("GetByBlogID: query context failed: %w", err)
	}

	result := []entities.Topic{}
	for {
		topic := entities.Topic{}
		if !rows.Next() {
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
				return []entities.Topic{}, fmt.Errorf("GetByBlogID: close rows error: %w", err)
			}
			return []entities.Topic{}, fmt.Errorf("GetByBlogID: scan error: %w", err)
		}
		result = append(result, topic)
	}

	if err := rows.Err(); err != nil {
		return []entities.Topic{}, fmt.Errorf("GetByBlogID: rows iteration error: %w", err)
	}

	return result, nil
}

func (t *Topics) List(ctx context.Context, db *sql.DB) ([]entities.Topic, error) {
	stmt := `SELECT * FROM topics;`
	util.LogQuery(ctx, "ListTopics:", stmt)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return []entities.Topic{}, fmt.Errorf("List: query failed: %w", err)
	}

	result := []entities.Topic{}
	for {
		if !rows.Next() {
			break
		}
		topic := entities.Topic{}
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
				return []entities.Topic{}, fmt.Errorf("List: close rows failed: %w", err)
			}
			return []entities.Topic{}, fmt.Errorf("List: scan failed: %w", err)
		}
		result = append(result, topic)
	}

	if err := rows.Err(); err != nil {
		return []entities.Topic{}, fmt.Errorf("List: rows iteration error: %w", err)
	}

	return result, nil
}
func (t *Topics) Get(ctx context.Context, db *sql.DB, id int) (*entities.Topic, error) {
	stmt := `SELECT * FROM topics WHERE id = ?;`
	util.LogQuery(ctx, "GetTopic:", stmt)

	row := db.QueryRowContext(ctx, stmt, id)
	if err := row.Err(); err != nil {
		return &entities.Topic{}, fmt.Errorf("Get: query failed: %w", err)
	}

	topic := entities.Topic{}
	err := row.Scan(
		&topic.ID,
		&topic.Created_at,
		&topic.Updated_at,
		&topic.Name,
		&topic.Description,
		&topic.Slug,
	)
	if err != nil {
		return &entities.Topic{}, fmt.Errorf("Get: row scan failed: %w", err)
	}

	return &topic, nil
}
func (t *Topics) Update(ctx context.Context, tx *sql.Tx, topic entities.Topic) (*entities.Topic, error) {
	stmt := `
	UPDATE topics
	SET
		name = ?,
		description = ?,
		slug = ?
	WHERE 
		id = ?
	RETURNING *;
	`
	util.LogQuery(ctx, "UpdateTopic:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		topic.Name,
		topic.Description,
		topic.Slug,
		topic.ID,
	)
	if err := row.Err(); err != nil {
		return &entities.Topic{}, fmt.Errorf("Update: update query failed: %w", err)
	}

	newTopic := entities.Topic{}
	scanErr := row.Scan(
		&newTopic.ID,
		&newTopic.Created_at,
		&newTopic.Updated_at,
		&newTopic.Name,
		&newTopic.Description,
		&newTopic.Slug,
	)
	if scanErr != nil {
		return &entities.Topic{}, fmt.Errorf("Update: scan error: %w", scanErr)
	}

	return &newTopic, nil
}

func (t *Topics) Delete(ctx context.Context, tx *sql.Tx, id int) (int, error) {
	stmt := `
	DELETE FROM topics WHERE id = ?;
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
