package validator

import (
	"github.com/bamaas/gofit/internal/data"
)

func (v *Validator) ValidateUser(user *data.User) {

	// Check email
	v.Check(user.Email != "", "email", "email is required")
	v.Check(v.Matches(user.Email, EmailRX), "email", "email is not valid")

	// Check password
	v.Check(user.Password.Plaintext != nil, "password", "password is required")
	v.Check(len(*user.Password.Plaintext) >= 8, "password", "password must be at least 8 bytes long")
	v.Check(len(*user.Password.Plaintext) <= 72, "password", "password not be more than 72 bytes")
	if user.Password.Hash == nil {
		panic("missing password hash")
	}
}