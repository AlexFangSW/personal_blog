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
type blogsRepository interface {
	Create(ctx context.Context, blog entities.InBlog) (*entities.OutBlog, error)
	Update(ctx context.Context, blog entities.InBlog, id int) (*entities.OutBlog, error)

	// This group of functions will only return rows with 'visible=true' and 'deleted_at=""'
	Get(ctx context.Context, id int) (*entities.OutBlog, error)
	List(ctx context.Context) ([]entities.OutBlog, error)
	ListByTopicIDs(ctx context.Context, topicID []int) ([]entities.OutBlog, error)
	ListByTopicAndTagIDs(ctx context.Context, topicID, tagID []int) ([]entities.OutBlog, error)

	// Returns all rows regardless of visiblility and soft delete status
	AdminGet(ctx context.Context, id int) (*entities.OutBlog, error)
	AdminList(ctx context.Context) ([]entities.OutBlog, error)
	AdminListByTopicIDs(ctx context.Context, topicID []int) ([]entities.OutBlog, error)
	AdminListByTopicAndTagIDs(ctx context.Context, topicID, tagID []int) ([]entities.OutBlog, error)

	SoftDelete(ctx context.Context, id int) (int, error)
	// blogs need to be soft deleted first to be deleted
	Delete(ctx context.Context, id int) (int, error)
	RestoreDeleted(ctx context.Context, id int) (*entities.OutBlog, error)
}

type Blogs struct {
	repo blogsRepository
}

func NewBlogs(repo blogsRepository) *Blogs {
	return &Blogs{
		repo: repo,
	}
}

