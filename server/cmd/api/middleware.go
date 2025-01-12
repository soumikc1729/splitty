package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/soumikc1729/splitty/server/internal/data"
	"github.com/soumikc1729/splitty/server/internal/validator"
)

func (app *App) AuthenticateGroup(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.NotFoundResponse(w, r)
			default:
				app.ServerErrorResponse(w, r, err)
			}
			return
		}

		r = app.ContextSetGroup(r, group)
		next.ServeHTTP(w, r)
	}
}

func readID(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("groupID"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func readToken(r *http.Request) (string, error) {
	token := r.Header.Get("X-Group-Token")
	if token == "" {
		return "", errors.New("invalid token header")
	}

	return token, nil
}
