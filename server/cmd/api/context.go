package main

import (
	"context"
	"net/http"

	"github.com/soumikc1729/splitty/server/internal/data"
)

type ContextKey string

const (
	GroupContextKey ContextKey = "group"
)

func (app *App) ContextSetGroup(r *http.Request, group *data.Group) *http.Request {
	ctx := context.WithValue(r.Context(), GroupContextKey, group)
	return r.WithContext(ctx)
}

func (app *App) ContextGetGroup(r *http.Request) *data.Group {
	group, ok := r.Context().Value(GroupContextKey).(*data.Group)
	if !ok {
		panic("missing group value in request context")
	}

	return group
}
