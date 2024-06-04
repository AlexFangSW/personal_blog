package handlers

import (
	"blog/entities"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type topicsRepository interface {
	CreateTopic(ctx context.Context, topic entities.Topic) (*entities.Topic, error)
	GetTopicsByBlogID(ctx context.Context, blog_id int) ([]entities.Topic, error)
}

type Topics struct {
	repo topicsRepository
}

func NewTopics(repo topicsRepository) *Topics {
	return &Topics{
		repo: repo,
	}
}

func (t *Topics) CreateTopic(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("CreateTopic")

	body := &entities.Topic{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		slog.Error("CreateTopic: decode failed", "error", err.Error())
		return writeJSON(w, err, nil, http.StatusBadRequest)
	}
	inTopic := entities.NewTopic(
		body.Name,
		body.Description,
	)

	outTopic, err := t.repo.CreateTopic(r.Context(), *inTopic)
	if err != nil {
		slog.Error("CreateTopic: create topic failed", "error", err.Error())
		return writeJSON(w, err, nil, http.StatusInternalServerError)
	}

	return writeJSON(w, nil, outTopic, http.StatusOK)
}

func (t *Topics) ListTopics(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (t *Topics) GetTopic(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (t *Topics) UpdateTopic(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (t *Topics) DeleteTopic(w http.ResponseWriter, r *http.Request) error {
	return nil
}
