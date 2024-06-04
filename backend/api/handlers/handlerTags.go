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

type tagsRepository interface {
	CreateTag(ctx context.Context, blog entities.Tag) (*entities.OutBlog, error)
	ListTags(ctx context.Context) ([]entities.Tag, error)
	GetTag(ctx context.Context, id int) (*entities.Tag, error)
}

type Tags struct {
	repo tagsRepository
}

func NewTags(repo tagsRepository) *Tags {
	return &Tags{
		repo: repo,
	}
}

func (t *Tags) CreateTag(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("CreateTag")

	body := &entities.Tag{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		slog.Error("CreateTag: decode failed", "error", err.Error())
		return writeJSON(w, err, nil, http.StatusBadRequest)
	}
	inTag := entities.NewTag(
		body.Name,
		body.Description,
	)

	outTag, err := t.repo.CreateTag(r.Context(), *inTag)
	if err != nil {
		slog.Error("CreateTag: create tag failed", "error", err.Error())
		return writeJSON(w, err, nil, http.StatusInternalServerError)
	}

	return writeJSON(w, nil, outTag, http.StatusOK)
}

func (t *Tags) ListTags(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("ListTags")

	tags, err := t.repo.ListTags(r.Context())
	if err != nil {
		slog.Error("ListTags: list tags failed", "error", err)
		return writeJSON(w, err, nil, http.StatusInternalServerError)
	}

	return writeJSON(w, nil, tags, http.StatusOK)
}

func (t *Tags) GetTag(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("GetTag")

	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return writeJSON(w, err, nil, http.StatusBadRequest)
	}

	tag, err := t.repo.GetTag(r.Context(), id)
	if err != nil {
		// differentiate if it's db error or that the user supplied id dosen't exist
		slog.Error("GetTag: get tag failed", "error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return writeJSON(w, err, nil, http.StatusBadRequest)
		} else {
			return writeJSON(w, err, nil, http.StatusInternalServerError)
		}
	}

	return writeJSON(w, nil, tag, http.StatusOK)
}

func (t *Tags) UpdateTag(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (t *Tags) DeleteTag(w http.ResponseWriter, r *http.Request) error {
	return nil
}
