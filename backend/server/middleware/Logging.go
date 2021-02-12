package middleware

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var LogPath = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/metrics" {
			IP := r.Header.Get("X-Real-IP") // depends on nginx
			log.Info(fmt.Sprintf("%s: %s %s (%s)", IP, r.Method, r.RequestURI, r.Host))
		}
		next.ServeHTTP(w, r)
	})
}
