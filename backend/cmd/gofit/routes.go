package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/", app.logRequests(http.HandlerFunc(app.spaHandler)))		// TODO: wrap mux and wrap all requests

	// Healthcheck
	mux.HandleFunc("GET /api/v1/health", app.healthcheckHandler)

	// Check-ins
	mux.Handle("GET /api/v1/check-ins", app.authenticate(http.HandlerFunc(app.getCheckInsHandler)))
	mux.Handle("GET /api/v1/check-ins/{uuid}", app.authenticate(http.HandlerFunc(app.getCheckInHandler)))
	mux.Handle("POST /api/v1/check-ins", app.authenticate(http.HandlerFunc(app.createCheckIn)))
	mux.Handle("DELETE /api/v1/check-ins/{uuid}", app.authenticate(http.HandlerFunc(app.deleteCheckIn)))
	mux.Handle("PUT /api/v1/check-ins/", app.authenticate(http.HandlerFunc(app.updateCheckIn)))

	// Stats
	mux.Handle("GET /api/v1/stats", app.authenticate(http.HandlerFunc(app.getStatsHandler)))

	// Users
	mux.HandleFunc("POST /api/v1/users", app.registerUserHandler)

	// Authentication
	mux.HandleFunc("POST /api/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return mux
}
