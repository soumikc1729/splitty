package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/soumikc1729/splitty/server/internal/data"
	"github.com/soumikc1729/splitty/server/internal/util"
	"github.com/soumikc1729/splitty/server/internal/validator"
)

func (app *App) CreateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	group := app.ContextGetGroup(r)
	if transaction := app.validateTransactionInput(w, r, group); transaction != nil {
		if err := app.Data.Transactions.Insert(transaction, app.Config.Data.QueryTimeout); err != nil {
			app.ServerErrorResponse(w, r, err)
			return
		}

		header := make(http.Header)
		header.Set("Location", fmt.Sprintf("/v1/transactions/%d", transaction.ID))

		if err := util.WriteJSON(w, http.StatusCreated, util.Envelope{"transaction": transaction}, header); err != nil {
			app.ServerErrorResponse(w, r, err)
			return
		}

		app.Logger.Info().Int64("group-id", group.ID).Int64("transaction-id", transaction.ID).Msg("created transaction")
	}
}

func (app *App) ListTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	group := app.ContextGetGroup(r)

	after, err := strconv.ParseInt(r.URL.Query().Get("after"), 10, 64)
	if err != nil {
		app.Logger.Info().Msg("using 0 (default) since no valid value for after specified")
		after = 0
	}

	transactions, err := app.Data.Transactions.GetAllAfterID(after, group.ID, app.Config.Data.QueryTimeout)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	err = util.WriteJSON(w, http.StatusOK, util.Envelope{"transactions": transactions}, nil)
	if err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	app.Logger.Info().Int64("group-id", group.ID).Int64("transactions-after", after).Msg("retrieved transactions")
}

func (app *App) UpdateTransactionHandler(w http.ResponseWriter, r *http.Request) {
	group := app.ContextGetGroup(r)
	if updatedTransaction := app.validateTransactionInput(w, r, group); updatedTransaction != nil {
		id, err := util.ReadParam("transactionID", r)
		if err != nil {
			app.BadRequestResponse(w, r, err)
			return
		}

		transaction, err := app.Data.Transactions.Get(id, group.ID, app.Config.Data.QueryTimeout)
		if err != nil {
			app.DataErrorResponse(w, r, err)
			return
		}

		updatedTransaction.ID = id
		updatedTransaction.Version = transaction.Version

		if err := app.Data.Transactions.Update(updatedTransaction, app.Config.Data.QueryTimeout); err != nil {
			app.DataErrorResponse(w, r, err)
			return
		}

		if err := util.WriteJSON(w, http.StatusOK, util.Envelope{"transaction": updatedTransaction}, nil); err != nil {
			app.ServerErrorResponse(w, r, err)
			return
		}

		app.Logger.Info().Int64("group-id", group.ID).Int64("transaction-id", transaction.ID).Msg("updated transaction")
	}
}

func (app *App) DeleteTransactionHandler(w http.ResponseWriter, r *http.Request) {
	group := app.ContextGetGroup(r)

	id, err := util.ReadParam("transactionID", r)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return
	}

	if err = app.Data.Transactions.Delete(id, group.ID, app.Config.Data.QueryTimeout); err != nil {
		app.DataErrorResponse(w, r, err)
		return
	}

	if err := util.WriteJSON(w, http.StatusOK, util.Envelope{"message": "transaction successfully deleted"}, nil); err != nil {
		app.ServerErrorResponse(w, r, err)
		return
	}

	app.Logger.Info().Int64("group-id", group.ID).Int64("transaction-id", id).Msg("deleted transaction")
}

func (app *App) validateTransactionInput(w http.ResponseWriter, r *http.Request, group *data.Group) *data.Transaction {
	var input struct {
		Title    string `json:"title"`
		Payments []struct {
			Amount float64 `json:"amount"`
			Payer  string  `json:"payer"`
		} `json:"payments"`
	}

	err := util.ReadJSON(r, &input)
	if err != nil {
		app.BadRequestResponse(w, r, err)
		return nil
	}

	payments := []data.Payment{}
	for _, p := range input.Payments {
		payments = append(payments, data.Payment{
			Amount: p.Amount,
			Payer:  p.Payer,
		})
	}
	transaction := &data.Transaction{
		Title:    input.Title,
		Payments: payments,
		GroupID:  group.ID,
	}

	v := validator.New()

	if data.ValidateTransaction(v, transaction, group); !v.Valid() {
		app.FailedValidationResponse(w, r, v.Errors)
		return nil
	}

	return transaction
}
