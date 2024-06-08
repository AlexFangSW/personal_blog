package repositories

import (
	"blog/config"
	"blog/db/models/interfaces"
	"blog/entities"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type BlogRepoModels struct {
	blog       interfaces.BlogsModel
	blogTags   interfaces.BlogTagsModel
	blogTopics interfaces.BlogTopicsModel
	tags       interfaces.TagsModel
	topics     interfaces.TopicsModel
}

func NewBlogsRepoModels(
	blog interfaces.BlogsModel,
	blogTags interfaces.BlogTagsModel,
	blogTopics interfaces.BlogTopicsModel,
	tags interfaces.TagsModel,
	topics interfaces.TopicsModel,
) *BlogRepoModels {

	return &BlogRepoModels{
		blog:       blog,
		blogTags:   blogTags,
		blogTopics: blogTopics,
		tags:       tags,
		topics:     topics,
	}
}

type Blogs struct {
	db     *sql.DB
	config config.DBSetting
	models BlogRepoModels
}

func NewBlogs(db *sql.DB, config config.DBSetting, models BlogRepoModels) *Blogs {
	return &Blogs{
		db:     db,
		config: config,
		models: models,
	}
}

func (b *Blogs) Create(ctx context.Context, blog entities.InBlog) (*entities.OutBlog, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	tx, err := b.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("Create: begin transaction error: %w", err)
	}

	newBlog, err := b.models.blog.Create(ctxTimeout, tx, blog)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.OutBlog{}, fmt.Errorf("Create: query rollback error: %w", err)
		}
		return &entities.OutBlog{}, fmt.Errorf("Create: models create blog failed: %w", err)
	}

	if err := b.models.blogTags.Upsert(ctxTimeout, tx, newBlog.ID, blog.Tags); err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.OutBlog{}, fmt.Errorf("Create: model create blog_tags rollback error: %w", err)
		}
		return &entities.OutBlog{}, fmt.Errorf("Create: model create blog_tags error: %w", err)
	}

	if err := b.models.blogTopics.Upsert(ctxTimeout, tx, newBlog.ID, blog.Topics); err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.OutBlog{}, fmt.Errorf("Create: model create blog_topics rollback error: %w", err)
		}
		return &entities.OutBlog{}, fmt.Errorf("Create: model create blog_topics error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return &entities.OutBlog{}, fmt.Errorf("Create: commit error: %w", err)
	}

	outBlog, err := b.fillOutBlog(ctxTimeout, *newBlog)
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("Create: fill OutBlog failed: %w", err)
	}

	return outBlog, nil
}

func (b *Blogs) Update(ctx context.Context, blog entities.InBlog, id int) (*entities.OutBlog, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	tx, err := b.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("Update: begin transaction error: %w", err)
	}

	// Update blog
	newBlog, err := b.models.blog.Update(ctxTimeout, tx, blog, id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.OutBlog{}, fmt.Errorf("Update: query rollback error: %w", err)
		}
		return &entities.OutBlog{}, fmt.Errorf("Update: models update blog failed: %w", err)
	}

	// Update many-to-many table
	// Uses upsert + inverse delete
	if err := b.models.blogTags.Upsert(ctxTimeout, tx, newBlog.ID, blog.Tags); err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.OutBlog{}, fmt.Errorf("Update: model update blog_tags rollback error: %w", err)
		}
		return &entities.OutBlog{}, fmt.Errorf("Update: model update blog_tags error: %w", err)
	}

	if err := b.models.blogTags.InverseDelete(ctxTimeout, tx, newBlog.ID, blog.Tags); err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.OutBlog{}, fmt.Errorf("Update: model inverse delete blog_tags rollback error: %w", err)
		}
		return &entities.OutBlog{}, fmt.Errorf("Update: model inverse delete blog_tags error: %w", err)
	}

	if err := b.models.blogTopics.Upsert(ctxTimeout, tx, newBlog.ID, blog.Topics); err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.OutBlog{}, fmt.Errorf("Update: model update blog_topics rollback error: %w", err)
		}
		return &entities.OutBlog{}, fmt.Errorf("Update: model update blog_topics error: %w", err)
	}

	if err := b.models.blogTopics.InverseDelete(ctxTimeout, tx, newBlog.ID, blog.Topics); err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.OutBlog{}, fmt.Errorf("Update: model inverse delete blog_topics rollback error: %w", err)
		}
		return &entities.OutBlog{}, fmt.Errorf("Update: model inverse delete blog_topics error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return &entities.OutBlog{}, fmt.Errorf("Update: commit error: %w", err)
	}

	outBlog, err := b.fillOutBlog(ctxTimeout, *newBlog)
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("Update: fill OutBlog failed: %w", err)
	}

	return outBlog, nil
}

