package handlers

import (
	"blog/entities"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
)

type blogsRepository interface {
	Create(ctx context.Context, blog entities.InBlog) (*entities.OutBlog, error)
}

type Blogs struct {
	repo blogsRepository
}

func NewBlogs(repo blogsRepository) *Blogs {
	return &Blogs{
		repo: repo,
	}
}

func (b *Blogs) CreateBlog(w http.ResponseWriter, r *http.Request) error {
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

	outBlog, err := b.repo.Create(r.Context(), *inBlog)
	if err != nil {
		slog.Error("CreateBlog: repo create failed", "error", err.Error())
		return writeJSON(w, err, nil, http.StatusInternalServerError)
	}

	return writeJSON(w, nil, outBlog, http.StatusOK)
}

func (b *Blogs) ListBlogs(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (b *Blogs) GetBlog(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (b *Blogs) UpdateBlog(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (b *Blogs) SoftDeleteBlog(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (b *Blogs) ListDeletedBlogs(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (b *Blogs) GetDeletedBlog(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (b *Blogs) RestoreDeletedBlog(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (b *Blogs) DeleteBlog(w http.ResponseWriter, r *http.Request) error {
	return nil
}
