package main

import (
	"math"
	"math/rand"
	"net/http"
)

func (app *application) getStatsHandler(w http.ResponseWriter, r *http.Request) {

	var stats = make(map[string]map[string]float64)
	weight_diff := map[string]float64{
		"week_ago":    math.Floor(rand.Float64()*100)/100,
		"90_days_ago": float64((rand.Intn(4 - 2 + 1) + 2)),
		"all_time":    float64((rand.Intn(11 - 4 + 1) + 4)),
	}
	stats["weight_diff"] = weight_diff

	app.writeJSON(w, http.StatusOK, envelope{"stats": stats})
}