// CreateBlog
//
//	@Summary		Create blog
//	@Description	blogs must have unique titles
//	@Tags			blogs
//	@Accept			json
//	@Produce		json
//	@Param			blog	body		entities.ReqInBlog	true	"new blog contents"
//	@Success		200		{object}	entities.RetSuccess[entities.OutBlog]
//	@Failure		400		{object}	entities.RetFailed
//	@Failure		500		{object}	entities.RetFailed
//	@Router			/blogs [post]
func (b *Blogs) CreateBlog(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("CreateTag")

	body := &entities.ReqInBlog{}
	if err := json.NewDecoder(r.Body).Decode(body); err != nil {
		slog.Error("CreateBlog: decode failed", "error", err.Error())
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
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
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	return entities.NewRetSuccess[entities.OutBlog](*outBlog).WriteJSON(w)
}

// ListBlogs
//
//	@Summary		List blogs
//	@Description	list blogs
//	@Tags			blogs
//	@Accept			json
//	@Produce		json
//	@Param			all		query		bool	false	"show all blogs regardless of visibility or soft delete status"	default(false)
//	@Param			topic	query		int		false	"filter by topic ids, can be multiple ids. ex: ?topic=1&topic=2"
//	@Param			tag		query		int		false	"filter by tag ids, can be multiple ids, must be use with topic. ex: ?tag=1&tag=2"
//	@Success		200		{object}	entities.RetSuccess[[]entities.OutBlog]
//	@Failure		400		{object}	entities.RetFailed
//	@Failure		500		{object}	entities.RetFailed
//	@Router			/blogs [get]
func (b *Blogs) ListBlogs(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("ListBlogs")

	// process queries
	queries := r.URL.Query()
	slog.Debug("got queries", "queries", queries)

	rowAll := queries["all"]
	all, err := strListToBool(rowAll)
	if err != nil {
		slog.Error("ListBlogs: 'all' string list to bool failed", "error", err)
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}
	all = removeDuplicate[bool](all)

	rowTopicIDs := queries["topic"]
	topicIDs, err := strListToInt(rowTopicIDs)
	if err != nil {
		slog.Error("ListBlogs: 'topic' string list to int failed", "error", err)
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}
	topicIDs = removeDuplicate[int](topicIDs)

	rowTagIDs := queries["tag"]
	tagIDs, err := strListToInt(rowTagIDs)
	if err != nil {
		slog.Error("ListBlogs: 'tag' string list to int failed", "error", err)
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}
	tagIDs = removeDuplicate[int](tagIDs)

	// admin list
	if len(all) > 0 && all[0] == true {
		// admin list by topic and tag ids
		if len(topicIDs) > 0 && len(tagIDs) > 0 {
			blogs, err := b.repo.AdminListByTopicAndTagIDs(r.Context(), topicIDs, tagIDs)
			if err != nil {
				slog.Error("ListBlogs: admin list by topic and tag ids failed", "error", err)
				if errors.Is(err, sql.ErrNoRows) {
					return entities.NewRetSuccess[[]entities.OutBlog]([]entities.OutBlog{}).WriteJSON(w)
				}
				return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
			}
			return entities.NewRetSuccess[[]entities.OutBlog](blogs).WriteJSON(w)
		}

		// admin list by topic ids
		if len(topicIDs) > 0 {
			blogs, err := b.repo.AdminListByTopicIDs(r.Context(), topicIDs)
			if err != nil {
				slog.Error("ListBlogs: admin list by topic ids failed", "error", err)
				if errors.Is(err, sql.ErrNoRows) {
					return entities.NewRetSuccess[[]entities.OutBlog]([]entities.OutBlog{}).WriteJSON(w)
				}
				return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
			}
			return entities.NewRetSuccess[[]entities.OutBlog](blogs).WriteJSON(w)
		}

		// admin list
		blogs, err := b.repo.AdminList(r.Context())
		if err != nil {
			slog.Error("ListBlogs: admin list failed", "error", err)
			if errors.Is(err, sql.ErrNoRows) {
				return entities.NewRetSuccess[[]entities.OutBlog]([]entities.OutBlog{}).WriteJSON(w)
			}
			return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
		}
		return entities.NewRetSuccess[[]entities.OutBlog](blogs).WriteJSON(w)
	}

	// normal list blogs, only list blogs that are visible and not soft deleted

	// list by topic and tag ids
	if len(topicIDs) > 0 && len(tagIDs) > 0 {
		blogs, err := b.repo.ListByTopicAndTagIDs(r.Context(), topicIDs, tagIDs)
		if err != nil {
			slog.Error("ListBlogs: list by topic and tag ids failed", "error", err)
			if errors.Is(err, sql.ErrNoRows) {
				return entities.NewRetSuccess[[]entities.OutBlog]([]entities.OutBlog{}).WriteJSON(w)
			}
			return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
		}
		return entities.NewRetSuccess[[]entities.OutBlog](blogs).WriteJSON(w)
	}

	// list by topic ids
	if len(topicIDs) > 0 {
		blogs, err := b.repo.ListByTopicIDs(r.Context(), topicIDs)
		if err != nil {
			slog.Error("ListBlogs: list by topic ids failed", "error", err)
			if errors.Is(err, sql.ErrNoRows) {
				return entities.NewRetSuccess[[]entities.OutBlog]([]entities.OutBlog{}).WriteJSON(w)
			}
			return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
		}
		return entities.NewRetSuccess[[]entities.OutBlog](blogs).WriteJSON(w)
	}

	// list
	blogs, err := b.repo.List(r.Context())
	if err != nil {
		slog.Error("ListBlogs: list failed", "error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewRetSuccess[[]entities.OutBlog]([]entities.OutBlog{}).WriteJSON(w)
		}
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}
	return entities.NewRetSuccess[[]entities.OutBlog](blogs).WriteJSON(w)
}

// GetBlog
//
//	@Summary		Get blog
//	@Description	get blog
//	@Tags			blogs
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int		true	"target blog id"
//	@Param			all	query		bool	false	"show all blogs regardless of visibility or soft delete status"	default(false)
//	@Success		200	{object}	entities.RetSuccess[entities.OutBlog]
//	@Failure		400	{object}	entities.RetFailed
//	@Failure		500	{object}	entities.RetFailed
//	@Router			/blogs/{id} [get]
func (b *Blogs) GetBlog(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("GetBlog")

	// process path param
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("GetBlog: id path param to int failed", "error", err)
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}

	// process queries
	queries := r.URL.Query()
	slog.Debug("got queries", "queries", queries)

	rowAll := queries["all"]
	all, err := strListToBool(rowAll)
	if err != nil {
		slog.Error("GetBlog: 'all' string list to bool failed", "error", err)
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}

	// admin get
	if len(all) > 0 && all[0] == true {
		blog, err := b.repo.AdminGet(r.Context(), id)
		if err != nil {
			slog.Error("GetBlog: admin get failed", "error", err)
			if errors.Is(err, sql.ErrNoRows) {
				return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
			}
			return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
		}
		return entities.NewRetSuccess[entities.OutBlog](*blog).WriteJSON(w)
	}

	// normal get
	blog, err := b.repo.Get(r.Context(), id)
	if err != nil {
		slog.Error("GetBlog: get failed", "error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
		}
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}
	return entities.NewRetSuccess[entities.OutBlog](*blog).WriteJSON(w)
}

