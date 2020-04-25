package accounts

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

var ErrAlreadyExists = errors.New("account already exists")

// TODO: all create model functions MUST return id or full object?
func (app *app) createAccount(ctx context.Context, a Account) (int, error) {
	baseErr := "createAccount fails: %w"

	var id int
	err := app.resources.Db.
		QueryRow(
			ctx,
			"INSERT INTO accounts (name, email, password, created, updated) VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id",
			a.Name,
			a.Email,
			a.password,
		).Scan(&id)

	// TODO: write helpers?
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.ConstraintName == "accounts_email_key" {
				return 0, fmt.Errorf(baseErr, ErrAlreadyExists)
			}
		}

		return 0, fmt.Errorf(baseErr, err)
	}
	return id, nil
}

var ErrAccountDoesNotExist = errors.New("account does not exist")

func (app *app) RetrieveByEmail(ctx context.Context, email string) (a Account, err error) {
	baseErr := "accounts.RetrieveByEmail fails: %v"

	err = app.resources.Db.
		QueryRow(
			ctx,
			"SELECT id, name, email, password FROM accounts WHERE email = $1",
			strings.ToLower(email),
		).Scan(&a.ID, &a.Name, &a.Email, &a.password)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			err = ErrAccountDoesNotExist
		}
		return a, fmt.Errorf(baseErr, err)
	}
	return a, nil
}

// TODO: do not use param like part of result object, always scan?
func (app *app) RetrieveByID(ctx context.Context, id int) (Account, error) {
	baseErr := "accounts.RetrieveByID fails: %v"

	var a Account
	if err := app.resources.Db.
		QueryRow(
			ctx,
			"SELECT id, name, email FROM accounts WHERE id = $1",
			id,
		).Scan(&a.ID, &a.Name, &a.Email); err != nil {
			return a, fmt.Errorf(baseErr, err)
	}
	return a, nil
}
