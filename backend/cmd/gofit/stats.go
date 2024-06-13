package main

import (
	"net/http"
	"time"

	"github.com/bamaas/gofit/internal/data"
)

// func (app *application) getStatsHandler(w http.ResponseWriter, r *http.Request) {

// 	// Get user
// 	user := app.contextGetUser(r)

// 	// Get stats
// 	stats, err := app.models.Stats.GetStats(user.ID)
// 	if err != nil {
// 		app.logger.Error(err.Error())
// 		http.Error(w, "error getting stats", http.StatusInternalServerError)
// 		return
// 	}

// 	// Write response
// 	app.writeJSON(w, http.StatusOK, envelope{"stats": stats})
// }

func (app *application) getWeightDifferenceStatsHandler(w http.ResponseWriter, r *http.Request) {

	// Get the input
	var input struct {
		data.Filters
	}

	qs := r.URL.Query()
	format := "2006-01-02"
	start := app.readString(qs, "start_time", "")
	if start != "" {
		startTime, err := time.Parse(format, start)
		if err != nil {
			app.logger.Error(err.Error())
			http.Error(w, "invalid start_time", http.StatusBadRequest)
			return
		}
		input.Filters.StartTime = startTime
	}
	
	endTime, err := time.Parse(format, app.readString(qs, "end_time", "2021-11-21"))
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "invalid end_time", http.StatusBadRequest)
		return
	}
	input.Filters.EndTime = endTime

	// Validate the input
	// TODO:

	// Get user
	user := app.contextGetUser(r)

	// Get stats
	diff, err := app.models.Stats.GetWeightDifference(user.ID, input.Filters)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error getting stats", http.StatusInternalServerError)
		return
	}

	// Write response
	app.writeJSON(w, http.StatusOK, envelope{"weight_difference": &diff})
}