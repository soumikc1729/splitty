package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *App) Routes() http.Handler {
	router := httprouter.New()
	return router
}