/*
Only return blog with field values:

- visible: true

- deleted_at: ""
*/
func (b *Blogs) Get(ctx context.Context, id int) (*entities.OutBlog, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	blog, err := b.models.blog.Get(ctxTimeout, b.db, id)
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("Get: model get blog failed: %w", err)
	}

	outBlog, err := b.fillOutBlog(ctxTimeout, *blog)
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("Get: fill OutBlog failed: %w", err)
	}

	return outBlog, nil
}

/*
Only return blogs with field values:

- visible: true

- deleted_at: ""
*/
func (b *Blogs) List(ctx context.Context) ([]entities.OutBlog, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	blogs, err := b.models.blog.List(ctxTimeout, b.db)
	if err != nil {
		return []entities.OutBlog{}, fmt.Errorf("List: model list blogs failed: %w", err)
	}

	result := []entities.OutBlog{}

	for _, blog := range blogs {
		outBlog, err := b.fillOutBlog(ctxTimeout, blog)
		if err != nil {
			return []entities.OutBlog{}, fmt.Errorf("List: fill OutBlog failed: %w", err)
		}
		result = append(result, *outBlog)
	}

	return result, nil
}

/*
Only return blogs with field values:

- visible: true

- deleted_at: ""
*/
func (b *Blogs) ListByTopicIDs(ctx context.Context, topicID []int) ([]entities.OutBlog, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	blogs, err := b.models.blog.ListByTopicIDs(ctxTimeout, b.db, topicID)
	if err != nil {
		return []entities.OutBlog{}, fmt.Errorf("ListByTopicIDs: model list blogs by topic id failed: %w", err)
	}

	result := []entities.OutBlog{}

	for _, blog := range blogs {
		outBlog, err := b.fillOutBlog(ctxTimeout, blog)
		if err != nil {
			return []entities.OutBlog{}, fmt.Errorf("ListByTopicIDs: fill OutBlog failed: %w", err)
		}
		result = append(result, *outBlog)
	}

	return result, nil
}

/*
Only return blogs with field values:

- visible: true

- deleted_at: ""
*/
func (b *Blogs) ListByTopicAndTagIDs(ctx context.Context, topicID, tagID []int) ([]entities.OutBlog, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	blogs, err := b.models.blog.ListByTopicAndTagIDs(ctxTimeout, b.db, topicID, tagID)
	if err != nil {
		return []entities.OutBlog{}, fmt.Errorf("ListByTopicAndTagIDs: model list blogs by topic and tag ids failed: %w", err)
	}

	result := []entities.OutBlog{}

	for _, blog := range blogs {
		outBlog, err := b.fillOutBlog(ctxTimeout, blog)
		if err != nil {
			return []entities.OutBlog{}, fmt.Errorf("ListByTopicAndTagIDs: fill OutBlog failed: %w", err)
		}
		result = append(result, *outBlog)
	}

	return result, nil
}

// Get any blog regardless of visiblity and delete timestamp
func (b *Blogs) AdminGet(ctx context.Context, id int) (*entities.OutBlog, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	blog, err := b.models.blog.AdminGet(ctxTimeout, b.db, id)
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("AdminGet: model admin get blog failed: %w", err)
	}

	outBlog, err := b.fillOutBlog(ctxTimeout, *blog)
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("AdminGet: fill OutBlog failed: %w", err)
	}

	return outBlog, nil
}

// Returns all blogs
func (b *Blogs) AdminList(ctx context.Context) ([]entities.OutBlog, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	blogs, err := b.models.blog.AdminList(ctxTimeout, b.db)
	if err != nil {
		return []entities.OutBlog{}, fmt.Errorf("AdminList: model list blogs failed: %w", err)
	}

	result := []entities.OutBlog{}

	for _, blog := range blogs {
		outBlog, err := b.fillOutBlog(ctxTimeout, blog)
		if err != nil {
			return []entities.OutBlog{}, fmt.Errorf("AdminList: fill OutBlog failed: %w", err)
		}
		result = append(result, *outBlog)
	}

	return result, nil
}

