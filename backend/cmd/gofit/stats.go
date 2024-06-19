package main

import (
	"net/http"
	"time"

	"github.com/bamaas/gofit/internal/data"
	"github.com/bamaas/gofit/internal/validator"
)

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

	end := app.readString(qs, "end_time", "")
	if end != "" {
		endTime, err := time.Parse(format, end)
		if err != nil {
			app.logger.Error(err.Error())
			http.Error(w, "invalid end_time", http.StatusBadRequest)
			return
		}
		input.Filters.EndTime = endTime
	}

	// Validate the input
	if !input.Filters.StartTime.IsZero() && !input.Filters.EndTime.IsZero() {
		v := validator.New()
		v.Check(input.Filters.StartTime.Before(input.Filters.EndTime), "start_time", "must be before end_time")
		if !v.Valid() {
			app.writeJSON(w, http.StatusBadRequest, envelope{"invalid_query_parameters": v.Errors})
			return
		}
	}

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

func (app *application) getWeightAverageStatsHandler(w http.ResponseWriter, r *http.Request) {
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
	v := validator.New()
	v.Check(input.Filters.StartTime.Before(input.Filters.EndTime) || input.Filters.StartTime == input.Filters.EndTime, "start_time", "must be before or equal to end_time")
	if !v.Valid() {
		http.Error(w, "invalid query parameters", http.StatusBadRequest)
		return
	}

	// Get user
	user := app.contextGetUser(r)

	// Get stats
	weight, err := app.models.Stats.GetWeightAverage(user.ID, input.Filters)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error getting average weight", http.StatusInternalServerError)
		return
	}

	// Write response
	app.writeJSON(w, http.StatusOK, envelope{"weight_average": &weight})
}

func (app *application) getWeightAverageByWeekStatsHandler(w http.ResponseWriter, r *http.Request) {

	// // Get the input
	// var input struct {
	// 	data.Filters
	// }

	// qs := r.URL.Query()
	// format := "2006-01-02"
	
	// start := app.readString(qs, "start_time", "")
	// if start != "" {
	// 	startTime, err := time.Parse(format, start)
	// 	if err != nil {
	// 		app.logger.Error(err.Error())
	// 		http.Error(w, "invalid start_time", http.StatusBadRequest)
	// 		return
	// 	}
	// 	input.Filters.StartTime = startTime
	// }

	// end := app.readString(qs, "end_time", "")
	// if end != "" {
	// 	endTime, err := time.Parse(format, end)
	// 	if err != nil {
	// 		app.logger.Error(err.Error())
	// 		http.Error(w, "invalid end_time", http.StatusBadRequest)
	// 		return
	// 	}
	// 	input.Filters.EndTime = endTime
	// }

	// // Validate the input
	// if !input.Filters.StartTime.IsZero() && !input.Filters.EndTime.IsZero() {
	// 	v := validator.New()
	// 	v.Check(input.Filters.StartTime.Before(input.Filters.EndTime), "start_time", "must be before end_time")
	// 	if !v.Valid() {
	// 		http.Error(w, "invalid query parameters", http.StatusBadRequest)
	// 		return
	// 	}
	// }

	// Get user
	user := app.contextGetUser(r)

	// Get stats
	stats, err := app.models.Stats.GetWeightAverageByWeek(user.ID)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error getting stats", http.StatusInternalServerError)
		return
	}

	// Write response
	app.writeJSON(w, http.StatusOK, stats)
}