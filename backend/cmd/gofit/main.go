package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/bamaas/gofit/internal/config"
	"github.com/bamaas/gofit/internal/database"
	"github.com/bamaas/gofit/internal/logger"
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

func main() {

	// Retrieve config
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}

	// Setup logger
	logger, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}

	// Database
	db, err := database.New(logger)
	if err != nil{
		panic(err)
	}
	defer db.Close()

	// Setup mux & routes
	mux := http.NewServeMux()

	mux.HandleFunc("GET /entries/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Getting entries")
		entries, err := db.GetEntries()
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "error getting entries", http.StatusInternalServerError)
			return
		}
		renderJSON(w, entries)
	})

	mux.HandleFunc("GET /entry/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		logger.Info("Getting entry", "id", id)
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid id", http.StatusInternalServerError)
			return
		}
		e, err := db.GetEntry(idInt)
		if err != nil {
			http.Error(w, "entry not found", http.StatusNotFound)
			return
		}
		renderJSON(w, e)
	})

	mux.HandleFunc("POST /entry/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Creating entry")
		body, err := io.ReadAll(r.Body)
		var e database.Entry
		if err != nil {
			http.Error(w, "error reading body", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &e)
		logger.Debug("Creating entry", "entry", e)
		if err != nil {
			http.Error(w, "error parsing body", http.StatusInternalServerError)
			return
		}
		if err = db.InsertEntry(e); err != nil {
			http.Error(w, "error inserting record into database", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("DELETE /entry/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		logger.Info("Deleting entry", "id", id)
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid id", http.StatusNotAcceptable)
			return
		}
		if err = db.DeleteEntry(idInt); err != nil {
			logger.Error(err.Error())
			http.Error(w, "error deleting", http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("PUT /entry/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Updating entry")
		body, err := io.ReadAll(r.Body)
		var e database.Entry
		if err != nil {
			http.Error(w, "error reading body", http.StatusInternalServerError)
			return
		}
		err = json.Unmarshal(body, &e)
		if err != nil {
			http.Error(w, "error parsing body", http.StatusInternalServerError)
			return
		}
		if err = db.UpdateEntry(e); err != nil {
			logger.Error(err.Error())
			http.Error(w, "error updating", http.StatusInternalServerError)
		}
	})

	logger.Info("Starting server")
	if err = http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
		panic(err)
	}
}
