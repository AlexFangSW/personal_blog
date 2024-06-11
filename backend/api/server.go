package api

import (
	"blog/api/handlers"
	"blog/config"
	_ "blog/docs"
	"context"
	"fmt"
	"log/slog"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	server *http.Server
	config config.ServerSetting
	blogs  handlers.Blogs
	topics handlers.Topics
	tags   handlers.Tags
	users  handlers.Users
}

func NewServer(
	config config.ServerSetting,
	blogs handlers.Blogs,
	tags handlers.Tags,
	topics handlers.Topics,
	users handlers.Users) *Server {
	return &Server{
		config: config,
		blogs:  blogs,
		tags:   tags,
		topics: topics,
		users:  users,
	}
}

func (s *Server) Start() error {
	// routes
	mux := http.NewServeMux()

	// api specification
	slog.Info("API specification at: /docs/*")
	filepath := fmt.Sprintf("http://localhost%s/docs/doc.json", s.config.Port)
	mux.HandleFunc("GET /docs/*", withMiddleware(apiHandlerWrapper(
		httpSwagger.Handler(httpSwagger.URL(filepath)))))

	// authentication
	mux.HandleFunc(s.post("/login"), withMiddleware(s.users.Login))
	mux.HandleFunc(s.post("/logout"), withMiddleware(s.users.Logout))

	// TODO: block 'list, get ?all=true' requests that dosen't have token
	mux.HandleFunc(s.post("/blogs"), withMiddleware(s.blogs.CreateBlog))
	mux.HandleFunc(s.get("/blogs"), withMiddleware(s.blogs.ListBlogs))
	mux.HandleFunc(s.get("/blogs/{id}"), withMiddleware(s.blogs.GetBlog))
	mux.HandleFunc(s.put("/blogs/{id}"), withMiddleware(s.blogs.UpdateBlog))
	mux.HandleFunc(s.delete("/blogs/{id}"), withMiddleware(s.blogs.SoftDeleteBlog))
	mux.HandleFunc(s.delete("/blogs/deleted/{id}"), withMiddleware(s.blogs.DeleteBlog))
	mux.HandleFunc(s.patch("/blogs/deleted/{id}"), withMiddleware(s.blogs.RestoreDeletedBlog))

	mux.HandleFunc(s.post("/tags"), withMiddleware(s.tags.CreateTag))
	mux.HandleFunc(s.get("/tags"), withMiddleware(s.tags.ListTags))
	mux.HandleFunc(s.get("/tags/{id}"), withMiddleware(s.tags.GetTag))
	mux.HandleFunc(s.put("/tags/{id}"), withMiddleware(s.tags.UpdateTag))
	mux.HandleFunc(s.delete("/tags/{id}"), withMiddleware(s.tags.DeleteTag))

	mux.HandleFunc(s.post("/topics"), withMiddleware(s.topics.CreateTopic))
	mux.HandleFunc(s.get("/topics"), withMiddleware(s.topics.ListTopics))
	mux.HandleFunc(s.get("/topics/{id}"), withMiddleware(s.topics.GetTopic))
	mux.HandleFunc(s.put("/topics/{id}"), withMiddleware(s.topics.UpdateTopic))
	mux.HandleFunc(s.delete("/topics/{id}"), withMiddleware(s.topics.DeleteTopic))

	s.server = &http.Server{
		Addr:    s.config.Port,
		Handler: mux,
	}
	slog.Info("Server is listening on", "port", s.config.Port)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Warn("Stop: server shutting down")
	return s.server.Shutdown(ctx)
}
