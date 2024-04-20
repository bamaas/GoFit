package validator

import (
	"github.com/bamaas/gofit/internal/data"
)

func (v *Validator) ValidateCheckIn(c *data.CheckIn) {
	v.Check(c.Datetime != 0, "datetime", "must be provided")
	v.Check(c.Weight >= 20, "weight", "must be at least 20")
	v.Check(c.Weight <= 300, "weight", "must be at most 300")
	v.Check(len(c.Notes) <= 100, "notes", "can't be longer than 100 characters")
}