// UpdateBlog
//
//	@Summary		Update blog
//	@Description	update blog
//	@Tags			blogs
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"target blog id"
//	@Param			blog	body		entities.ReqInBlog	true	"new blog content"
//	@Success		200		{object}	entities.RetSuccess[entities.OutBlog]
//	@Failure		400		{object}	entities.RetFailed
//	@Failure		500		{object}	entities.RetFailed
//	@Router			/blogs/{id} [put]
func (b *Blogs) UpdateBlog(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("UpdateBlog")

	// process path param
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("UpdateBlog: id path param to int failed", "error", err)
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}

	// process body
	blog := &entities.ReqInBlog{}
	if err := json.NewDecoder(r.Body).Decode(blog); err != nil {
		slog.Error("UpdateBlog: parse body param failed", "error", err)
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}
	newBlog := entities.NewBlog(
		blog.Title,
		blog.Content,
		blog.Description,
		blog.Pined,
		blog.Visible,
	)
	inBlog := entities.NewInBlog(
		*newBlog,
		blog.Tags,
		blog.Topics,
	)

	// update
	updatedBlog, err := b.repo.Update(r.Context(), *inBlog, id)
	if err != nil {
		slog.Error("UpdateBlog: update failed", "error", err)
		if errors.Is(err, sql.ErrNoRows) {
			return entities.NewRetFailed(ErrorTargetNotFound, http.StatusBadRequest).WriteJSON(w)
		}
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}
	return entities.NewRetSuccess[entities.OutBlog](*updatedBlog).WriteJSON(w)
}

// SoftDeleteBlog
//
//	@Summary		Soft delete blog
//	@Description	update blog
//	@Tags			blogs
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"target blog id"
//	@Success		200	{object}	entities.RetSuccess[entities.RowsAffected]
//	@Failure		400	{object}	entities.RetFailed
//	@Failure		500	{object}	entities.RetFailed
//	@Router			/blogs/{id} [delete]
func (b *Blogs) SoftDeleteBlog(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("SoftDeleteBlog")

	// parse path param
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("SoftDeleteBlog: id path param to int failed", "error", err)
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}

	// soft delete blog
	affectedRows, err := b.repo.SoftDelete(r.Context(), id)
	if err != nil {
		slog.Error("SoftDeleteBlog: soft delete failed", "error", err)
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}
	if affectedRows == 0 {
		slog.Error("SoftDeleteBlog: target not failed", "error", err)
		return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
	}

	return entities.NewRetSuccess[entities.RowsAffected](*entities.NewRowsAffected(affectedRows)).WriteJSON(w)
}

// RestoreDeletedBlog
//
//	@Summary		Restore delete blog
//	@Description	restore delete blog
//	@Tags			blogs
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"target blog id"
//	@Success		200	{object}	entities.RetSuccess[entities.OutBlog]
//	@Failure		400	{object}	entities.RetFailed
//	@Failure		500	{object}	entities.RetFailed
//	@Router			/blogs/deleted/{id} [patch]
func (b *Blogs) RestoreDeletedBlog(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("RestoreDeletedBlog")

	// parse path param
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("RestoreDeletedBlog: id path param to int failed", "error", err)
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}

	// restore soft deleted blog
	blog, err := b.repo.RestoreDeleted(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("RestoreDeletedBlog: restore blog failed", "error", err)
			return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
		}
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}

	return entities.NewRetSuccess[entities.OutBlog](*blog).WriteJSON(w)
}

// DeleteBlog
//
//	@Summary		Delete blog
//	@Description	delete blog
//	@Tags			blogs
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"target blog id"
//	@Success		200	{object}	entities.RetSuccess[entities.OutBlog]
//	@Failure		400	{object}	entities.RetFailed
//	@Failure		500	{object}	entities.RetFailed
//	@Router			/blogs/deleted/{id} [delete]
func (b *Blogs) DeleteBlog(w http.ResponseWriter, r *http.Request) error {
	slog.Debug("DeleteBlog")

	// parse path param
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error("DeleteBlog: id path param to int failed", "error", err)
		return entities.NewRetFailed(err, http.StatusBadRequest).WriteJSON(w)
	}

	// delete blog
	affectedRows, err := b.repo.Delete(r.Context(), id)
	if err != nil {
		slog.Error("DeleteBlog: delete failed", "error", err)
		return entities.NewRetFailed(err, http.StatusInternalServerError).WriteJSON(w)
	}
	if affectedRows == 0 {
		slog.Error("DeleteBlog: target not failed", "error", err)
		return entities.NewRetFailed(ErrorTargetNotFound, http.StatusNotFound).WriteJSON(w)
	}

	return entities.NewRetSuccess[entities.RowsAffected](*entities.NewRowsAffected(affectedRows)).WriteJSON(w)
}
