package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"errors"
	"io"
)

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// Limit the size of the request body to 1MB
	var maxBytes int64 = 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	// Initialize the json decoder
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode the request body to the provided interface
	err := decoder.Decode(v)
	if err != nil {
		return err
	}

	// Ensure that there are no additional fields in the request body
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("request body must only contain a single JSON object")
	}

	return nil
}

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