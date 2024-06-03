package models

import (
	"blog/util"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type BlogTopic struct {
	BlogID  int `json:"blog_id"`
	TopicID int `json:"tag_id"`
}

func NewBlogTopic(blogID, topicID int) *BlogTopic {
	blogTopic := &BlogTopic{
		BlogID:  blogID,
		TopicID: topicID,
	}
	return blogTopic
}

// Used in conjunction with CreateBlog.
// DOES NOT rollback or commit transaction
func (m *Models) createBlogTopics(ctx context.Context, tx *sql.Tx, blogID int, topicIDs []int) error {
	if len(topicIDs) == 0 {
		return nil
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(m.config.Timeout)*time.Second)
	defer cancel()

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

	util.LogQuery(ctxTimeout, "createBlogTopics:", stmt)

	_, insertErr := tx.ExecContext(
		ctxTimeout,
		stmt,
	)
	if insertErr != nil {
		return fmt.Errorf("createBlogTopics: insert blog_topics failed: %w", insertErr)
	}

	return nil
}
