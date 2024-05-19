package api

import (
	"log/slog"
	"net/http"
)

// Adds middleware on top of base handler func
// A middleware for logging paths will be added by default
func withMiddleware(
	base http.HandlerFunc,
	handlers ...func(http.HandlerFunc) http.HandlerFunc,
) http.HandlerFunc {

	var finalHandler = mLogPath(base)

	for index, handler := range handlers {
		slog.Info("handler", "number", index)
		finalHandler = handler(finalHandler)
	}

	return finalHandler
}

// Middleware for logging request path
func mLogPath(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info(r.Method + " " + r.URL.String())
		next(w, r)
	}
}
