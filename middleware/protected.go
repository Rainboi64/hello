package middleware

import (
	"example/hello/util"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func ProtectedHandler() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			tokenString := r.Header.Get("Authorization")
			
			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Missing authorization header")
				return
			}

			const bearerPrefix = "Bearer "
			if len(tokenString) < len(bearerPrefix) || tokenString[:len(bearerPrefix)] != bearerPrefix {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Invalid authorization header format")
				return
			}
			tokenString = tokenString[len(bearerPrefix):]

			err := util.VerifyToken(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprint(w, "Invalid token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
