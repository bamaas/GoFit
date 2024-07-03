package main

import (
	"net/http"
	"strings"

	"github.com/bamaas/gofit/internal/data"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Add the Vary: Authorization header to the response to indicate that the response may vary based on the value of the Authorization header
		w.Header().Add("Vary", "Authorization")

		authHeader := r.Header.Get("Authorization")

		// Check if the authorization header is empty
		if authHeader == "" {
			app.writeJSON(w, http.StatusUnauthorized, envelope{"error": "authorization header required"})
			return
		}

		// Split the authorization header
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			app.writeJSON(w, http.StatusUnauthorized, envelope{"error": "invalid authorization header"})
			return
		}

		// Check if the authorization type is not Bearer
		if headerParts[0] != "Bearer" {
			app.writeJSON(w, http.StatusUnauthorized, envelope{"error": "authorization header must be a bearer token"})
			return
		}

		// Check if the token is empty
		token := headerParts[1]
		if token == "" {
			app.writeJSON(w, http.StatusUnauthorized, envelope{"error": "authorization token required"})
			return
		}

		// Get user for token
		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {		// TODO: implement better error handling
			app.writeJSON(w, http.StatusUnauthorized, envelope{"error": "error fetching user from token"})
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)

	})
}

func (app *application) logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Debug("request", "method", r.Method, "url", r.URL.String(), )
		next.ServeHTTP(w, r)
	})
}