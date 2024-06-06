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

type BlogTopics struct{}

func NewBlogTopics() *BlogTopics {
	return &BlogTopics{}
}

func (b *BlogTopics) Upsert(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error {
	if len(topicIDs) == 0 {
		return nil
	}

	var values strings.Builder
	for i, id := range topicIDs {
		values.WriteString("(" + fmt.Sprint(blogID) + "," + fmt.Sprint(id) + ")")
		if i < len(topicIDs)-1 {
			values.WriteString(",")
		}
	}

	stmt := `
	REPLACE INTO blog_topics
	(
		blog_id,
		topic_id
	)
	VALUES 
	` + values.String() + ";"

	util.LogQuery(ctx, "CreateBlogTopics:", stmt)

	_, insertErr := tx.ExecContext(
		ctx,
		stmt,
	)
	if insertErr != nil {
		return fmt.Errorf("Create: insert blog_topics failed: %w", insertErr)
	}

	return nil
}

func (b *BlogTopics) Delete(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error {
	stmt := `DELETE FROM blog_topics WHERE blog_id = ? `

	if len(topicIDs) == 0 {
		// DELETE FROM blog_topics WHERE blog_id = ? ;
		stmt += ";"

	} else {
		var inIDs strings.Builder

		// DELETE FROM blog_topics WHERE blog_id = ? AND topic_id IN (x,x,x,x,x);
		inIDs.WriteString("AND topic_id IN (")
		for i, id := range topicIDs {
			inIDs.WriteString(strconv.Itoa(id))
			if i != len(topicIDs)-1 {
				inIDs.WriteString(",")
			}
		}
		inIDs.WriteString(");")

		stmt += inIDs.String()
	}
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

	var values strings.Builder
	for i, id := range topicIDs {
		values.WriteString(strconv.Itoa(id))
		if i != len(topicIDs) {
			values.WriteString(",")
		}
	}

	stmt := `
	DELETE FROM blog_topics
	WHERE 
		blog_id = ?
	AND topic_id NOT IN (` + values.String() + `)`

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
