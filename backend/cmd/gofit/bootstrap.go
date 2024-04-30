package main

import (
	"time"

	"github.com/bamaas/gofit/internal/data"
	"github.com/google/uuid"
)

func (app *application) injectUser(email string, password string) error {
	user := &data.User{
		Email: email,
		Activated: true,
	}
	user.Password.Set(password)
	return app.models.Users.Insert(user)
}

func (app *application) injectSampleData() error {

	var checkIns []data.CheckIn

	for i := 1; i <= 60; i++ {
		uuid, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		checkIn := data.CheckIn{
			UUID:     uuid.String(),
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
	err := app.injectUser("hi@gofit.nl", "gofit123")
	if err != nil {
		return err
	}
	err = app.injectSampleData()
	if err != nil {
		return err
	}
	return nil
}