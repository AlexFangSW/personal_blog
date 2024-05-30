package api

import (
	"blog/structs"
	"log/slog"
	"net/http"
)

type Server struct {
	Config *structs.Config
}

func NewServer(config structs.Config) *Server {
	// set up db connection here...
	return &Server{
		Config: &config,
	}
}

func (s *Server) Start() error {
	// routes
	mux := &http.ServeMux{}

	mux.HandleFunc(s.post("/blogs"), withMiddleware(s.hBlogsCreate))
	mux.HandleFunc(s.get("/blogs"), withMiddleware(s.hBlogsGetAll))
	mux.HandleFunc(s.get("/blogs/{id}"), withMiddleware(s.hBlogsGetOne))
	mux.HandleFunc(s.patch("/blogs/{id}"), withMiddleware(s.hBlogsUpdate))
	mux.HandleFunc(s.delete("/blogs/{id}"), withMiddleware(s.hBlogsDelete))
	mux.HandleFunc(s.get("/blogs/deleted"), withMiddleware(s.hBlogsGetDeleteAll))
	mux.HandleFunc(s.get("/blogs/deleted/{id}"), withMiddleware(s.hBlogsGetDeleteOne))
	mux.HandleFunc(s.patch("/blogs/deleted"), withMiddleware(s.hBlogsDeleteRestore))
	mux.HandleFunc(s.delete("/blogs/deleted"), withMiddleware(s.hBlogsTrueDelete))

	mux.HandleFunc(s.post("/tags"), withMiddleware(s.hTagsCreate))
	mux.HandleFunc(s.get("/tags"), withMiddleware(s.hTagsGetAll))
	mux.HandleFunc(s.get("/tags/{id}"), withMiddleware(s.hTagsGetOne))
	mux.HandleFunc(s.patch("/tags/{id}"), withMiddleware(s.hTagsUpdate))
	mux.HandleFunc(s.delete("/tags/{id}"), withMiddleware(s.hTagsDelete))

	mux.HandleFunc(s.post("/topics"), withMiddleware(s.hTopicsCreate))
	mux.HandleFunc(s.get("/topics"), withMiddleware(s.hTopicsGetAll))
	mux.HandleFunc(s.get("/topics/{id}"), withMiddleware(s.hTopicsGetOne))
	mux.HandleFunc(s.patch("/topics/{id}"), withMiddleware(s.hTopicsUpdate))
	mux.HandleFunc(s.delete("/topics/{id}"), withMiddleware(s.hTopicsDelete))

	slog.Info("Server is listening on", "port", s.Config.Server.Port)
	return http.ListenAndServe(s.Config.Server.Port, mux)
}
