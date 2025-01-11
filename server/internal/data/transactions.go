package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/soumikc1729/splitty/server/internal/validator"
)

type Payment struct {
	Amount float64 `json:"amount"`
	Payer  string  `json:"payer"`
}

type Transaction struct {
	ID       int64     `json:"id"`
	Title    string    `json:"title"`
	Payments []Payment `json:"payments"`
	GroupID  int64     `json:"group_id"`
	Version  int       `json:"-"`
}

func ValidateTransaction(v *validator.Validator, transaction *Transaction, group *Group) {
	v.Check(validator.Matches(transaction.Title, ShortTextRX), "title", "must be 3-50 characters long and contain only letters, numbers, spaces, hyphens, and underscores")

	var payers []string
	var amount float64
	for _, p := range transaction.Payments {
		v.Check(validator.In(p.Payer, group.Users...), "payments", fmt.Sprintf("%s not one of the group users", p.Payer))
		payers = append(payers, p.Payer)
		amount += p.Amount
	}
	v.Check(validator.Unique(payers), "payments", "must not contain duplicate payers")
	v.Check(amount == 0, "payments", "sum of all payments must be 0")

	v.Check(transaction.GroupID == group.ID, "group_id", "must be same as the id of the group")
}

type TransactionModel struct {
	DB *sql.DB
}

func (t *TransactionModel) Insert(transaction *Transaction, timeout time.Duration) error {
	query := `
        INSERT INTO transactions (title, payments, group_id)
        VALUES ($1, $2, $3)
        RETURNING id, version`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	args := []interface{}{transaction.Title, pq.Array(transaction.Payments), transaction.GroupID}

	return t.DB.QueryRowContext(ctx, query, args...).Scan(&transaction.ID, &transaction.Version)
}

func (t *TransactionModel) GetAllAfterID(id int64, groupID int64, timeout time.Duration) (*[]Transaction, error) {
	query := `
        SELECT id, title, payments, group_id, version
        FROM transactions
        WHERE id > $1 AND group_id = $2
        ORDER BY id ASC`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rows, err := t.DB.QueryContext(ctx, query, id, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		var transaction Transaction
		var payments []Payment

		err := rows.Scan(
			&transaction.ID,
			&transaction.Title,
			pq.Array(&payments),
			&transaction.GroupID,
			&transaction.Version,
		)

		if err != nil {
			return nil, err
		}

		transaction.Payments = payments
		transactions = append(transactions, transaction)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &transactions, nil
}

func (t *TransactionModel) Update(transaction *Transaction, timeout time.Duration) error {
	query := `
        UPDATE transactions
        SET title = $1, payments = $2, version = version + 1
        WHERE id = $3 AND group_id = $4 AND version = $5
        RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	args := []interface{}{transaction.Title, pq.Array(transaction.Payments), transaction.ID, transaction.GroupID, transaction.Version}

	err := t.DB.QueryRowContext(ctx, query, args...).Scan(&transaction.Version)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrEditConflict
		}
		return err
	}

	return nil
}

func (t *TransactionModel) Delete(id int64, groupID int64, timeout time.Duration) error {
	query := `
        DELETE FROM transactions
        WHERE id = $1 AND group_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := t.DB.ExecContext(ctx, query, id, groupID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
