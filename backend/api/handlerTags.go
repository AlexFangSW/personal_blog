package api

import (
	"blog/db/models"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *Server) CreateTag(w http.ResponseWriter, r *http.Request) error {
	body := &models.Tag{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		slog.Error("CreateTag: decode error", "error", err.Error())
		return writeJSON(w, err, "", http.StatusBadRequest)
	}
	inTag := models.NewTag(
		body.Name,
		body.Description,
	)

	ctx := context.Background()
	outTag, err := s.models.CreateTag(ctx, *inTag)
	if err != nil {
		slog.Error("CreateTag: create tag error", "error", err.Error())
		return writeJSON(w, err, "", http.StatusInternalServerError)
	}

	return writeJSON(w, nil, outTag, http.StatusOK)
}

func (s *Server) ListTag(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) GetTag(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) UpdateTag(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) DeleteTag(w http.ResponseWriter, r *http.Request) error {
	return nil
}
