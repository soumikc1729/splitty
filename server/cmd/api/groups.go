package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
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
	id, err := readID(r)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	token, err := readToken(r)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateToken(v, token); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return
	}

	group, err := app.Data.Groups.GetByIDAndToken(id, token, app.Config.Data.QueryTimeout)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"group": group}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	app.Logger.Info().Int64("id", group.ID).Msg("retrieved group")
}

func readID(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func readToken(r *http.Request) (string, error) {
	token := r.URL.Query().Get("token")
	if token == "" {
		return "", errors.New("missing token parameter")
	}

	return token, nil
}
