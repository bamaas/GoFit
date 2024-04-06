package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/bamaas/gofit/internal/database"

	"github.com/google/uuid"
)

func (app *application) getCheckInsHandler(w http.ResponseWriter, r *http.Request) {
	checkIns, err := app.database.GetCheckIns()
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
	checkIn, err := app.database.GetCheckIn(uuid)
	if err != nil {
		http.Error(w, "check-in not found", http.StatusNotFound)
		return
	}
	app.writeJSON(w, http.StatusOK, checkIn)
}

func (app *application) createCheckIn(w http.ResponseWriter, r *http.Request) {
	time.Sleep(5 * time.Second)
	body, err := io.ReadAll(r.Body)
	var checkIn database.CheckIn
	if err != nil {
		http.Error(w, "error reading body", http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &checkIn)
	if err != nil {
		http.Error(w, "error parsing body", http.StatusInternalServerError)
		return
	}
	if checkIn.Weight < 20 {
		http.Error(w, "weight must be greater than 20", http.StatusBadRequest)
		return
	}
	uuid, err := uuid.NewRandom()
	if err != nil {
		http.Error(w, "error generating UUID", http.StatusInternalServerError)
		return
	}
	checkIn.UUID = uuid.String()
	app.logger.Debug("Creating check-in", "check-in", checkIn)
	if err = app.database.InsertCheckIn(checkIn); err != nil {
		http.Error(w, "error inserting record into database", http.StatusInternalServerError)
	}
}

func (app *application) deleteCheckIn(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	app.logger.Info("Deleting check-in", "uuid", uuid)
	if err := app.database.DeleteCheckIn(uuid); err != nil {
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
