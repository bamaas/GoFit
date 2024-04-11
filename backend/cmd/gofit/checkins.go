package main

import (
	"encoding/json"
	"io"
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

	checkIns, err := app.models.CheckIns.List(input.Filters)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error getting check-ins", http.StatusInternalServerError)
		return
	}
	app.writeJSON(w, http.StatusOK, checkIns)
}

func (app *application) getCheckInHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	app.logger.Info("Getting check-in", "UUID", uuid)
	checkIn, err := app.models.CheckIns.Get(uuid)
	if err != nil {
		http.Error(w, "check-in not found", http.StatusNotFound)
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
		http.Error(w, "error generating UUID", http.StatusInternalServerError)
		return
	}
	input.UUID = uuid.String()
	app.logger.Debug("Creating check-in", "check-in", input)
	if err = app.models.CheckIns.Insert(input); err != nil {
		http.Error(w, "error inserting record into database", http.StatusInternalServerError)
	}
}

func (app *application) deleteCheckIn(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	app.logger.Info("Deleting check-in", "uuid", uuid)
	if err := app.models.CheckIns.Delete(uuid); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error deleting", http.StatusInternalServerError)
	}
}

func (app *application) updateCheckIn(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Updating check-in")
	body, err := io.ReadAll(r.Body)
	var e data.CheckIn
	if err != nil {
		http.Error(w, "error reading body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &e)
	if err != nil {
		http.Error(w, "error parsing body", http.StatusInternalServerError)
		return
	}
	if err = app.models.CheckIns.Update(e); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error updating", http.StatusInternalServerError)
	}
}
