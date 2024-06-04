package sqlite

import (
	"blog/entities"
	"blog/util"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (m *Models) CreateTopic(ctx context.Context, topic entities.Topic) (entities.Topic, error) {
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

	util.LogQuery(ctxTimeout, "CreateTopic:", stmt)

	tx, err := m.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return entities.Topic{}, fmt.Errorf("CreateTopic: begin transaction error: %w", err)
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
			return entities.Topic{}, fmt.Errorf("CreateTopic: rollback error: %w", err)
		}
		return entities.Topic{}, fmt.Errorf("CreateTopic: insert topic failed: %w", err)
	}

	newTopic := &entities.Topic{}
	scanErr := row.Scan(
		&newTopic.ID,
		&newTopic.Created_at,
		&newTopic.Updated_at,
		&newTopic.Name,
		&newTopic.Description,
		&newTopic.Slug,
	)
	if scanErr != nil {
		return entities.Topic{}, fmt.Errorf("CreateTopic: scan error: %w", scanErr)
	}

	if err := tx.Commit(); err != nil {
		return entities.Topic{}, fmt.Errorf("CreateTopic: commit error: %w", err)
	}

	return *newTopic, nil
}

func (m *Models) GetTopicsByBlogID(ctx context.Context, blog_id int) ([]entities.Topic, error) {
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

	util.LogQuery(ctxTimeout, "GetTopicsByBlogID:", stmt)

	rows, err := m.db.QueryContext(
		ctxTimeout,
		stmt,
		blog_id,
	)
	if err != nil {
		return []entities.Topic{}, fmt.Errorf("GetTopicsByBlogID: query context failed: %w", err)
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
				return []entities.Topic{}, fmt.Errorf("GetTopicsByBlogID: close rows error: %w", err)
			}
			return []entities.Topic{}, fmt.Errorf("GetTopicsByBlogID: scan error: %w", err)
		}
		result = append(result, topic)
	}

	if err := rows.Err(); err != nil {
		return []entities.Topic{}, fmt.Errorf("GetTopicsByBlogID: rows iteration error: %w", err)
	}

	return result, nil
}
