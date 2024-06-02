package api

import (
	"blog/db/models"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *Server) CreateBlog(w http.ResponseWriter, r *http.Request) error {
	body := &models.InBlog{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		slog.Error("CreateBlog: decode error", "error", err.Error())
		return writeJSON(w, err, "", http.StatusBadRequest)
	}
	blog := models.NewBlog(
		body.Title,
		body.Content,
		body.Description,
		body.Pined,
		body.Visible,
	)
	inBlog := models.NewInBlog(
		*blog,
		body.Tags,
		body.Topics,
	)

	ctx := context.Background()
	outBlog, err := s.models.CreateBlog(ctx, *inBlog)
	if err != nil {
		slog.Error("CreateBlog: create blog error", "error", err.Error())
		return writeJSON(w, err, "", http.StatusInternalServerError)
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
