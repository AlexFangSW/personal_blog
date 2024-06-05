package api

import (
	"log/slog"
	"net/http"
)

type apiHandler func(w http.ResponseWriter, r *http.Request) error

// Adds middleware on top of base handler func
// Default middlewares:
// - error handling
// - logging
func withMiddleware(
	base apiHandler,
	handlers ...func(http.HandlerFunc) http.HandlerFunc,
) http.HandlerFunc {

	var finalHandler = internalError(base)
	finalHandler = logPath(finalHandler)

	for index, handler := range handlers {
		slog.Info("handler", "number", index)
		finalHandler = handler(finalHandler)
	}

	return finalHandler
}

// logging request path
func logPath(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info(r.Method + " " + r.URL.String())
		next(w, r)
	}
}

// process unexpected error
// THIS SHOULD BE THE FIRST MIDDLEWARE
func internalError(next apiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	}
}
