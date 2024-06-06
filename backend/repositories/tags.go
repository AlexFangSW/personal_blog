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

func (t *Tags) Create(ctx context.Context, tag entities.Tag) (*entities.Tag, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tx, err := t.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return &entities.Tag{}, fmt.Errorf("Create: begin transaction failed: %w", err)
	}

	newTag, err := t.models.tags.Create(ctxTimeout, tx, tag)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.Tag{}, fmt.Errorf("Create: model create tag rollback failed: %w", err)
		}
		return &entities.Tag{}, fmt.Errorf("Create: model create tag failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return &entities.Tag{}, fmt.Errorf("Create: commit failed: %w", err)
	}

	return newTag, nil
}

func (t *Tags) List(ctx context.Context) ([]entities.Tag, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tags, err := t.models.tags.List(ctxTimeout, t.db)
	if err != nil {
		return []entities.Tag{}, fmt.Errorf("List: model list tags failed: %w", err)
	}

	return tags, nil
}

func (t *Tags) Get(ctx context.Context, id int) (*entities.Tag, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tag, err := t.models.tags.Get(ctxTimeout, t.db, id)
	if err != nil {
		return &entities.Tag{}, fmt.Errorf("Get: model get tags failed: %w", err)
	}

	return tag, nil
}

func (t *Tags) Update(ctx context.Context, tag entities.Tag) (*entities.Tag, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tx, err := t.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return &entities.Tag{}, fmt.Errorf("Update: begin transaction failed: %w", err)
	}

	newTag, err := t.models.tags.Update(ctxTimeout, tx, tag)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.Tag{}, fmt.Errorf("Update: model update tag rollback failed: %w", err)
		}
		return &entities.Tag{}, fmt.Errorf("Update: model update tag failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return &entities.Tag{}, fmt.Errorf("Update: commit failed: %w", err)
	}

	return newTag, nil
}

func (t *Tags) Delete(ctx context.Context, id int) (int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tx, err := t.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("Delete: begin transaction failed: %w", err)
	}

	affectedRows, err := t.models.tags.Delete(ctxTimeout, tx, id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, fmt.Errorf("Delete: model delete tag rollback failed: %w", err)
		}
		return 0, fmt.Errorf("Delete: model delete tag failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("Delete: commit failed: %w", err)
	}

	return affectedRows, nil
}
