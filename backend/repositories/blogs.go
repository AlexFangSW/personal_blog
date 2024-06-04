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

	if err := b.models.blogTags.Create(ctxTimeout, tx, newBlog.ID, blog.Tags); err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.OutBlog{}, fmt.Errorf("Create: model create blog_tags rollback error: %w", err)
		}
		return &entities.OutBlog{}, fmt.Errorf("Create: model create blog_tags error: %w", err)
	}

	if err := b.models.blogTopics.Create(ctxTimeout, tx, newBlog.ID, blog.Topics); err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.OutBlog{}, fmt.Errorf("Create: model create blog_topics rollback error: %w", err)
		}
		return &entities.OutBlog{}, fmt.Errorf("Create: model create blog_topics error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return &entities.OutBlog{}, fmt.Errorf("Create: commit error: %w", err)
	}

	tags, err := b.models.tags.GetByBlogID(ctxTimeout, b.db, newBlog.ID)
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("Create: model get tags by blog id error: %w", err)
	}
	topics, err := b.models.topics.GetByBlogID(ctxTimeout, b.db, newBlog.ID)
	if err != nil {
		return &entities.OutBlog{}, fmt.Errorf("Create: model get topics by blog id error: %w", err)
	}

	outBlog := entities.NewOutBlog(*newBlog, tags, topics)
	return outBlog, nil
}

func (b *Blogs) Update(ctx context.Context, blog entities.InBlog) (*entities.OutBlog, error) {
	return &entities.OutBlog{}, nil
}

func (b *Blogs) Get(ctx context.Context, id int) (*entities.OutBlog, error) {
	return &entities.OutBlog{}, nil
}
func (b *Blogs) AdminGet(ctx context.Context, id int) (*entities.OutBlog, error) {

	return &entities.OutBlog{}, nil
}

func (b *Blogs) List(ctx context.Context) ([]entities.OutBlog, error) {

	return []entities.OutBlog{}, nil
}
func (b *Blogs) ListByTopicID(ctx context.Context, topicID int) ([]entities.OutBlog, error) {

	return []entities.OutBlog{}, nil
}

func (b *Blogs) AdminList(ctx context.Context) ([]entities.OutBlog, error) {

	return []entities.OutBlog{}, nil
}

func (b *Blogs) AdminListByTopicID(ctx context.Context, topicID int) ([]entities.OutBlog, error) {

	return []entities.OutBlog{}, nil
}

func (b *Blogs) SoftDelele(ctx context.Context, id int) error {
	return nil
}
func (b *Blogs) Delele(ctx context.Context, id int) error {
	return nil
}
func (b *Blogs) RestoreDeleted(ctx context.Context, id int) (*entities.OutBlog, error) {
	return &entities.OutBlog{}, nil
}
