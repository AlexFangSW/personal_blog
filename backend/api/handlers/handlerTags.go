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
type tagsRepository interface {
	Create(ctx context.Context, tag entities.Tag) (*entities.Tag, error)
	List(ctx context.Context) ([]entities.Tag, error)
	Get(ctx context.Context, id int) (*entities.Tag, error)
	Update(ctx context.Context, tag entities.Tag, id int) (*entities.Tag, error)
	Delete(ctx context.Context, id int) (int, error)
}

type Tags struct {
	repo tagsRepository
}

func NewTags(repo tagsRepository) *Tags {
	return &Tags{
		repo: repo,
	}
}

// CreateTag
//
//	@Summary		Create tag
//	@Description	tags must have unique names
//	@Tags			tags
//	@Accept			json
//	@Produce		json
//	@Param			tag	body		entities.InTag	true	"new tag contents"
//	@Success		200	{object}	entities.RetSuccess[entities.Tag]
//	@Failure		400	{object}	entities.RetFailed
//	@Failure		500	{object}	entities.RetFailed
//	@Router			/tags [post]
func (t *Tags) CreateTag(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("CreateTag")

	body := &entities.InTag{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		slog.Error("CreateTag: decode failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}
	inTag := entities.NewTag(
		body.Name,
		body.Description,
	)

	outTag, err := t.repo.Create(r.Context(), *inTag)
	if err != nil {
		slog.Error("CreateTag: repo create failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	return entities.NewRetSuccess[entities.Tag](*outTag).WriteJSON(w)
}

// ListTags
//
//	@Summary		List tags
//	@Description	list all tags
//	@Tags			tags
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	entities.RetSuccess[[]entities.Tag]
//	@Failure		500	{object}	entities.RetFailed
//	@Router			/tags [get]
func (t *Tags) ListTags(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("ListTags")

	tags, err := t.repo.List(r.Context())
	if err != nil {
		slog.Error("ListTags: repo list failed", "error", err)
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	return entities.NewRetSuccess[[]entities.Tag](tags).WriteJSON(w)
}

// GetTag
//
//	@Summary		Get tags
//	@Description	get tag by id
//	@Tags			tags
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"target tag id"
//	@Success		200	{object}	entities.RetSuccess[entities.Tag]
//	@Failure		400	{object}	entities.RetFailed
//	@Failure		404	{object}	entities.RetFailed
//	@Failure		500	{object}	entities.RetFailed
//	@Router			/tags/{id} [get]
func (t *Tags) GetTag(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("GetTag")

	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}

	tag, err := t.repo.Get(r.Context(), id)
	if err != nil {
		// differentiate if it's db error or that the user supplied id dosen't exist
		slog.Error("GetTag: repo get failed", "error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
		} else {
			return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
		}
	}

	return entities.NewRetSuccess[entities.Tag](*tag).WriteJSON(w)
}

// UpdateTag
//
//	@Summary		Update tag
//	@Description	update tag
//	@Tags			tags
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int				true	"target tag id"
//	@Param			tag	body		entities.InTag	true	"new tag content"
//	@Success		200	{object}	entities.RetSuccess[entities.Tag]
//	@Failure		400	{object}	entities.RetFailed
//	@Failure		404	{object}	entities.RetFailed
//	@Failure		500	{object}	entities.RetFailed
//	@Router			/tags/{id} [put]
func (t *Tags) UpdateTag(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("UpdateTag")

	// load body
	body := &entities.InTag{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		slog.Error("UpdateTag: decode failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}
	inTag := entities.NewTag(
		body.Name,
		body.Description,
	)

	// get target id
	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		slog.Error("UpdateTag: id string to int failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}

	outTag, err := t.repo.Update(r.Context(), *inTag, id)
	if err != nil {
		slog.Error("UpdateTag: repo update failed", "error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
		} else {
			return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
		}
	}

	return entities.NewRetSuccess[entities.Tag](*outTag).WriteJSON(w)
}

// DeleteTag
//
//	@Summary		Delete tag
//	@Description	delete tag
//	@Tags			tags
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"target tag id"
//	@Success		200	{object}	entities.RetSuccess[entities.RowsAffected]
//	@Failure		400	{object}	entities.RetFailed
//	@Failure		404	{object}	entities.RetFailed
//	@Failure		500	{object}	entities.RetFailed
//	@Router			/tags/{id} [delete]
func (t *Tags) DeleteTag(w http.ResponseWriter, r *http.Request) error {
	slog.Info("DeleteTag")

	// get target id
	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		slog.Error("DeleteTag: id string to int failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}

	affectedRows, err := t.repo.Delete(r.Context(), id)
	if err != nil {
		slog.Error("DeleteTag: repo delete failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	if affectedRows == 0 {
		return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
	}

	return entities.NewRetSuccess[entities.RowsAffected](*entities.NewRowsAffected(affectedRows)).WriteJSON(w)
}
