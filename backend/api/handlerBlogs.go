package api

import (
	"blog/entities"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *Server) CreateBlog(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("CreateTag")

	body := &entities.InBlog{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		slog.Error("CreateBlog: decode failed", "error", err.Error())
		return writeJSON(w, err, nil, http.StatusBadRequest)
	}
	blog := entities.NewBlog(
		body.Title,
		body.Content,
		body.Description,
		body.Pined,
		body.Visible,
	)
	inBlog := entities.NewInBlog(
		*blog,
		body.Tags,
		body.Topics,
	)

	outBlog, err := s.models.CreateBlog(r.Context(), *inBlog)
	if err != nil {
		slog.Error("CreateBlog: create blog failed", "error", err.Error())
		return writeJSON(w, err, nil, http.StatusInternalServerError)
	}

	return writeJSON(w, nil, outBlog, http.StatusOK)
}

func (s *Server) ListBlogs(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) GetBlog(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) UpdateBlog(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) SoftDeleteBlog(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) ListDeletedBlogs(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) GetDeletedBlog(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) RestoreDeletedBlog(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Server) DeleteBlog(w http.ResponseWriter, r *http.Request) error {
	return nil
}
