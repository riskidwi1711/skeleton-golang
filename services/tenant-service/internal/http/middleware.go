package http

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

func withRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get("X-Request-ID")
		if rid == "" {
			buf := make([]byte, 8)
			_, _ = rand.Read(buf)
			rid = hex.EncodeToString(buf)
		}
		w.Header().Set("X-Request-ID", rid)
		next.ServeHTTP(w, r)
	})
}
