package validator

import (
	"github.com/bamaas/gofit/internal/data"
)

func (v *Validator) ValidateEmail(email string) {
	v.Check(email != "", "email", "email is required")
	v.Check(v.Matches(email, EmailRX), "email", "email is not valid")
}

func (v *Validator) ValidatePlainTextPassword(password string) {
	v.Check(password != "", "password", "password is required")
	v.Check(len(password) >= 8, "password", "password must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "password not be more than 72 bytes")
}

func (v *Validator) ValidateUser(user *data.User) {
	v.ValidateEmail(user.Email)
	v.ValidatePlainTextPassword(*user.Password.Plaintext)
	if user.Password.Hash == nil {
		panic("missing password hash")
	}
}