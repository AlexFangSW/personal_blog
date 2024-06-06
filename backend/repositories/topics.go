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

type TopicsRepoModels struct {
	blogTopics interfaces.BlogTopicsModel
	topics     interfaces.TopicsModel
}

func NewTopicsRepoModels(
	blogTopics interfaces.BlogTopicsModel,
	topics interfaces.TopicsModel,
) *TopicsRepoModels {

	return &TopicsRepoModels{
		blogTopics: blogTopics,
		topics:     topics,
	}
}

type Topics struct {
	db     *sql.DB
	config config.DBSetting
	models TopicsRepoModels
}

func NewTopics(db *sql.DB, config config.DBSetting, models TopicsRepoModels) *Topics {
	return &Topics{
		db:     db,
		config: config,
		models: models,
	}
}

func (t *Topics) Create(ctx context.Context, topic entities.Topic) (*entities.Topic, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tx, err := t.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return &entities.Topic{}, fmt.Errorf("Create: begin transaction failed: %w", err)
	}

	newTopic, err := t.models.topics.Create(ctxTimeout, tx, topic)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.Topic{}, fmt.Errorf("Create: model create topic rollback failed: %w", err)
		}
		return &entities.Topic{}, fmt.Errorf("Create: model create topic failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return &entities.Topic{}, fmt.Errorf("Create: commit failed: %w", err)
	}

	return newTopic, nil
}

func (t *Topics) List(ctx context.Context) ([]entities.Topic, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	topics, err := t.models.topics.List(ctxTimeout, t.db)
	if err != nil {
		return []entities.Topic{}, fmt.Errorf("List: model list topics failed: %w", err)
	}

	return topics, nil
}

func (t *Topics) Get(ctx context.Context, id int) (*entities.Topic, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	topic, err := t.models.topics.Get(ctxTimeout, t.db, id)
	if err != nil {
		return &entities.Topic{}, fmt.Errorf("Get: model get topic failed: %w", err)
	}

	return topic, nil
}

func (t *Topics) Update(ctx context.Context, topic entities.Topic) (*entities.Topic, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tx, err := t.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return &entities.Topic{}, fmt.Errorf("Update: begin transaction failed: %w", err)
	}

	newTopic, err := t.models.topics.Update(ctxTimeout, tx, topic)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.Topic{}, fmt.Errorf("Update: model update topic rollback failed: %w", err)
		}
		return &entities.Topic{}, fmt.Errorf("Update: model update topic failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return &entities.Topic{}, fmt.Errorf("Update: commit failed: %w", err)
	}

	return newTopic, nil
}

func (t *Topics) Delete(ctx context.Context, id int) (int, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Duration(t.config.Timeout)*time.Second)
	defer cancel()

	tx, err := t.db.BeginTx(ctxTimeout, &sql.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("Delete: begin transaction failed: %w", err)
	}

	affectedRows, err := t.models.topics.Delete(ctxTimeout, tx, id)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return 0, fmt.Errorf("Delete: model delete topic rollback failed: %w", err)
		}
		return 0, fmt.Errorf("Delete: model delete topic failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("Delete: commit failed: %w", err)
	}

	return affectedRows, nil
}
