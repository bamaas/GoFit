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
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	// Validate the user
	v := validator.New()
	v.ValidateUser(user)
	if !v.Valid() {
		http.Error(w, "failed validating user", http.StatusBadRequest)
		return
	}

	// Insert the user
	err = app.models.Users.Insert(user)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "error inserting record", http.StatusInternalServerError)
		return
	}

	// Respond
	app.writeJSON(w, http.StatusCreated, envelope{"user": user})
}