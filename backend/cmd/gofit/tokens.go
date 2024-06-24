package main

import (
	"net/http"
	"time"

	"github.com/bamaas/gofit/internal/data"
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
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": "failed to read input"})
		return
	}

	// Validate the input
	v := validator.New()
	v.ValidateEmail(input.Email)
	v.ValidatePlainTextPassword(input.Password)
	if !v.Valid() {
		app.writeJSON(w, http.StatusBadRequest, envelope{"error": v.Errors})
	}

	// Get the user
	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		// TODO: implement better error handling
		app.writeJSON(w, http.StatusUnauthorized, envelope{"error": "invalid credentials"})
		return
	}

	// Check the password
	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.writeJSON(w, http.StatusUnauthorized, envelope{"error": "invalid credentials"})
		return
	}
	if !match {
		app.writeJSON(w, http.StatusUnauthorized, envelope{"error": "invalid credentials"})
		return
	}

	// Create the token
	token, err := app.models.Tokens.New(user.ID, 1*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.logger.Error(err.Error())
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "error creating token"})
		return
	}

	// Respond
	app.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token})
}

func (app *application) retractAllAuthenticationTokensHandler(w http.ResponseWriter, r *http.Request) {

	// Get the user
	user := app.contextGetUser(r)

	// Retract the token
	err := app.models.Tokens.DeleteAllForUser(data.ScopeAuthentication, user.ID)
	if err != nil {
		app.logger.Error(err.Error())
		app.writeJSON(w, http.StatusInternalServerError, envelope{"error": "error retracting tokens"})
		return
	}

	// Respond
	app.writeJSON(w, http.StatusOK, envelope{"message": "tokens retracted"})
}