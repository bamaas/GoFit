package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
	"io"
)

type entry struct {
	ID int `json:"id"`
	Weight float64 `json:"weight"`
}

var entries = []entry{
	{ID: 1, Weight: 1.0},
	{ID: 2, Weight: 2.0},
	{ID: 3, Weight: 3.0},
}

func getEntries() []entry {
	fmt.Println("Getting entries")
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

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /entries/", func(w http.ResponseWriter, r *http.Request) {
		entries = getEntries()
		renderJSON(w, entries)
	})

	mux.HandleFunc("GET /entry/{id}/", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "handling entry with id=%v\n", id)
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
		fmt.Fprintf(w, "handling POST entry\n")
		body, err := io.ReadAll(r.Body)
		var e entry
		if err != nil {
			http.Error(w, "error reading body", http.StatusInternalServerError)
		}
		err = json.Unmarshal(body, &e)
		if err != nil {
			http.Error(w, "error parsing body", http.StatusInternalServerError)
		}
		entries = append(entries, e)
	})

	http.ListenAndServe("localhost:8080", mux)
  }
