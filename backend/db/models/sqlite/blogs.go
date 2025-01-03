package sqlite

import (
	"blog/entities"
	"blog/util"
	"context"
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"
)

type Blogs struct{}

func NewBlogs() *Blogs {
	return &Blogs{}
}

func (b *Blogs) Create(ctx context.Context, tx *sql.Tx, blog entities.InBlog) (*entities.Blog, error) {

	stmt := `
	INSERT INTO blogs
	(
		title,
		content,
		content_md5,
		description,
		slug,
		pined,
		visible
	)
	VALUES
	( ?, ?, ?, ?, ?, ?, ?)
	RETURNING *;
	`

	util.LogQuery(ctx, "CreateBlog:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		blog.Title,
		blog.Content,
		blog.ContentMD5,
		blog.Description,
		blog.Slug,
		blog.Pined,
		blog.Visible,
	)
	if err := row.Err(); err != nil {
		return &entities.Blog{}, fmt.Errorf("Create: insert blog failed: %w", err)
	}

	newBlog, scanErr := scanBlog(row)
	if scanErr != nil {
		return &entities.Blog{}, fmt.Errorf("Create: scan error: %w", scanErr)
	}

	return newBlog, nil
}
func (b *Blogs) CreateWithID(ctx context.Context, tx *sql.Tx, blog entities.InBlog, id int) (*entities.Blog, error) {

	stmt := `
	INSERT INTO blogs
	(
		id,
		title,
		content,
		content_md5,
		description,
		slug,
		pined,
		visible
	)
	VALUES
	( ?, ?, ?, ?, ?, ?, ?, ?)
	RETURNING *;
	`

	util.LogQuery(ctx, "CreateWithID:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		id,
		blog.Title,
		blog.Content,
		blog.ContentMD5,
		blog.Description,
		blog.Slug,
		blog.Pined,
		blog.Visible,
	)
	if err := row.Err(); err != nil {
		return &entities.Blog{}, fmt.Errorf("CreateWithID: insert blog failed: %w", err)
	}

	newBlog, scanErr := scanBlog(row)
	if scanErr != nil {
		return &entities.Blog{}, fmt.Errorf("CreateWithID: scan error: %w", scanErr)
	}

	return newBlog, nil
}

func (b *Blogs) Update(ctx context.Context, tx *sql.Tx, blog entities.InBlog, id int) (*entities.Blog, error) {
	stmt := `
	UPDATE blogs 
	SET
		title = ?,
		content = ?,
		content_md5 = ?,
		description = ?,
		slug = ?,
		pined = ?,
		visible = ?
	WHERE 
		id = ?
	RETURNING *;
	`
	util.LogQuery(ctx, "UpdateBlog:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		blog.Title,
		blog.Content,
		blog.ContentMD5,
		blog.Description,
		blog.Slug,
		blog.Pined,
		blog.Visible,
		id,
	)
	if err := row.Err(); err != nil {
		return &entities.Blog{}, fmt.Errorf("Update: update blog failed: %w", err)
	}

	newBlog, err := scanBlog(row)
	if err != nil {
		return &entities.Blog{}, fmt.Errorf("Update: scan blog failed: %w", err)
	}

	return newBlog, nil
}

// only return visible and none soft deleted blogs
func (b *Blogs) Get(ctx context.Context, db *sql.DB, id int) (*entities.Blog, error) {
	stmt := `
	SELECT * FROM blogs WHERE id = ? AND visible = 1 AND deleted_at = "";
	`
	util.LogQuery(ctx, "GetBlog:", stmt)

	row := db.QueryRowContext(ctx, stmt, id)
	if err := row.Err(); err != nil {
		return &entities.Blog{}, fmt.Errorf("Get: get blog failed: %w", err)
	}

	blog, err := scanBlog(row)
	if err != nil {
		return &entities.Blog{}, fmt.Errorf("Get: scan blog failed: %w", err)
	}

	return blog, nil
}

// only return visible and none soft deleted blogs
func (b *Blogs) List(ctx context.Context, db *sql.DB) ([]entities.Blog, error) {
	stmt := `
	SELECT
		id,
		created_at,
		updated_at,
		deleted_at,
		title,
		content_md5,
		description,
		slug,
		pined,
		visible
	FROM blogs 
	WHERE visible = 1 AND deleted_at = ""
	ORDER BY updated_at DESC;
	`
	util.LogQuery(ctx, "ListBlogs:", stmt)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return []entities.Blog{}, fmt.Errorf("List: list blogs failed: %w", err)
	}

	result := []entities.Blog{}
	for {
		if !rows.Next() {
			break
		}
		blog, err := scanBlogRows(rows)
		if err != nil {
			return []entities.Blog{}, fmt.Errorf("List: scan blogs failed: %w", err)
		}
		result = append(result, *blog)
	}

	return result, nil
}

