package middlewares

import (
	"fmt"
	"net/http"
)

var allowedOrigins = []string{"http://localhost:300", "http://localhost:301"}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		h := w.Header()

		if isAllowedOrigin(origin) {
			h.Set("Access-Control-Allow-Origin", origin)
		} else {
			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
			return
		}

		h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		h.Set("Access-Control-Expose-Headers", "Authorization")
		h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		h.Set("Access-Control-Allow-Credentials", "true")
		h.Set("Access-Control-Max-Age", "3600")

		fmt.Println("r.Method", r.Method)
		fmt.Println(" http.MethodOptions", http.MethodOptions)

		if r.Method == http.MethodOptions {

			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(origin string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if allowedOrigin == origin {
			return true
		}
	}

	return false
}
