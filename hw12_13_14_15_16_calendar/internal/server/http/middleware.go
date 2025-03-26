package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.HandlerFunc, log Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Info(fmt.Sprintf("%s  %s  %s  %s  %s  %v",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			r.Proto,
			r.UserAgent(),
			time.Since(start)))
	}
}