// only return visible and none soft deleted blogs
func (b *Blogs) ListByTopicIDs(ctx context.Context, db *sql.DB, topicIDs []int) ([]entities.Blog, error) {
	valueStrings := make([]string, 0, len(topicIDs))
	valueArgs := make([]any, 0, len(topicIDs)+1)

	for _, id := range topicIDs {
		valueStrings = append(valueStrings, "?")
		valueArgs = append(valueArgs, id)
	}

	valueArgs = append(valueArgs, len(topicIDs))

	// Only return blogs that has relation with all input topics
	stmt := fmt.Sprintf(
		`
	SELECT 
		id,
		created_at,
		updated_at,
		deleted_at,
		title,
		content_md5,
		description,
		slug,
		pined,
		visible
	FROM blogs
	WHERE id IN (
		SELECT blog_id FROM (
			SELECT blog_id,COUNT(blog_id) as count FROM blog_topics
			WHERE topic_id IN (%s) 
			GROUP BY blog_id
		) WHERE count = ?
	)
	AND visible = 1 AND deleted_at = ""
	ORDER BY updated_at DESC`,
		strings.Join(valueStrings, ","),
	)

	util.LogQuery(ctx, "ListBlogsByTopicIDs:", stmt)

	rows, err := db.QueryContext(ctx, stmt, valueArgs...)
	if err != nil {
		return []entities.Blog{}, fmt.Errorf("ListByTopicIDs: list blogs failed: %w", err)
	}

	result := []entities.Blog{}
	for {
		if !rows.Next() {
			break
		}
		blog, err := scanBlogRows(rows)
		if err != nil {
			return []entities.Blog{}, fmt.Errorf("ListByTopicIDs: scan row failed: %w", err)
		}
		result = append(result, *blog)
	}

	return result, nil
}

// only return visible and none soft deleted blogs
func (b *Blogs) ListByTopicAndTagIDs(ctx context.Context, db *sql.DB, topicIDs, tagIDs []int) ([]entities.Blog, error) {
	topicValueStrings := make([]string, 0, len(topicIDs))
	topicValueArgs := make([]any, 0, len(topicIDs))
	for _, id := range topicIDs {
		topicValueStrings = append(topicValueStrings, "?")
		topicValueArgs = append(topicValueArgs, id)
	}

	tagValueStrings := make([]string, 0, len(tagIDs))
	tagValueArgs := make([]any, 0, len(tagIDs))
	for _, id := range tagIDs {
		tagValueStrings = append(tagValueStrings, "?")
		tagValueArgs = append(tagValueArgs, id)
	}

	valueArgs := slices.Concat(topicValueArgs, tagValueArgs)
	valueArgs = append(valueArgs, len(topicIDs)*len(tagIDs))

	stmt := fmt.Sprintf(
		`
	SELECT 
		id,
		created_at,
		updated_at,
		deleted_at,
		title,
		content_md5,
		description,
		slug,
		pined,
		visible
	FROM blogs
	WHERE id IN (
		SELECT blog_id FROM (
			SELECT blog_id, COUNT(blog_id) AS count FROM ( 
				SELECT * FROM blog_topics JOIN blog_tags ON blog_topics.blog_id = blog_tags.blog_id 
			)
			WHERE topic_id IN (%s) AND tag_id IN (%s)
			GROUP BY blog_id
		) WHERE count = ?
	) 
	AND visible = 1 
	AND deleted_at = ""
	ORDER BY updated_at DESC;`,
		strings.Join(topicValueStrings, ","),
		strings.Join(tagValueStrings, ","),
	)

	util.LogQuery(ctx, "ListBlogsByTopicAndTagIDs:", stmt)

	rows, err := db.QueryContext(ctx, stmt, valueArgs...)
	if err != nil {
		return []entities.Blog{}, fmt.Errorf("ListByTopicAndTagIDs: list blogs failed: %w", err)
	}

	result := []entities.Blog{}
	for {
		if !rows.Next() {
			break
		}
		blog, err := scanBlogRows(rows)
		if err != nil {
			return []entities.Blog{}, fmt.Errorf("ListByTopicAndTagIDs: scan row failed: %w", err)
		}
		result = append(result, *blog)
	}

	return result, nil
}

// return blogs regardless of visiblility and soft delete status
func (b *Blogs) AdminGet(ctx context.Context, db *sql.DB, id int) (*entities.Blog, error) {
	stmt := `
	SELECT * FROM blogs WHERE id = ?;
	`
	util.LogQuery(ctx, "AdminGetBlog:", stmt)

	row := db.QueryRowContext(ctx, stmt, id)
	if err := row.Err(); err != nil {
		return &entities.Blog{}, fmt.Errorf("AdminGet: get blog failed: %w", err)
	}

	blog, err := scanBlog(row)
	if err != nil {
		return &entities.Blog{}, fmt.Errorf("AdminGet: scan blog failed: %w", err)
	}

	return blog, nil
}

