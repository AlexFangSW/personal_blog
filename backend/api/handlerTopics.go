package api

import (
	"blog/db/models"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *Server) CreateTopic(w http.ResponseWriter, r *http.Request) error {
	body := &models.Topic{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		slog.Error("CreateTopic: decode failed", "error", err.Error())
		return writeJSON(w, err, "", http.StatusBadRequest)
	}
	inTopic := models.NewTopic(
		body.Name,
		body.Description,
	)

	outTopic, err := s.models.CreateTopic(r.Context(), *inTopic)
	if err != nil {
		slog.Error("CreateTopic: create topic failed", "error", err.Error())
		return writeJSON(w, err, "", http.StatusInternalServerError)
	}

	return writeJSON(w, nil, outTopic, http.StatusOK)
}

func (s *Server) ListTopics(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) GetTopic(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) UpdateTopic(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) DeleteTopic(w http.ResponseWriter, r *http.Request) error {
	return nil
}
