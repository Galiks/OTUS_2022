package internalhttp

import (
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type logMiddleware struct {
	http.ResponseWriter
	statusCode int
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &logMiddleware{w, http.StatusOK}
		next.ServeHTTP(lrw, r)
		log.Info().
			Str("ip", strings.Split(r.RemoteAddr, ":")[0]).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("version", r.Proto).
			Int("status_code", lrw.statusCode).
			Str("latency", time.Since(start).String()).
			Str("agent", r.UserAgent()).
			Msg("")
	})
}
