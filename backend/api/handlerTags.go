package api

import (
	"blog/entities"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
)

func (s *Server) CreateTag(w http.ResponseWriter, r *http.Request) error {
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

	outTag, err := s.models.CreateTag(r.Context(), *inTag)
	if err != nil {
		slog.Error("CreateTag: create tag failed", "error", err.Error())
		return writeJSON(w, err, nil, http.StatusInternalServerError)
	}

	return writeJSON(w, nil, outTag, http.StatusOK)
}

func (s *Server) ListTags(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("ListTags")

	tags, err := s.models.ListTags(r.Context())
	if err != nil {
		slog.Error("ListTags: list tags failed", "error", err)
		return writeJSON(w, err, nil, http.StatusInternalServerError)
	}

	return writeJSON(w, nil, tags, http.StatusOK)
}

func (s *Server) GetTag(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("GetTag")

	rawID := r.PathValue("id")
	id, err := strconv.Atoi(rawID)
	if err != nil {
		return writeJSON(w, err, nil, http.StatusBadRequest)
	}

	tag, err := s.models.GetTag(r.Context(), id)
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

func (s *Server) UpdateTag(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) DeleteTag(w http.ResponseWriter, r *http.Request) error {
	return nil
}
