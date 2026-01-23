package main

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

func (app *application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read auth header
			// parse it
			// decode it
			// check it
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				app.unauthorizedBasicErrorResponse(w, r, errors.New("empty authorization header"))
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				app.unauthorizedBasicErrorResponse(w, r, errors.New("authorization header format must be Basic"))
				return
			}
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				app.unauthorizedBasicErrorResponse(w, r, err)
				return
			}
			creds := strings.Split(string(decoded), ":")
			if len(creds) != 2 {
				app.unauthorizedBasicErrorResponse(w, r, errors.New("invalid credential"))
				return
			}
			username := creds[0]
			password := creds[1]
			if username != app.config.auth.username || password != app.config.auth.password {
				app.unauthorizedBasicErrorResponse(w, r, errors.New("invalid credential"))
			}
			next.ServeHTTP(w, r)
		})
	}
}
