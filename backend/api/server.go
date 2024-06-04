package api

import (
	"blog/config"
	"log/slog"
	"net/http"
)

type blogsHandler interface {
	CreateBlog(w http.ResponseWriter, r *http.Request) error
	ListBlogs(w http.ResponseWriter, r *http.Request) error
	GetBlog(w http.ResponseWriter, r *http.Request) error
	UpdateBlog(w http.ResponseWriter, r *http.Request) error
	SoftDeleteBlog(w http.ResponseWriter, r *http.Request) error
	ListDeletedBlogs(w http.ResponseWriter, r *http.Request) error
	GetDeletedBlog(w http.ResponseWriter, r *http.Request) error
	RestoreDeletedBlog(w http.ResponseWriter, r *http.Request) error
	DeleteBlog(w http.ResponseWriter, r *http.Request) error
}

type topicsHandler interface {
	CreateTopic(w http.ResponseWriter, r *http.Request) error
	ListTopics(w http.ResponseWriter, r *http.Request) error
	GetTopic(w http.ResponseWriter, r *http.Request) error
	UpdateTopic(w http.ResponseWriter, r *http.Request) error
	DeleteTopic(w http.ResponseWriter, r *http.Request) error
}

type tagsHandler interface {
	CreateTag(w http.ResponseWriter, r *http.Request) error
	ListTags(w http.ResponseWriter, r *http.Request) error
	GetTag(w http.ResponseWriter, r *http.Request) error
	UpdateTag(w http.ResponseWriter, r *http.Request) error
	DeleteTag(w http.ResponseWriter, r *http.Request) error
}

type Server struct {
	config config.ServerSetting
	blogs  blogsHandler
	topics topicsHandler
	tags   tagsHandler
}

func NewServer(config config.ServerSetting, blogs blogsHandler, tags tagsHandler, topics topicsHandler) *Server {
	return &Server{
		config: config,
		blogs:  blogs,
		tags:   tags,
		topics: topics,
	}
}

func (s *Server) Start() error {
	// routes
	mux := &http.ServeMux{}

	mux.HandleFunc(s.post("/blogs"), withMiddleware(s.blogs.CreateBlog))
	// mux.HandleFunc(s.get("/blogs"), withMiddleware(s.ListBlogs))
	// mux.HandleFunc(s.get("/blogs/{id}"), withMiddleware(s.GetBlog))
	// mux.HandleFunc(s.patch("/blogs/{id}"), withMiddleware(s.UpdateBlog))
	// mux.HandleFunc(s.delete("/blogs/{id}"), withMiddleware(s.SoftDeleteBlog))
	// mux.HandleFunc(s.get("/blogs/deleted"), withMiddleware(s.ListDeletedBlogs))
	// mux.HandleFunc(s.get("/blogs/deleted/{id}"), withMiddleware(s.GetDeletedBlog))
	// mux.HandleFunc(s.patch("/blogs/deleted"), withMiddleware(s.RestoreDeletedBlog))
	// mux.HandleFunc(s.delete("/blogs/deleted"), withMiddleware(s.DeleteBlog))

	mux.HandleFunc(s.post("/tags"), withMiddleware(s.tags.CreateTag))
	mux.HandleFunc(s.get("/tags"), withMiddleware(s.tags.ListTags))
	mux.HandleFunc(s.get("/tags/{id}"), withMiddleware(s.tags.GetTag))
	// mux.HandleFunc(s.patch("/tags/{id}"), withMiddleware(s.TagsUpdate))
	// mux.HandleFunc(s.delete("/tags/{id}"), withMiddleware(s.TagsDelete))
	//
	mux.HandleFunc(s.post("/topics"), withMiddleware(s.topics.CreateTopic))
	// mux.HandleFunc(s.get("/topics"), withMiddleware(s.TopicsGetAll))
	// mux.HandleFunc(s.get("/topics/{id}"), withMiddleware(s.TopicsGetOne))
	// mux.HandleFunc(s.patch("/topics/{id}"), withMiddleware(s.TopicsUpdate))
	// mux.HandleFunc(s.delete("/topics/{id}"), withMiddleware(s.TopicsDelete))

	slog.Info("Server is listening on", "port", s.config.Port)
	return http.ListenAndServe(s.config.Port, mux)
}