// return blogs regardless of visiblility and soft delete status
func (b *Blogs) AdminList(ctx context.Context, db *sql.DB) ([]entities.Blog, error) {
	stmt := `
	SELECT 
		id,
		created_at,
		updated_at,
		deleted_at,
		title,
		content_md5,
		description,
		slug,
		pined,
		visible
	FROM blogs ORDER BY updated_at DESC;`

	util.LogQuery(ctx, "AdminListBlogs:", stmt)

	rows, err := db.QueryContext(ctx, stmt)
	if err != nil {
		return []entities.Blog{}, fmt.Errorf("AdminList: list blogs failed: %w", err)
	}

	result := []entities.Blog{}
	for {
		if !rows.Next() {
			break
		}
		blog, err := scanBlogRows(rows)
		if err != nil {
			return []entities.Blog{}, fmt.Errorf("AdminList: scan blogs failed: %w", err)
		}
		result = append(result, *blog)
	}

	return result, nil
}

// return blogs regardless of visiblility and soft delete status
func (b *Blogs) AdminListByTopicIDs(ctx context.Context, db *sql.DB, topicIDs []int) ([]entities.Blog, error) {
	valueStrings := make([]string, 0, len(topicIDs))
	valueArgs := make([]any, 0, len(topicIDs)+1)

	for _, id := range topicIDs {
		valueStrings = append(valueStrings, "?")
		valueArgs = append(valueArgs, id)
	}

	valueArgs = append(valueArgs, len(topicIDs))

	// Only return blogs that has relation with all input topics
	stmt := fmt.Sprintf(
		`
	SELECT 
		id,
		created_at,
		updated_at,
		deleted_at,
		title,
		content_md5,
		description,
		slug,
		pined,
		visible
	FROM blogs
	WHERE id IN (
		SELECT blog_id FROM (
			SELECT blog_id,COUNT(blog_id) as count FROM blog_topics
			WHERE topic_id IN (%s)
			GROUP BY blog_id
		) WHERE count = ?
	)
	ORDER BY updated_at DESC`,
		strings.Join(valueStrings, ","),
	)

	util.LogQuery(ctx, "AdminListBlogsByTopicIDs:", stmt)

	rows, err := db.QueryContext(ctx, stmt, valueArgs...)
	if err != nil {
		return []entities.Blog{}, fmt.Errorf("AdminListBlogsByTopicIDs: list blogs failed: %w", err)
	}

	result := []entities.Blog{}
	for {
		if !rows.Next() {
			break
		}
		blog, err := scanBlogRows(rows)
		if err != nil {
			return []entities.Blog{}, fmt.Errorf("AdminListBlogsByTopicIDs: scan row failed: %w", err)
		}
		result = append(result, *blog)
	}

	return result, nil
}

// return blogs regardless of visiblility and soft delete status
func (b *Blogs) AdminListByTopicAndTagIDs(ctx context.Context, db *sql.DB, topicIDs, tagIDs []int) ([]entities.Blog, error) {
	topicValueStrings := make([]string, 0, len(topicIDs))
	topicValueArgs := make([]any, 0, len(topicIDs))
	for _, id := range topicIDs {
		topicValueStrings = append(topicValueStrings, "?")
		topicValueArgs = append(topicValueArgs, id)
	}

	tagValueStrings := make([]string, 0, len(tagIDs))
	tagValueArgs := make([]any, 0, len(tagIDs))
	for _, id := range tagIDs {
		tagValueStrings = append(tagValueStrings, "?")
		tagValueArgs = append(tagValueArgs, id)
	}

	valueArgs := slices.Concat(topicValueArgs, tagValueArgs)
	valueArgs = append(valueArgs, len(topicIDs)*len(tagIDs))

	stmt := fmt.Sprintf(
		`
	SELECT 
		id,
		created_at,
		updated_at,
		deleted_at,
		title,
		content_md5,
		description,
		slug,
		pined,
		visible
	FROM blogs
	WHERE id IN (
		SELECT blog_id FROM (
			SELECT blog_id, COUNT(blog_id) AS count FROM ( 
				SELECT * FROM blog_topics JOIN blog_tags ON blog_topics.blog_id = blog_tags.blog_id 
			)
			WHERE topic_id IN (%s) AND tag_id IN (%s)
			GROUP BY blog_id
		) WHERE count = ?	)
	ORDER BY updated_at DESC`,
		strings.Join(topicValueStrings, ","),
		strings.Join(tagValueStrings, ","),
	)

	util.LogQuery(ctx, "AdminListBlogsByTopicAndTagIDs:", stmt)

	rows, err := db.QueryContext(ctx, stmt, valueArgs...)
	if err != nil {
		return []entities.Blog{}, fmt.Errorf("AdminListByTopicAndTagIDs: list blogs failed: %w", err)
	}

	result := []entities.Blog{}
	for {
		if !rows.Next() {
			break
		}
		blog, err := scanBlogRows(rows)
		if err != nil {
			return []entities.Blog{}, fmt.Errorf("AdminListByTopicAndTagIDs: scan row failed: %w", err)
		}
		result = append(result, *blog)
	}

	return result, nil
}

