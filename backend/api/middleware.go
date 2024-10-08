package api

import (
	"log/slog"
	"net/http"

	"golang.org/x/time/rate"
)

type apiHandler func(w http.ResponseWriter, r *http.Request) error

// Adds middleware on top of base handler func
// Default middlewares:
// - error handling
// - logging
func WithMiddleware(
	base apiHandler,
	handlers ...func(http.HandlerFunc) http.HandlerFunc,
) http.HandlerFunc {

	var finalHandler = internalError(base)
	finalHandler = logPath(finalHandler, "INFO")

	for index, handler := range handlers {
		slog.Info("handler", "number", index)
		finalHandler = handler(finalHandler)
	}

	return finalHandler
}

func WithMiddlewareDebugAccessLog(
	base apiHandler,
	handlers ...func(http.HandlerFunc) http.HandlerFunc,
) http.HandlerFunc {

	var finalHandler = internalError(base)
	finalHandler = logPath(finalHandler, "DEBUG")

	for index, handler := range handlers {
		slog.Info("handler", "number", index)
		finalHandler = handler(finalHandler)
	}

	return finalHandler
}

type RateLimit struct {
	limiter *rate.Limiter
}

func NewRateLimit(average, burst int) RateLimit {
	slog.Debug("new rate limit", "average", average, "burst", burst)
	return RateLimit{
		limiter: rate.NewLimiter(rate.Limit(average), burst),
	}
}

func (rlimit *RateLimit) RateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if rlimit.limiter.Allow() == false {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next(w, r)
	}
}

// logging request path
func logPath(next http.HandlerFunc, level string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if level == "DEBUG" {
			slog.Debug(r.Method + " " + r.URL.String())
		}
		if level == "INFO" {
			slog.Info(r.Method + " " + r.URL.String())
		}
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

// Wrap handlerFunc into apiHandler
func apiHandlerWrapper(next http.HandlerFunc) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		next.ServeHTTP(w, r)
		return nil
	}
}
