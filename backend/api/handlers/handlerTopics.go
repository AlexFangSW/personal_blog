package handlers

import (
	"blog/entities"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

// Concrete implementations are at repository/<name>
type topicsRepository interface {
	Create(ctx context.Context, topic entities.Topic) (*entities.Topic, error)
	List(ctx context.Context) ([]entities.Topic, error)
	Get(ctx context.Context, id int) (*entities.Topic, error)
	Update(ctx context.Context, topic entities.Topic) (*entities.Topic, error)
	Delete(ctx context.Context, id int) (int, error)
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
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}
	inTopic := entities.NewTopic(
		body.Name,
		body.Description,
	)

	outTopic, err := t.repo.Create(r.Context(), *inTopic)
	if err != nil {
		slog.Error("CreateTopic: repo create failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	return entities.NewRetSuccess[entities.Topic](*outTopic).WriteJSON(w)
}

func (t *Topics) ListTopics(w http.ResponseWriter, r *http.Request) error {
	slog.Info("ListTopics")

	topics, err := t.repo.List(r.Context())
	if err != nil {
		slog.Error("ListTopics: repo list failed", "error", err)
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	return entities.NewRetSuccess[[]entities.Topic](topics).WriteJSON(w)
}

func (t *Topics) GetTopic(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("GetTopic")

	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}

	topic, err := t.repo.Get(r.Context(), id)
	if err != nil {
		// differentiate if it's db error or that the user supplied id dosen't exist
		slog.Error("GetTopic: repo get failed", "error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
		} else {
			return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
		}
	}

	return entities.NewRetSuccess[entities.Topic](*topic).WriteJSON(w)
}

func (t *Topics) UpdateTopic(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("UpdateTopic")

	// load body
	body := &entities.Topic{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		slog.Error("UpdateTopic: decode failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}
	inTopic := entities.NewTopic(
		body.Name,
		body.Description,
	)

	// get target id
	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		slog.Error("UpdateTopic: id string to int failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}
	inTopic.ID = id

	outTopic, err := t.repo.Update(r.Context(), *inTopic)
	if err != nil {
		// differentiate if it's db error or that the user supplied id dosen't exist
		slog.Error("UpdateTopic: repo update failed", "error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
		} else {
			return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
		}
	}

	return entities.NewRetSuccess[entities.Topic](*outTopic).WriteJSON(w)
}

func (t *Topics) DeleteTopic(w http.ResponseWriter, r *http.Request) error {
	slog.Info("DeleteTopic")

	// get target id
	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		slog.Error("DeleteTopic: id string to int failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}

	affectedRows, err := t.repo.Delete(r.Context(), id)
	if err != nil {
		slog.Error("DeleteTopic: repo delete failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	if affectedRows == 0 {
		return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
	}
	return entities.NewRetSuccess[entities.RowsAffected](*entities.NewRowsAffected(affectedRows)).WriteJSON(w)
}
