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
	Update(ctx context.Context, topic entities.Topic, id int) (*entities.Topic, error)
	Delete(ctx context.Context, id int) (int, error)
}

type Topics struct {
	repo topicsRepository
	auth authHelper
}

func NewTopics(repo topicsRepository, auth authHelper) *Topics {
	return &Topics{
		repo: repo,
		auth: auth,
	}
}

// CreateTopic
//
//	@Summary		Create topic
//	@Description	topics must have unique names
//	@Tags			topics
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"jwt token"
//	@Param			topic			body		entities.InTopic	true	"new topic contents"
//	@Success		200				{object}	entities.RetSuccess[entities.Topic]
//	@Failure		400				{object}	entities.RetFailed
//	@Failure		500				{object}	entities.RetFailed
//	@Router			/topics [post]
func (t *Topics) CreateTopic(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("CreateTopic")

	// authorization
	authorized, err := t.auth.Verify(r)
	if err != nil || !authorized {
		slog.Warn("CreateTopic: authorization failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusForbidden).WriteJSON(w)
	}

	body := &entities.InTopic{}
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

		if sqliteErr, ok := getSQLiteError(err); ok {
			slog.Error("got sqlite error", "error code", sqliteErr.Code, "extended error code", sqliteErr.ExtendedCode)
			return entities.NewRetFailedCustom(err, int(sqliteErr.ExtendedCode), http.StatusInternalServerError).WriteJSON(w)
		}

		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	return entities.NewRetSuccess(*outTopic).WriteJSON(w)
}

// ListTopics
//
//	@Summary		List topics
//	@Description	list all topics
//	@Tags			topics
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entities.RetSuccess[[]entities.Topic]
//	@Failure		500	{object}	entities.RetFailed
//	@Router			/topics [get]
func (t *Topics) ListTopics(w http.ResponseWriter, r *http.Request) error {
	slog.Info("ListTopics")

	topics, err := t.repo.List(r.Context())
	if err != nil {
		slog.Error("ListTopics: repo list failed", "error", err)
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	return entities.NewRetSuccess(topics).WriteJSON(w)
}

// GetTopic
//
//	@Summary		Get topic
//	@Description	get topic by id
//	@Tags			topics
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"target topic id"
//	@Success		200	{object}	entities.RetSuccess[entities.Topic]
//	@Failure		400	{object}	entities.RetFailed
//	@Failure		404	{object}	entities.RetFailed
//	@Failure		500	{object}	entities.RetFailed
//	@Router			/topics/{id} [get]
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
		}

		if sqliteErr, ok := getSQLiteError(err); ok {
			slog.Error("got sqlite error", "error code", sqliteErr.Code, "extended error code", sqliteErr.ExtendedCode)
			return entities.NewRetFailedCustom(err, int(sqliteErr.ExtendedCode), http.StatusInternalServerError).WriteJSON(w)
		}

		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	return entities.NewRetSuccess(*topic).WriteJSON(w)
}

// UpdateTopic
//
//	@Summary		Update topic
//	@Description	update topic
//	@Tags			topics
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int					true	"target tag id"
//	@Param			Authorization	header		string				true	"jwt token"
//	@Param			topic			body		entities.InTopic	true	"new topic content"
//	@Success		200				{object}	entities.RetSuccess[entities.Topic]
//	@Failure		400				{object}	entities.RetFailed
//	@Failure		404				{object}	entities.RetFailed
//	@Failure		500				{object}	entities.RetFailed
//	@Router			/topics/{id} [put]
func (t *Topics) UpdateTopic(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("UpdateTopic")

	// authorization
	authorized, err := t.auth.Verify(r)
	if err != nil || !authorized {
		slog.Warn("UpdateTopic: authorization failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusForbidden).WriteJSON(w)
	}

	// load body
	body := &entities.InTopic{}
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

	outTopic, err := t.repo.Update(r.Context(), *inTopic, id)
	if err != nil {
		// differentiate if it's db error or that the user supplied id dosen't exist
		slog.Error("UpdateTopic: repo update failed", "error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
		}

		if sqliteErr, ok := getSQLiteError(err); ok {
			slog.Error("got sqlite error", "error code", sqliteErr.Code, "extended error code", sqliteErr.ExtendedCode)
			return entities.NewRetFailedCustom(err, int(sqliteErr.ExtendedCode), http.StatusInternalServerError).WriteJSON(w)
		}

		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	return entities.NewRetSuccess(*outTopic).WriteJSON(w)
}

// DeleteTopic
//
//	@Summary		Delete topic
//	@Description	delete topic
//	@Tags			topics
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int		true	"target topic id"
//	@Param			Authorization	header		string	true	"jwt token"
//	@Success		200				{object}	entities.RetSuccess[entities.RowsAffected]
//	@Failure		400				{object}	entities.RetFailed
//	@Failure		404				{object}	entities.RetFailed
//	@Failure		500				{object}	entities.RetFailed
//	@Router			/topics/{id} [delete]
func (t *Topics) DeleteTopic(w http.ResponseWriter, r *http.Request) error {
	slog.Info("DeleteTopic")

	// authorization
	authorized, err := t.auth.Verify(r)
	if err != nil || !authorized {
		slog.Warn("DeleteTopic: authorization failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusForbidden).WriteJSON(w)
	}

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

		if sqliteErr, ok := getSQLiteError(err); ok {
			slog.Error("got sqlite error", "error code", sqliteErr.Code, "extended error code", sqliteErr.ExtendedCode)
			return entities.NewRetFailedCustom(err, int(sqliteErr.ExtendedCode), http.StatusInternalServerError).WriteJSON(w)
		}

		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	if affectedRows == 0 {
		return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
	}
	return entities.NewRetSuccess(*entities.NewRowsAffected(affectedRows)).WriteJSON(w)
}
