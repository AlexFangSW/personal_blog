package api

import (
	"blog/db/models"
	"blog/structs"
	"log/slog"
	"net/http"
)

type Server struct {
	Config structs.Config
	models models.Models
}

func NewServer(config structs.Config, models models.Models) *Server {
	return &Server{
		Config: config,
		models: models,
	}
}

func (s *Server) Start() error {
	// routes
	mux := &http.ServeMux{}

	mux.HandleFunc(s.post("/blogs"), withMiddleware(s.CreateBlog))
	// mux.HandleFunc(s.get("/blogs"), withMiddleware(s.ListBlogs))
	// mux.HandleFunc(s.get("/blogs/{id}"), withMiddleware(s.GetBlog))
	// mux.HandleFunc(s.patch("/blogs/{id}"), withMiddleware(s.UpdateBlog))
	// mux.HandleFunc(s.delete("/blogs/{id}"), withMiddleware(s.SoftDeleteBlog))
	// mux.HandleFunc(s.get("/blogs/deleted"), withMiddleware(s.ListDeletedBlogs))
	// mux.HandleFunc(s.get("/blogs/deleted/{id}"), withMiddleware(s.GetDeletedBlog))
	// mux.HandleFunc(s.patch("/blogs/deleted"), withMiddleware(s.RestoreDeletedBlog))
	// mux.HandleFunc(s.delete("/blogs/deleted"), withMiddleware(s.DeleteBlog))

	mux.HandleFunc(s.post("/tags"), withMiddleware(s.CreateTag))
	// mux.HandleFunc(s.get("/tags"), withMiddleware(s.TagsGetAll))
	// mux.HandleFunc(s.get("/tags/{id}"), withMiddleware(s.TagsGetOne))
	// mux.HandleFunc(s.patch("/tags/{id}"), withMiddleware(s.TagsUpdate))
	// mux.HandleFunc(s.delete("/tags/{id}"), withMiddleware(s.TagsDelete))
	//
	mux.HandleFunc(s.post("/topics"), withMiddleware(s.CreateTopic))
	// mux.HandleFunc(s.get("/topics"), withMiddleware(s.TopicsGetAll))
	// mux.HandleFunc(s.get("/topics/{id}"), withMiddleware(s.TopicsGetOne))
	// mux.HandleFunc(s.patch("/topics/{id}"), withMiddleware(s.TopicsUpdate))
	// mux.HandleFunc(s.delete("/topics/{id}"), withMiddleware(s.TopicsDelete))

	slog.Info("Server is listening on", "port", s.Config.Server.Port)
	return http.ListenAndServe(s.Config.Server.Port, mux)
}
