package models

import (
	"blog/util"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type BlogTag struct {
	BlogID int `json:"blog_id"`
	TagID  int `json:"tag_id"`
}

func NewBlogTag(blogID, tagID int) *BlogTag {
	blogTag := &BlogTag{
		BlogID: blogID,
		TagID:  tagID,
	}
	return blogTag
}

// Used in conjunction with CreateBlog.
// DOES NOT rollback or commit transaction
func (m *Models) createBlogTags(ctx context.Context, tx *sql.Tx, blogID int, tagIDs []int) error {
	if len(tagIDs) == 0 {
		return nil
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(m.config.Timeout)*time.Second)
	defer cancel()

	var values strings.Builder
	for i, id := range tagIDs {
		values.WriteString("(" + fmt.Sprint(blogID) + "," + fmt.Sprint(id) + ")")
		if i < len(tagIDs)-1 {
			values.WriteString(",")
		}
	}

	stmt := `
	INSERT INTO blog_tags
	(
		blog_id,
		tag_id
	)
	VALUES 
	` + values.String() + ";"

	util.LogQuery(ctxTimeout, "createBlogTags:", stmt)

	_, insertErr := tx.ExecContext(
		ctxTimeout,
		stmt,
	)
	if insertErr != nil {
		return fmt.Errorf("createBlogTags: insert blog_tags failed: %w", insertErr)
	}

	return nil
}
