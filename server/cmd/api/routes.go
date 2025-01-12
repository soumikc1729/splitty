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
	router.HandlerFunc(http.MethodGet, "/v1/groups/:groupID", app.AuthenticateGroup(app.GetGroupHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/groups/:groupID", app.AuthenticateGroup(app.UpdateGroupHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/groups/:groupID", app.AuthenticateGroup(app.DeleteGroupHandler))

	return router
}
