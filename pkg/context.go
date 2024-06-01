package pkg

import (
	"context"
	"damir/internal/entity"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *Applicaiton) contextSetUser(r *http.Request, user *entity.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *Applicaiton) contextGetUser(r *http.Request) *entity.User {
	user, ok := r.Context().Value(userContextKey).(*entity.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