// mark deleted_at with current timestamp (ISO 8061)
func (b *Blogs) SoftDelete(ctx context.Context, tx *sql.Tx, id int) (int, error) {
	ts := time.Now().UTC().Format("2006-01-02T15:04:05-07:00")
	stmt := `
	UPDATE blogs SET deleted_at = ? WHERE id = ? AND deleted_at = '';
	`
	util.LogQuery(ctx, "SoftDeleteBlog:", stmt)

	res, err := tx.ExecContext(
		ctx,
		stmt,
		ts,
		id,
	)
	if err != nil {
		return 0, fmt.Errorf("SoftDelete: mark timestamp failed: %w", err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("SoftDelete: get affected rows failed: %w", err)
	}

	return int(affectedRows), nil
}

func (b *Blogs) Delete(ctx context.Context, tx *sql.Tx, id int) (int, error) {
	stmt := `
	DELETE FROM blogs WHERE id = ? AND deleted_at <> '';
	`
	util.LogQuery(ctx, "DeleteBlog:", stmt)

	res, err := tx.ExecContext(
		ctx,
		stmt,
		id,
	)
	if err != nil {
		return 0, fmt.Errorf("DeleteBlog: delete blog failed: %w", err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("DeleteBlog: get affected rows failed: %w", err)
	}

	return int(affectedRows), nil
}

func (b *Blogs) DeleteNow(ctx context.Context, tx *sql.Tx, id int) (int, error) {
	stmt := `
	DELETE FROM blogs WHERE id = ?;
	`
	util.LogQuery(ctx, "DeleteBlogNow:", stmt)

	res, err := tx.ExecContext(
		ctx,
		stmt,
		id,
	)
	if err != nil {
		return 0, fmt.Errorf("DeleteBlogNow: delete blog failed: %w", err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("DeleteBlogNow: get affected rows failed: %w", err)
	}

	return int(affectedRows), nil
}

func (b *Blogs) RestoreDeleted(ctx context.Context, tx *sql.Tx, id int) (*entities.Blog, error) {
	stmt := `
	UPDATE blogs SET deleted_at = "" WHERE id = ? RETURNING *;
	`
	util.LogQuery(ctx, "RestoreDeletedBlog:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		id,
	)
	if err := row.Err(); err != nil {
		return &entities.Blog{}, fmt.Errorf("RestoreDeletedBlog: restrore deleted blog failed: %w", err)
	}

	blog, err := scanBlog(row)
	if err != nil {
		return &entities.Blog{}, fmt.Errorf("RestoreDeletedBlog: scan blog failed: %w", err)
	}

	return blog, nil
}

// Helper for scanning blog
func scanBlog(row *sql.Row) (*entities.Blog, error) {
	newBlog := entities.Blog{}
	err := row.Scan(
		&newBlog.ID,
		&newBlog.Created_at,
		&newBlog.Updated_at,
		&newBlog.Deleted_at,
		&newBlog.Title,
		&newBlog.Content,
		&newBlog.ContentMD5,
		&newBlog.Description,
		&newBlog.Slug,
		&newBlog.Pined,
		&newBlog.Visible,
	)
	if err != nil {
		return &entities.Blog{}, fmt.Errorf("scanBlog: scan blog failed: %w", err)
	}
	return &newBlog, nil
}

func scanBlogRows(rows *sql.Rows) (*entities.Blog, error) {
	newBlog := entities.Blog{}
	err := rows.Scan(
		&newBlog.ID,
		&newBlog.Created_at,
		&newBlog.Updated_at,
		&newBlog.Deleted_at,
		&newBlog.Title,
		&newBlog.ContentMD5,
		&newBlog.Description,
		&newBlog.Slug,
		&newBlog.Pined,
		&newBlog.Visible,
	)
	if err != nil {
		return &entities.Blog{}, fmt.Errorf("scanBlog: scan blog failed: %w", err)
	}
	return &newBlog, nil
}
