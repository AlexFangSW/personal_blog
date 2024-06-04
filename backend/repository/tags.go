package repository

import (
	"blog/config"
	"blog/db/models/interfaces"
	"blog/entities"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type TagsRepoModels struct {
	blogTags interfaces.BlogTagsModel
	tags     interfaces.TagsModel
}

func NewTagsRepoModels(
	blogTags interfaces.BlogTagsModel,
	tags interfaces.TagsModel,
) *TagsRepoModels {

	return &TagsRepoModels{
		blogTags: blogTags,
		tags:     tags,
	}
}

type Tags struct {
	db     *sql.DB
	config config.DBSetting
	models TagsRepoModels
}

func NewTags(db *sql.DB, config config.DBSetting, models TagsRepoModels) *Tags {
	return &Tags{
		db:     db,
		config: config,
		models: models,
	}
}

func (t *Tags) CreateTag(ctx context.Context, tag entities.Tag) (*entities.Tag, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tx, err := t.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return &entities.Tag{}, fmt.Errorf("RepoCreateTag: begin transaction error: %w", err)
	}

	newTag, err := t.models.tags.Create(ctxTimeout, tx, tag)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.Tag{}, fmt.Errorf("RepoCreateTag: rollback error: %w", err)
		}
		return &entities.Tag{}, fmt.Errorf("RepoCreateTag: insert tag failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return &entities.Tag{}, fmt.Errorf("RepoCreateTag: commit error: %w", err)
	}

	return newTag, nil
}

func (t *Tags) GetTagsByBlogID(ctx context.Context, blog_id int) ([]entities.Tag, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tags, err := t.models.tags.GetByBlogID(ctxTimeout, t.db, blog_id)
	if err != nil {
		return []entities.Tag{}, fmt.Errorf("RepoGetTagsByBlogID: query context failed: %w", err)
	}

	return tags, nil
}

func (t *Tags) ListTags(ctx context.Context) ([]entities.Tag, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tags, err := t.models.tags.List(ctxTimeout, t.db)
	if err != nil {
		return []entities.Tag{}, fmt.Errorf("RepoListTags: query failed: %w", err)
	}

	return tags, nil
}

func (t *Tags) GetTag(ctx context.Context, id int) (*entities.Tag, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tag, err := t.models.tags.Get(ctxTimeout, t.db, id)
	if err != nil {
		return &entities.Tag{}, fmt.Errorf("RepoGetTag: query failed: %w", err)
	}

	return tag, nil
}
