package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

func (app *application) readInt(qs url.Values, key string, defaultValue int) int {
	value := qs.Get(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}