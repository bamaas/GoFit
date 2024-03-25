package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/bamaas/gofit/internal/database"
)

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (app *application) getCheckInsHandler(w http.ResponseWriter, r *http.Request) {
	ci, err := app.database.GetCheckIns()
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error getting check-ins", http.StatusInternalServerError)
		return
	}
	renderJSON(w, ci)
}

func (app *application) getCheckInHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	app.logger.Info("Getting check-in", "id", id)
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusInternalServerError)
		return
	}
	e, err := app.database.GetCheckIn(idInt)
	if err != nil {
		http.Error(w, "check-in not found", http.StatusNotFound)
		return
	}
	renderJSON(w, e)
}

func (app *application) createCheckIn(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Creating check-in")
	body, err := io.ReadAll(r.Body)
	var e database.CheckIn
	if err != nil {
		http.Error(w, "error reading body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &e)
	app.logger.Debug("Creating check-in", "check-in", e)
	if err != nil {
		http.Error(w, "error parsing body", http.StatusInternalServerError)
		return
	}
	if err = app.database.InsertCheckIn(e); err != nil {
		http.Error(w, "error inserting record into database", http.StatusInternalServerError)
	}
}

func (app *application) deleteCheckIn(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	app.logger.Info("Deleting check-in", "id", id)
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusNotAcceptable)
		return
	}
	if err = app.database.DeleteCheckIn(idInt); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error deleting", http.StatusInternalServerError)
	}
}

func (app *application) updateCheckIn(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("Updating check-in")
	body, err := io.ReadAll(r.Body)
	var e database.CheckIn
	if err != nil {
		http.Error(w, "error reading body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &e)
	if err != nil {
		http.Error(w, "error parsing body", http.StatusInternalServerError)
		return
	}
	if err = app.database.UpdateCheckIn(e); err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error updating", http.StatusInternalServerError)
	}
}
