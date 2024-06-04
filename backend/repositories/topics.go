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
		return &entities.Topic{}, fmt.Errorf("Create: begin transaction error: %w", err)
	}

	newTopic, err := t.models.topics.Create(ctxTimeout, tx, topic)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return &entities.Topic{}, fmt.Errorf("Create: model create topic rollback error: %w", err)
		}
		return &entities.Topic{}, fmt.Errorf("Create: model create topic failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return &entities.Topic{}, fmt.Errorf("RepoCreateTopic: commit error: %w", err)
	}

	return newTopic, nil
}

func (t *Topics) List(ctx context.Context) ([]entities.Topic, error) {
	return []entities.Topic{}, nil
}

func (t *Topics) Get(ctx context.Context, id int) (*entities.Topic, error) {
	return &entities.Topic{}, nil
}

func (t *Topics) Update(ctx context.Context, topic entities.Topic) (*entities.Topic, error) {
	return &entities.Topic{}, nil
}

func (t *Topics) Delete(ctx context.Context, id int) error {
	return nil
}
