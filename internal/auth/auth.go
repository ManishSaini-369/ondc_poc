package auth

import (
	"net/http"
    "log"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

        // TODO: Implement the full signature verification logic here.
        log.Println("Auth middleware is not fully implemented yet.")

		next.ServeHTTP(w, r)
	}
}