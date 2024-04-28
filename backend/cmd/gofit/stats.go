package main

import (
	"net/http"
)

func (app *application) getStatsHandler(w http.ResponseWriter, r *http.Request) {

	stats, err := app.models.Stats.GetStats()
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error getting stats", http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"stats": stats})
}