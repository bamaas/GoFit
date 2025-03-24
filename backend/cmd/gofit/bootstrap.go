package main

import (
	"errors"
	"time"

	"github.com/bamaas/gofit/internal/data"
	"github.com/google/uuid"
)

func (app *application) injectUser(user *data.User) error {

	// Check if user already exists
	_, err := app.models.Users.GetByEmail(user.Email)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				return app.models.Users.Insert(user)
			default:
				return err
		}
	}
	return nil
}

func (app *application) injectSampleData() error {

	var checkIns []data.CheckIn

	for i := 1; i <= 120; i++ {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		checkIn := data.CheckIn{
			UUID:     uuid.String(),
			UserID:   1,
			Datetime: time.Now().AddDate(0, 0, -i).Unix(),
			Weight:   float64(i + 29),
			Notes:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam nulla sem.",
		}
		checkIns = append(checkIns, checkIn)
	}

	for _, c := range checkIns {
		err := app.models.CheckIns.Insert(c)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *application) Bootstrap() error {

	// Inject users
	for _, u := range app.config.Users {
		user := data.User{
			Email:     u.Email,
			Goal: 	   "cut",
		}
		err := user.Password.Set(u.Password)
		if err != nil {
			return err
		}
		err = app.injectUser(&user)
		if err != nil {
			return err
		}
	}

	// Inject sample data if running in development mode
	if app.config.DevelopmentMode {
		err := app.injectSampleData()
		if err != nil {
			return err
		}
	}

	return nil
}
