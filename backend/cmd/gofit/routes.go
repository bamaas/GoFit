package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/health", app.healthcheckHandler)
	mux.HandleFunc("GET /v1/check-ins", app.getCheckInsHandler)
	mux.HandleFunc("GET /v1/check-ins/{uuid}", app.getCheckInHandler)
	mux.HandleFunc("POST /v1/check-ins", app.createCheckIn)
	mux.HandleFunc("DELETE /check-ins/{uuid}", app.deleteCheckIn)
	mux.HandleFunc("PUT /check-ins/", app.updateCheckIn)

	return mux
}
