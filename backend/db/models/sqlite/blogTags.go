package sqlite

import (
	"blog/util"
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
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

	var values strings.Builder
	for i, id := range tagIDs {
		values.WriteString("(" + fmt.Sprint(blogID) + "," + fmt.Sprint(id) + ")")
		if i < len(tagIDs)-1 {
			values.WriteString(",")
		}
	}

	stmt := `
	REPLACE INTO blog_tags
	(
		blog_id,
		tag_id
	)
	VALUES 
	` + values.String() + ";"

	util.LogQuery(ctx, "CreateBlogTags:", stmt)

	_, insertErr := tx.ExecContext(
		ctx,
		stmt,
	)
	if insertErr != nil {
		return fmt.Errorf("Create: insert blog_tags failed: %w", insertErr)
	}

	return nil
}

func (b *BlogTags) Delete(ctx context.Context, tx *sql.Tx, blogID int, tagIDs []int) error {
	stmt := `DELETE FROM blog_tags WHERE blog_id = ? `

	if len(tagIDs) == 0 {
		// DELETE FROM blog_tags WHERE blog_id = ? ;
		stmt += ";"

	} else {
		var inIDs strings.Builder

		// DELETE FROM blog_tags WHERE blog_id = ? AND tag_id IN (x,x,x,x,x);
		inIDs.WriteString("AND tag_id IN (")
		for i, id := range tagIDs {
			inIDs.WriteString(strconv.Itoa(id))
			if i != len(tagIDs)-1 {
				inIDs.WriteString(",")
			}
		}
		inIDs.WriteString(");")

		stmt += inIDs.String()
	}
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

	var values strings.Builder
	for i, id := range tagIDs {
		values.WriteString(strconv.Itoa(id))
		if i != len(tagIDs) {
			values.WriteString(",")
		}
	}

	stmt := `
	DELETE FROM blog_tags
	WHERE 
		blog_id = ?
	AND tag_id NOT IN (` + values.String() + `)`

	util.LogQuery(ctx, "InverseDeleteBlogTags:", stmt)

	res, err := tx.ExecContext(ctx, stmt, blogID)
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
