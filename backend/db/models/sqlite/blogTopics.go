package sqlite

import (
	"blog/util"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
)

type BlogTopics struct{}

func NewBlogTopics() *BlogTopics {
	return &BlogTopics{}
}

func (b *BlogTopics) Upsert(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error {
	if len(topicIDs) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(topicIDs))
	valueArgs := make([]any, 0, len(topicIDs)*2)

	for _, id := range topicIDs {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, blogID)
		valueArgs = append(valueArgs, id)
	}

	stmt := fmt.Sprintf(
		`
	REPLACE INTO blog_topics
	(
		blog_id,
		topic_id
	)
	VALUES %s`,
		strings.Join(valueStrings, ","),
	)

	util.LogQuery(ctx, "CreateBlogTopics:", stmt)

	_, insertErr := tx.ExecContext(
		ctx,
		stmt,
		valueArgs...,
	)
	if insertErr != nil {
		return fmt.Errorf("Create: insert blog_topics failed: %w", insertErr)
	}

	return nil
}

func (b *BlogTopics) Delete(ctx context.Context, tx *sql.Tx, blogID int) error {
	stmt := `DELETE FROM blog_topics WHERE blog_id = ?;`
	util.LogQuery(ctx, "DeleteBlogTopics:", stmt)

	res, err := tx.ExecContext(ctx, stmt, blogID)
	if err != nil {
		return fmt.Errorf("Delete: exec context failed: %w", err)
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Delete: aquire affected rows failed: %w", err)
	}
	slog.Debug("affected rows", "rows", affectedRows)

	return nil
}
func (b *BlogTopics) InverseDelete(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error {
	if len(topicIDs) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(topicIDs))
	valueArgs := make([]any, 0, len(topicIDs)+1)
	valueArgs = append(valueArgs, blogID)

	for _, id := range topicIDs {
		valueStrings = append(valueStrings, "?")
		valueArgs = append(valueArgs, id)
	}

	stmt := fmt.Sprintf(
		`
	DELETE FROM blog_topics
	WHERE 
		blog_id = ?
	AND topic_id NOT IN (%s) `,
		strings.Join(valueStrings, ","),
	)

	util.LogQuery(ctx, "InverseDeleteBlogTags:", stmt)

	res, err := tx.ExecContext(ctx, stmt, valueArgs...)
	if err != nil {
		return fmt.Errorf("InverseDelete: exec context failed: %w", err)
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("InverseDelete: aquire affected rows failed: %w", err)
	}
	slog.Debug("affected rows", "rows", affectedRows)

	return nil
}
