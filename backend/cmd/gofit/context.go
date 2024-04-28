package main

import (
	"context"
	"net/http"

	"github.com/bamaas/gofit/internal/data"
)

type contextKey string
const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		return nil
		// TODO: panic here?
	}
	return user
}