package main

import (
	"net/http"

	"github.com/bamaas/gofit/internal/data"

	"github.com/google/uuid"
)

func (app *application) getCheckInsHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		data.Filters
	}

	qs := r.URL.Query()
	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 30)

	app.logger.Info("Getting check-ins")
	checkIns, metadata, err := app.models.CheckIns.List(input.Filters)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error retrieving records", http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"metadata": metadata, "data": checkIns})
}

func (app *application) getCheckInHandler(w http.ResponseWriter, r *http.Request) {
	
	input := r.PathValue("uuid")
	if input == "" {
		http.Error(w, "uuid is required", http.StatusBadRequest)
		return
	}

	app.logger.Info("Getting check-in", "UUID", input)
	checkIn, err := app.models.CheckIns.Get(input)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "record not found", http.StatusNotFound)
		return
	}
	app.writeJSON(w, http.StatusOK, checkIn)
}

func (app *application) createCheckIn(w http.ResponseWriter, r *http.Request) {

	// var input struct {
	// 	Datetime time.Time `json:"datetime"`
	// 	Weight float64 		`json:"weight"`
	// }

	var input data.CheckIn

	err := app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Weight < 20 {
		http.Error(w, "weight must be greater than 20", http.StatusBadRequest)
		return
	}
	uuid, err := uuid.NewRandom()
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error inserting record", http.StatusInternalServerError)
		return
	}
	input.UUID = uuid.String()

	app.logger.Info("Creating check-in", "UUID", input.UUID)
	if err = app.models.CheckIns.Insert(input); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error inserting record", http.StatusInternalServerError)
		return
	}
}

func (app *application) deleteCheckIn(w http.ResponseWriter, r *http.Request) {

	input := r.PathValue("uuid")
	if input == "" {
		http.Error(w, "uuid is required", http.StatusBadRequest)
		return
	}

	app.logger.Info("Deleting check-in", "uuid", input)
	if err := app.models.CheckIns.Delete(input); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error deleting record", http.StatusInternalServerError)
		return
	}
}

func (app *application) updateCheckIn(w http.ResponseWriter, r *http.Request) {

	var input data.CheckIn
	err := app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.logger.Info("Updating check-in")
	if err = app.models.CheckIns.Update(input); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error updating", http.StatusInternalServerError)
		return
	}
}
