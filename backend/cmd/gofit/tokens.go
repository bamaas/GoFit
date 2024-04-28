package main

import (
	"net/http"
	"time"

	"github.com/bamaas/gofit/internal/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {

	// Get the input
	var input struct {
		Email string 	`json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "failed to read input", http.StatusBadRequest)
		return
	}

	// Validate the input
	v := validator.New()
	v.ValidateEmail(input.Email)
	v.ValidatePlainTextPassword(input.Password)
	if !v.Valid() {
		http.Error(w, "failed validating input", http.StatusBadRequest)
	}

	// Get the user
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		// TODO: implement better error handling
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check the password
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if !match {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create the token
	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, "authentication")
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error creating token", http.StatusInternalServerError)
		return
	}

	// Respond
	app.writeJSON(w, http.StatusCreated, envelope{"token": token})
}