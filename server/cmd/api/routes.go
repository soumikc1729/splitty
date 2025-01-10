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
	router.HandlerFunc(http.MethodGet, "/v1/groups/:id", app.GetGroupHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/groups/:id", app.UpdateGroupHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/groups/:id", app.DeleteGroupHandler)

	return router
}