// Returns all matched blogs
func (b *Blogs) AdminListByTopicIDs(ctx context.Context, topicID []int) ([]entities.OutBlog, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	blogs, err := b.models.blog.AdminListByTopicIDs(ctxTimeout, b.db, topicID)
	if err != nil {
		return []entities.OutBlog{}, fmt.Errorf("AdminListByTopicIDs: model list blogs by topic id failed: %w", err)
	}

	result := []entities.OutBlog{}

	for _, blog := range blogs {
		outBlog, err := b.fillOutBlog(ctxTimeout, blog)
		if err != nil {
			return []entities.OutBlog{}, fmt.Errorf("AdminListByTopicIDs: fill OutBlog failed: %w", err)
		}
		result = append(result, *outBlog)
	}

	return result, nil
}

// Returns all matched blogs
func (b *Blogs) AdminListByTopicAndTagIDs(ctx context.Context, topicID, tagID []int) ([]entities.OutBlog, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	blogs, err := b.models.blog.AdminListByTopicAndTagIDs(ctxTimeout, b.db, topicID, tagID)
	if err != nil {
		return []entities.OutBlog{}, fmt.Errorf("AdminListByTopicAndTagIDs: model list blogs by topic and tag ids failed: %w", err)
	}

	result := []entities.OutBlog{}

	for _, blog := range blogs {
		outBlog, err := b.fillOutBlog(ctxTimeout, blog)
		if err != nil {
			return []entities.OutBlog{}, fmt.Errorf("AdminListByTopicAndTagIDs: fill OutBlog failed: %w", err)
		}
		result = append(result, *outBlog)
	}

	return result, nil
}

func (b *Blogs) SoftDelete(ctx context.Context, id int) (int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	tx, err := b.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("SoftDelete: begin transaction failed: %w", err)
	}

	affectedRows, err := b.models.blog.SoftDelete(ctxTimeout, tx, id)
	if err != nil {
		return 0, fmt.Errorf("SoftDelete: blog soft delete failed: %w", err)
	}

	return affectedRows, nil
}

func (b *Blogs) Delele(ctx context.Context, id int) (int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	tx, err := b.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("Delete: begin transaction failed: %w", err)
	}

	affectedRows, err := b.models.blog.Delete(ctxTimeout, tx, id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, fmt.Errorf("Delete: model delete blog rollback failed: %w", err)
		}
		return 0, fmt.Errorf("Delete: model delete blog failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("Delete: commit failed: %w", err)
	}

	return affectedRows, nil
}

func (b *Blogs) RestoreDeleted(ctx context.Context, id int) (*entities.OutBlog, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(b.config.Timeout)*time.Second)
	defer cancel()

	tx, err := b.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("RestoreDeleted: begin transaction failed: %w", err)
	}

	blog, err := b.models.blog.RestoreDeleted(ctxTimeout, tx, id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.OutBlog{}, fmt.Errorf("RestoreDeleted: model restore deleted blog rollback failed: %w", err)
		}
		return &entities.OutBlog{}, fmt.Errorf("RestoreDeleted: model restore deleted blog failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return &entities.OutBlog{}, fmt.Errorf("RestoreDeleted: commit failed: %w", err)
	}

	outBlog, err := b.fillOutBlog(ctxTimeout, *blog)
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("RestoreDeleted: fill OutBlog failed: %w", err)
	}

	return outBlog, nil
}

// Helper function to fill out OutBlog with tags and topics
func (b *Blogs) fillOutBlog(ctx context.Context, blog entities.Blog) (*entities.OutBlog, error) {
	tags, err := b.models.tags.GetByBlogID(ctx, b.db, blog.ID)
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("fillOutBlog: model get tags failed: %w", err)
	}

	topics, err := b.models.topics.GetByBlogID(ctx, b.db, blog.ID)
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("fillOutBlog: model get topics failed: %w", err)
	}

	outBlog := entities.NewOutBlog(blog, tags, topics)
	return outBlog, nil
}
