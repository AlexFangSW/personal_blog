package sqlite

import (
	"blog/util"
	"context"
	"database/sql"
	"fmt"
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
	return nil
}

func (b *BlogTags) InverseDelete(ctx context.Context, tx *sql.Tx, blogID int, tagIDs []int) error {
	return nil
}
