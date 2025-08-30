package main

import (
	"context"
	"net/http"

	"github.com/grglucastr/go-greenlight/internal/data"
)

type contextKey string

// Use this constant as the key for get and set user info
// in the request context
const userContextKey = contextKey("user")

// returns a new copy of the request with the provided
// User struct added to the context.
// We are using the userContextKey constant as the key.
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)

	if !ok {
		panic("missing user value in request context")
	}

	return user
}
