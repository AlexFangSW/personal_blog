package api

import (
	"blog/api/handlers"
	"blog/config"
	_ "blog/swagger_docs"
	"context"
	"fmt"
	"log/slog"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Server struct {
	server *http.Server
	config config.Config
	blogs  handlers.Blogs
	topics handlers.Topics
	tags   handlers.Tags
	users  handlers.Users
	probes handlers.Probes
}

func NewServer(
	config config.Config,
	blogs handlers.Blogs,
	tags handlers.Tags,
	topics handlers.Topics,
	users handlers.Users,
	probes handlers.Probes) *Server {
	return &Server{
		config: config,
		blogs:  blogs,
		tags:   tags,
		topics: topics,
		users:  users,
		probes: probes,
	}
}

func (s *Server) Start() error {
	// routes
	mux := http.NewServeMux()

	// api specification
	slog.Info("API specification at: /docs/*")
	filepath := fmt.Sprintf("http://localhost:%d/docs/doc.json", s.config.Server.Port)
	mux.HandleFunc("GET /docs/*", WithMiddleware(apiHandlerWrapper(
		httpSwagger.Handler(httpSwagger.URL(filepath)))))

	// authentication
	loginRateLimit := NewRateLimit(s.config.Login.RateLimit, s.config.Login.RateLimit*2)
	authCheckRateLimit := NewRateLimit(s.config.Login.RateLimit, s.config.Login.RateLimit*2)
	mux.HandleFunc(s.post("/login"), WithMiddleware(s.users.Login, loginRateLimit.RateLimit))
	mux.HandleFunc(s.post("/logout"), WithMiddleware(s.users.Logout))
	mux.HandleFunc(s.post("/auth-check"), WithMiddleware(s.users.AuthorizeCheck, authCheckRateLimit.RateLimit))

	mux.HandleFunc(s.post("/blogs"), WithMiddleware(s.blogs.CreateBlog))
	mux.HandleFunc(s.get("/blogs"), WithMiddleware(s.blogs.ListBlogs))
	mux.HandleFunc(s.get("/blogs/{id}"), WithMiddleware(s.blogs.GetBlog))
	mux.HandleFunc(s.put("/blogs/{id}"), WithMiddleware(s.blogs.UpdateBlog))
	mux.HandleFunc(s.post("/blogs/{id}"), WithMiddleware(s.blogs.CreateBlogWithID))
	mux.HandleFunc(s.delete("/blogs/{id}"), WithMiddleware(s.blogs.SoftDeleteBlog))
	mux.HandleFunc(s.delete("/blogs/deleted/{id}"), WithMiddleware(s.blogs.DeleteBlog))
	mux.HandleFunc(s.delete("/blogs/delete-now/{id}"), WithMiddleware(s.blogs.DeleteBlogNow))
	mux.HandleFunc(s.patch("/blogs/deleted/{id}"), WithMiddleware(s.blogs.RestoreDeletedBlog))

	mux.HandleFunc(s.post("/tags"), WithMiddleware(s.tags.CreateTag))
	mux.HandleFunc(s.get("/tags"), WithMiddleware(s.tags.ListTags))
	mux.HandleFunc(s.get("/tags/{id}"), WithMiddleware(s.tags.GetTag))
	mux.HandleFunc(s.put("/tags/{id}"), WithMiddleware(s.tags.UpdateTag))
	mux.HandleFunc(s.delete("/tags/{id}"), WithMiddleware(s.tags.DeleteTag))

	mux.HandleFunc(s.post("/topics"), WithMiddleware(s.topics.CreateTopic))
	mux.HandleFunc(s.get("/topics"), WithMiddleware(s.topics.ListTopics))
	mux.HandleFunc(s.get("/topics/{id}"), WithMiddleware(s.topics.GetTopic))
	mux.HandleFunc(s.put("/topics/{id}"), WithMiddleware(s.topics.UpdateTopic))
	mux.HandleFunc(s.delete("/topics/{id}"), WithMiddleware(s.topics.DeleteTopic))

	mux.HandleFunc(s.getRoot("/alive"), WithMiddlewareDebugAccessLog(s.probes.LivenessProbe))
	mux.HandleFunc(s.getRoot("/ready"), WithMiddlewareDebugAccessLog(s.probes.ReadinessProbe))

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Server.Port),
		Handler: mux,
	}
	slog.Info("Server is listening on", "port", s.config.Server.Port)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	slog.Warn("Stop: server shutting down")
	return s.server.Shutdown(ctx)
}
