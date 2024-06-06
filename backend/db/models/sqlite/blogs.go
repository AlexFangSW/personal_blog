package sqlite

import (
	"blog/entities"
	"blog/util"
	"context"
	"database/sql"
	"fmt"
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
		description,
		slug,
		pined,
		visible
	)
	VALUES
	( ?, ?, ?, ?, ?, ? )
	RETURNING *;
	`

	util.LogQuery(ctx, "CreateBlog:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		blog.Title,
		blog.Content,
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

func (b *Blogs) Update(ctx context.Context, tx *sql.Tx, blog entities.InBlog) (*entities.Blog, error) {
	stmt := `
	UPDATE blogs 
	SET
		title = ?,
		content = ?,
		description = ?,
		slug = ?,
		pined = ?,
		visible = ?
	WHERE 
		id = ?
	RETURING *;
	`
	util.LogQuery(ctx, "UpdateBlog:", stmt)

	row := tx.QueryRowContext(
		ctx,
		stmt,
		blog.Title,
		blog.Content,
		blog.Description,
		blog.Slug,
		blog.Pined,
		blog.Visible,
		blog.ID,
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
	SELECT * FROM blogs WHERE visible = 1 AND deleted_at = "";
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

func (b *Blogs) ListByTopicID(ctx context.Context, db *sql.DB, topicID int) ([]entities.Blog, error) {
	return []entities.Blog{}, nil
}

func (b *Blogs) ListByTopicAndTagIDs(ctx context.Context, db *sql.DB, topicID, tagID []int) ([]entities.Blog, error) {
	return []entities.Blog{}, nil
}

func (b *Blogs) AdminGet(ctx context.Context, db *sql.DB, id int) (*entities.Blog, error) {
	return &entities.Blog{}, nil
}

func (b *Blogs) AdminList(ctx context.Context, db *sql.DB) ([]entities.Blog, error) {
	return []entities.Blog{}, nil
}
func (b *Blogs) AdminListByTopicID(ctx context.Context, db *sql.DB, topicID int) ([]entities.Blog, error) {
	return []entities.Blog{}, nil
}

func (b *Blogs) AdminListByTopicAndTagIDs(ctx context.Context, db *sql.DB, topicID, tagID []int) ([]entities.Blog, error) {
	return []entities.Blog{}, nil
}
func (b *Blogs) SoftDelete(ctx context.Context, tx *sql.Tx, id int) (int, error) { return 0, nil }
func (b *Blogs) Delete(ctx context.Context, tx *sql.Tx, id int) (int, error)     { return 0, nil }
func (b *Blogs) RestoreDeleted(ctx context.Context, tx *sql.Tx, id int) (*entities.Blog, error) {
	return &entities.Blog{}, nil
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
		&newBlog.Content,
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
