package db

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"time"

	"github.com/lib/pq"
	"github.com/soumikc1729/splitty/server/internal/validator"
)

var (
	ShortTextRX = regexp.MustCompile(`^[a-zA-Z0-9 \-_]{3,50}$`)
)

type Group struct {
	ID      int64    `json:"id"`
	Name    string   `json:"name"`
	Token   string   `json:"token"`
	Users   []string `json:"users"`
	Version int      `json:"version"`
}

func ValidateGroup(v *validator.Validator, group *Group) {
	v.Check(validator.Matches(group.Name, ShortTextRX), "name", "must be 3-50 characters long and contain only letters, numbers, spaces, hyphens, and underscores")

	ValidateToken(v, group.Token)

	v.Check(validator.Unique(group.Users), "users", "must not contain duplicate values")
	v.Check(len(group.Users) >= 2, "users", "must contain at least two values")
	v.Check(len(group.Users) <= 50, "users", "must contain at most fifty values")
	v.Check(validator.All(group.Users, func(u string) bool {
		return validator.Matches(u, ShortTextRX)
	}), "users", "each value must be 3-50 characters long and contain only letters, numbers, spaces, hyphens, and underscores")
}

type GroupModel struct {
	Config *Config
	DB     *sql.DB
}

func (m GroupModel) Insert(group *Group, timeout time.Duration) error {
	query := `
		INSERT INTO groups (name, token, users)
		VALUES ($1, $2, $3)
		RETURNING id, version`

	retryCount := 3

	for range retryCount {
		token := GenerateRandomToken()
		args := []interface{}{group.Name, token, pq.Array(group.Users)}

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		err := m.DB.QueryRowContext(ctx, query, args...).Scan(&group.ID, &group.Version)
		if err != nil {
			switch {
			case err.Error() == `pq: duplicate key value violates unique constraint "groups_token_key"`:
				continue
			default:
				return err
			}
		}

		group.Token = token
		return nil
	}

	return ErrCannotGenerateUniqueToken
}

func (m GroupModel) GetByIDAndToken(id int64, token string, timeout time.Duration) (*Group, error) {
	var group Group

	query := `
		SELECT id, name, token, users, version
		FROM groups
		WHERE id = $1 AND token = $2`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	row := m.DB.QueryRowContext(ctx, query, id, token)
	err := row.Scan(&group.ID, &group.Name, &group.Token, pq.Array(&group.Users), &group.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &group, nil
}

func (m GroupModel) Update(group *Group, timeout time.Duration) error {
	query := `
		UPDATE groups
		SET name = $1, users = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version`

	args := []interface{}{group.Name, pq.Array(group.Users), group.ID, group.Version}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&group.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
	}

	return nil
}

func (m GroupModel) Delete(id int64, token string, timeout time.Duration) error {
	query := `
		DELETE FROM groups
		WHERE id = $1 AND token = $2`

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id, token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
