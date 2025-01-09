package main

import (
	"fmt"
	"net/http"

	"github.com/soumikc1729/splitty/server/internal/util"
)

func (app *App) LogError(r *http.Request, err error) {
	app.Logger.Err(err).Str("request_method", r.Method).Str("request_url", r.URL.String()).Msg("an error occurred")
}

func (app *App) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := util.Envelope{"error": message}
	err := util.WriteJSON(w, status, env, nil)
	if err != nil {
		app.LogError(r, err)
		w.WriteHeader(500)
	}
}

func (app *App) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.ErrorResponse(w, r, http.StatusBadRequest, "the requested resource could not be found")
}

func (app *App) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *App) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *App) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *App) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
}
