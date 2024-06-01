package main

import (
	"net/http"
)

func (app *application) getStatsHandler(w http.ResponseWriter, r *http.Request) {

	// Get user
	user := app.contextGetUser(r)

	// Get stats
	stats, err := app.models.Stats.GetStats(user.ID)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error getting stats", http.StatusInternalServerError)
		return
	}

	// Write response
	app.writeJSON(w, http.StatusOK, envelope{"stats": stats})
}