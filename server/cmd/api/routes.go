package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *App) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.MethodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/v1/groups", app.CreateGroupHandler)

	return router
}
