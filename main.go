package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

type entry struct {
	ID     int     `json:"id"`
	Weight float64 `json:"weight"`
}

var entries = []entry{
	{ID: 1, Weight: 1.0},
	{ID: 2, Weight: 2.0},
	{ID: 3, Weight: 3.0},
}

func getEntries() []entry {
	return entries
}

func getEntry(id int) (entry, bool) {
	for i, v := range entries {
		if v.ID == id {
			return entries[i], true
		}
	}
	return entry{}, false
}

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

var logLevel map[string]slog.Level = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

func main() {

	// Setup logger
	level := logLevel["INFO"]
	envLogLevel := os.Getenv("LOGLEVEL")
	if envLogLevel != "" {
		_, ok := logLevel[envLogLevel]
		if ok {
			level = logLevel[envLogLevel]
		} else {
			slog.Info("Invalid log level, falling back to default", "level", level)
		}
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))

	// Setup mux & routes
	mux := http.NewServeMux()
	mux.HandleFunc("GET /entries/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Getting entries")
		entries = getEntries()
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
		e, ok := getEntry(idInt)
		if !ok {
			http.Error(w, "entry not found", http.StatusNotFound)
			return
		}
		renderJSON(w, e)
	})

	mux.HandleFunc("POST /entry/", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Creating entry")
		body, err := io.ReadAll(r.Body)
		var e entry
		if err != nil {
			http.Error(w, "error reading body", http.StatusInternalServerError)
		}
		err = json.Unmarshal(body, &e)
		logger.Debug("Creating entry", "entry", e)
		if err != nil {
			http.Error(w, "error parsing body", http.StatusInternalServerError)
		}
		entries = append(entries, e)
	})

	mux.HandleFunc("DELETE /entry/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		logger.Info("Deleting entry", "id", id)
		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "invalid id", http.StatusInternalServerError)
			return
		}
		for i, v := range entries {
			if v.ID == idInt {
				entries = append(entries[:i], entries[i+1:]...)
				return
			}
		}
		http.Error(w, "entry not found", http.StatusNotFound)
	})

	logger.Info("Starting server")
	http.ListenAndServe("localhost:8080", mux)
}
