package main

import (
	"net/http"

	"github.com/bamaas/gofit/internal/data"
	"github.com/bamaas/gofit/internal/validator"
)

func (app * application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	// Get the input
	var input struct {
		Email string 	`json:"email"`
		Password string `json:"password"`
		Goal string 	`json:"goal"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &data.User{
		Email: input.Email,
		Activated: false,
	}

	// Hash the password
	err = user.Password.Set(input.Password)
	if err != nil {
		app.logger.Error(err.Error())
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "error creating user"})
		return
	}

	// Validate the user
	v := validator.New()
	v.ValidateUser(user)
	if !v.Valid() {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": v.Errors})
		return
	}

	// Insert the user
	err = app.models.Users.Insert(user)
	if err != nil {
		app.logger.Error(err.Error())
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "error inserting record"})
		return
	}

	// Respond
	app.writeJSON(w, http.StatusCreated, envelope{"data": user})
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	app.writeJSON(w, http.StatusOK, envelope{"data": user})
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	// Get the input
	var input struct {
		Goal string 	`json:"goal"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "failed to read input"})
		return
	}

	// Validate the input
	validGoals := map[string]bool{
		"cut": true,
		"bulk": true,
		"maintain": true,
	}
	v := validator.New()
	v.Check(validGoals[input.Goal], "goal", "must be 'cut', 'bulk', or 'maintain'");
	if !v.Valid() {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": v.Errors})
		return
	}

	// Update the user
	user.Goal = input.Goal
	err = app.models.Users.Update(user)
	if err != nil {
		app.logger.Error(err.Error())
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "error updating record"})
		return
	}

	// Respond
	app.writeJSON(w, http.StatusOK, envelope{"data": user})
}

// func (app *application) setUserGoalHandler(w http.ResponseWriter, r *http.Request) {
// 	user := app.contextGetUser(r)

// 	// Get the input
// 	var input struct {
// 		Goal string `json:"goal"`
// 	}

// 	err := app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "failed to read input"})
// 		return
// 	}

// 	// Validate the input
// 	validGoals := map[string]bool{
// 		"cut": true,
// 		"bulk": true,
// 		"maintain": true,
// 	}
// 	v := validator.New()
// 	v.Check(validGoals[input.Goal], "goal", "must be 'cut', 'bulk', or 'maintain'");
// 	if !v.Valid() {
// 		app.writeJSON(w, http.StatusBadRequest, envelope{"error": v.Errors})
// 		return
// 	}

// 	// Update the user
// 	user.Goal = input.Goal

// 	err = app.models.Users.Update(user)
// 	if err != nil {
// 		app.logger.Error(err.Error())
// 		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "error updating record"})
// 		return
// 	}

// 	// Respond
// 	app.writeJSON(w, http.StatusOK, envelope{"data": user})
// }