package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("GET /v1/check-ins", app.getCheckInsHandler)
	mux.HandleFunc("GET /v1/check-ins/{uuid}", app.getCheckInHandler)
	mux.HandleFunc("POST /v1/check-ins", app.createCheckIn)
	mux.HandleFunc("DELETE /entry/{uuid}/", app.deleteCheckIn)
	mux.HandleFunc("PUT /entry/", app.updateCheckIn)

	return mux
}
