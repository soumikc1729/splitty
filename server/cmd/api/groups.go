package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/soumikc1729/splitty/server/internal/data"
	"github.com/soumikc1729/splitty/server/internal/util"
	"github.com/soumikc1729/splitty/server/internal/validator"
)

func (app *App) CreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string   `json:"name"`
		Users []string `json:"users"`
	}

	err := util.ReadJSON(r, &input)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	group := &data.Group{Name: input.Name, Users: input.Users}

	v := validator.New()

	if data.ValidateGroup(v, group); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.Data.Groups.Insert(group, app.Config.Data.QueryTimeout)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("/v1/groups/%d", group.ID))

	err = util.WriteJSON(w, http.StatusCreated, util.Envelope{"group": group}, header)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	app.Logger.Info().Int64("id", group.ID).Msg("created group")
}

func (app *App) GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	group := app.ContextGetGroup(r)

	if err := util.WriteJSON(w, http.StatusOK, util.Envelope{"group": group}, nil); err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	app.Logger.Info().Int64("id", group.ID).Msg("retrieved group")
}

func (app *App) UpdateGroupHandler(w http.ResponseWriter, r *http.Request) {
	group := app.ContextGetGroup(r)

	var input struct {
		Name  string   `json:"name"`
		Users []string `json:"users"`
	}

	if err := util.ReadJSON(r, &input); err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	group.Name = input.Name

	for _, user := range group.Users {
		v.Check(validator.In(user, input.Users...), "users", fmt.Sprintf("cannot remove user '%s'", user))
	}

	if !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	group.Users = input.Users

	if data.ValidateGroup(v, group); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.Data.Groups.Update(group, app.Config.Data.QueryTimeout); err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.EditConflictResponse(w, r)
		default:
			app.ServerErrorResponse(w, r, err)
		}
		return
	}

	if err := util.WriteJSON(w, http.StatusOK, util.Envelope{"group": group}, nil); err != nil {
		app.ServerErrorResponse(w, r, err)
	}

	app.Logger.Info().Int64("id", group.ID).Msg("updated group")
}

func (app *App) DeleteGroupHandler(w http.ResponseWriter, r *http.Request) {
	group := app.ContextGetGroup(r)

	err := app.Data.Groups.Delete(group.ID, group.Token, app.Config.Data.QueryTimeout)
	if err != nil {
		app.DataErrorResponse(w, r, err)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"message": "group successfully deleted"}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
	}

	app.Logger.Info().Int64("id", group.ID).Msg("deleted group")
}
