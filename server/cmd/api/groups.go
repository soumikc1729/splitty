package main

import (
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
