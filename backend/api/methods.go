package api

import (
	"log/slog"
)

func (s *Server) get(path string) string {
	prefix := s.config.Prefix
	result := "GET " + prefix + path
	slog.Info("new route", "route", result)
	return result
}

func (s *Server) post(path string) string {
	prefix := s.config.Prefix
	result := "POST " + prefix + path
	slog.Info("new route", "route", result)
	return result
}

func (s *Server) delete(path string) string {
	prefix := s.config.Prefix
	result := "DELETE " + prefix + path
	slog.Info("new route", "route", result)
	return result
}

func (s *Server) patch(path string) string {
	prefix := s.config.Prefix
	result := "PATCH " + prefix + path
	slog.Info("new route", "route", result)
	return result
}

func (s *Server) put(path string) string {
	prefix := s.config.Prefix
	result := "PUT " + prefix + path
	slog.Info("new route", "route", result)
	return result
}
