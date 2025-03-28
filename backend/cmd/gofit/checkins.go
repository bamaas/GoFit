package main

import (
	"net/http"
	"time"

	"github.com/bamaas/gofit/internal/data"
	"github.com/bamaas/gofit/internal/validator"

	"github.com/google/uuid"
)

func (app *application) getCheckInsHandler(w http.ResponseWriter, r *http.Request) {

	// Get the input
	var input struct {
		data.Filters
	}

	qs := r.URL.Query()
	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 30)

	// TODO: Refactor this into a function
	format := "2006-01-02"
	
	start := app.readString(qs, "start_time", "")
	if start != "" {
		startTime, err := time.Parse(format, start)
		if err != nil {
			app.logger.Error(err.Error())
			app.writeJSON(w, http.StatusBadRequest, envelope{"error": "invalid start_time"})
			return
		}
		input.Filters.StartTime = startTime
	}

	end := app.readString(qs, "end_time", "")
	if end != "" {
		endTime, err := time.Parse(format, end)
		if err != nil {
			app.logger.Error(err.Error())
			app.writeJSON(w, http.StatusBadRequest, envelope{"error": "invalid end_time"})
			return
		}
		input.Filters.EndTime = endTime
	}

	// Validate the input
	v := validator.New()
	v.Check(input.Filters.Page >= 1, "page", "can't be less than 1")
	v.Check(input.Filters.PageSize >= 1, "pageSize", "can't be less than 1")
	v.Check(input.Filters.PageSize <= 100, "pageSize", "can't be greater than 100")
	if !v.Valid() {
		app.writeJSON(w, http.StatusBadRequest, envelope{"invalid_query_parameters": v.Errors})
		return
	}
	if !input.Filters.StartTime.IsZero() && !input.Filters.EndTime.IsZero() {
		v := validator.New()
		v.Check(input.Filters.StartTime.Before(input.Filters.EndTime) || input.Filters.StartTime == input.Filters.EndTime, "start_time", "must be before or equal to end_time")
		if !v.Valid() {
			app.writeJSON(w, http.StatusBadRequest, envelope{"invalid_query_parameters": v.Errors})
			return
		}
	}

	// Get the user
	user := app.contextGetUser(r)

	// Retrieve data
	app.logger.Info("Getting check-ins")
	checkIns, metadata, err := app.models.CheckIns.List(user.ID, input.Filters)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error retrieving records", http.StatusInternalServerError)
		return
	}

	// Write the response
	app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "data": checkIns})
}

func (app *application) getCheckInHandler(w http.ResponseWriter, r *http.Request) {
	
	// Get the input
	input := r.PathValue("uuid")
	if input == "" {
		http.Error(w, "uuid is required", http.StatusBadRequest)
		return
	}

	// Get the user
	user := app.contextGetUser(r)

	// Retrieve data
	app.logger.Info("Getting check-in", "UUID", input)
	checkIn, err := app.models.CheckIns.Get(user.ID, input)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "record not found", http.StatusNotFound)
		return
	}

	// Write the response
	app.writeJSON(w, http.StatusOK, checkIn)
}

func (app *application) createCheckIn(w http.ResponseWriter, r *http.Request) {

	// Get the input
	var input struct {
		Notes    string  `json:"notes"`
		Weight   float64 `json:"weight"`
		Datetime int64   `json:"datetime"`
	}
	err := app.readJSON(w, r, &input)	// first decode into input struct to prevent decoding of i.e. uuid field
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the user
	user := app.contextGetUser(r)

	// Validate the input
	c := &data.CheckIn{
		Datetime: input.Datetime,
		Weight: input.Weight,
		Notes: input.Notes,
		UserID: user.ID,
	}
	v := validator.New()
	v.ValidateCheckIn(c)
	if !v.Valid() {
		app.writeJSON(w, http.StatusUnprocessableEntity, envelope{"invalid_fields": v.Errors})
		return
	}

	// Generate uuid
	uuid, err := uuid.NewRandom()
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error inserting record", http.StatusInternalServerError)
		return
	}
	c.UUID = uuid.String()

	// Insert the record
	app.logger.Info("Creating check-in", "UUID", c.UUID)
	if err = app.models.CheckIns.Insert(*c); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error inserting record", http.StatusInternalServerError)
		return
	}

	// Write the response
	app.writeJSON(w, http.StatusCreated, envelope{"data": c})
}

func (app *application) deleteCheckIn(w http.ResponseWriter, r *http.Request) {

	input := r.PathValue("uuid")
	if input == "" {
		http.Error(w, "uuid is required", http.StatusBadRequest)
		return
	}

	// Get the user
	user := app.contextGetUser(r)

	app.logger.Info("Deleting check-in", "uuid", input)
	if err := app.models.CheckIns.Delete(user.ID, input); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error deleting record", http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "record deleted"})
}

func (app *application) updateCheckIn(w http.ResponseWriter, r *http.Request) {

	// Get the input
	var c data.CheckIn
	err := app.readJSON(w, r, &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the user and append to the input
	user := app.contextGetUser(r)
	c.UserID = user.ID

	// Validate the input
	v := validator.New()
	v.ValidateCheckIn(&c)
	if !v.Valid() {
		app.writeJSON(w, http.StatusUnprocessableEntity, envelope{"invalid_fields": v.Errors})
		return
	}

	// Update the record
	app.logger.Info("Updating check-in")
	if err = app.models.CheckIns.Update(c); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error updating", http.StatusInternalServerError)
		return
	}

	// Write the response
	app.writeJSON(w, http.StatusCreated, envelope{"data": c})
}
