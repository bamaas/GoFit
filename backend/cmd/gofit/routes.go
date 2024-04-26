package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", app.healthcheckHandler)

	// Check-ins
	mux.HandleFunc("GET /v1/check-ins", app.getCheckInsHandler)
	mux.HandleFunc("GET /v1/check-ins/{uuid}", app.getCheckInHandler)
	mux.HandleFunc("POST /v1/check-ins", app.createCheckIn)
	mux.HandleFunc("DELETE /v1/check-ins/{uuid}", app.deleteCheckIn)
	mux.HandleFunc("PUT /v1/check-ins/", app.updateCheckIn)

	// Stats
	mux.HandleFunc("GET /v1/stats", app.getStatsHandler)

	// Users
	mux.HandleFunc("POST /v1/users", app.registerUserHandler)

	return mux
}
