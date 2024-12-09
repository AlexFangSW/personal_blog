package sqlite

import (
	"blog/util"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
)

type BlogTags struct{}

func NewBlogTags() *BlogTags {
	return &BlogTags{}
}

func (b *BlogTags) Upsert(ctx context.Context, tx *sql.Tx, blogID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(tagIDs))
	valueArgs := make([]any, 0, len(tagIDs)*2)

	for _, id := range tagIDs {
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, blogID)
		valueArgs = append(valueArgs, id)
	}

	stmt := fmt.Sprintf(
		`
	REPLACE INTO blog_tags 
	(
		blog_id,
		tag_id
	)
	VALUES %s`,
		strings.Join(valueStrings, ","),
	)
	util.LogQuery(ctx, "CreateBlogTags:", stmt)

	_, insertErr := tx.ExecContext(
		ctx,
		stmt,
		valueArgs...,
	)
	if insertErr != nil {
		return fmt.Errorf("Create: insert blog_tags failed: %w", insertErr)
	}

	return nil
}

func (b *BlogTags) Delete(ctx context.Context, tx *sql.Tx, blogID int) error {
	stmt := `DELETE FROM blog_tags WHERE blog_id = ?;`

	util.LogQuery(ctx, "DeleteBlogTags:", stmt)

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

func (b *BlogTags) InverseDelete(ctx context.Context, tx *sql.Tx, blogID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}

	valueStrings := make([]string, 0, len(tagIDs))
	valueArgs := make([]any, 0, len(tagIDs)+1)
	valueArgs = append(valueArgs, blogID)

	for _, id := range tagIDs {
		valueStrings = append(valueStrings, "?")
		valueArgs = append(valueArgs, id)
	}

	stmt := fmt.Sprintf(
		`
	DELETE FROM blog_tags
	WHERE 
		blog_id = ?
	AND tag_id NOT IN (%s)`,
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
