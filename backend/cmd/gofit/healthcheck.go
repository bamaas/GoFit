package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, _ *http.Request) {

	data := map[string]string{
		"status":  "available",
		"version": version,
	}

	app.writeJSON(w, http.StatusOK, data)
}
