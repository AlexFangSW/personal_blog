package sqlite

import (
	"blog/util"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type BlogTopics struct{}

func NewBlogTopics() *BlogTopics {
	return &BlogTopics{}
}

func (b *BlogTopics) Create(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error {
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
	INSERT INTO blog_topics
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
