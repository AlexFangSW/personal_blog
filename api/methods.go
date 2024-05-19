package api

func (s *Server) get(path string) string {
	prefix := s.Config.Server.Prefix
	return "GET " + prefix + path
}

func (s *Server) post(path string) string {
	prefix := s.Config.Server.Prefix
	return "POST " + prefix + path
}

func (s *Server) delete(path string) string {
	prefix := s.Config.Server.Prefix
	return "DELETE " + prefix + path
}

func (s *Server) patch(path string) string {
	prefix := s.Config.Server.Prefix
	return "PATCH " + prefix + path
}
