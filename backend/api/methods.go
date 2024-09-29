package api

import (
	"log/slog"
)

// With no routing prefix prefix
func (s *Server) getRoot(path string) string {
	result := "GET " + path
	slog.Info("new route", "route", result)
	return result
}

func (s *Server) get(path string) string {
	prefix := s.config.Server.Prefix
	result := "GET " + prefix + path
	slog.Info("new route", "route", result)
	return result
}

func (s *Server) post(path string) string {
	prefix := s.config.Server.Prefix
	result := "POST " + prefix + path
	slog.Info("new route", "route", result)
	return result
}

func (s *Server) delete(path string) string {
	prefix := s.config.Server.Prefix
	result := "DELETE " + prefix + path
	slog.Info("new route", "route", result)
	return result
}

func (s *Server) patch(path string) string {
	prefix := s.config.Server.Prefix
	result := "PATCH " + prefix + path
	slog.Info("new route", "route", result)
	return result
}

func (s *Server) put(path string) string {
	prefix := s.config.Server.Prefix
	result := "PUT " + prefix + path
	slog.Info("new route", "route", result)
	return result
}
