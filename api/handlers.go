package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

func (s *Server) handlerHelloWorld(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello world")

	dummy := r.PathValue("dummy")
	slog.Info("got path value", "dummy", dummy)

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"dummy": dummy})
}
