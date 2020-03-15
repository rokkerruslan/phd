package accounts

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
)

const createQuery = `
	INSERT INTO accounts (email, password, created, updated) 
		VALUES($1, $2, NOW(), NOW())
	RETURNING id
`

// TODO: all create model functions MUST return id or full object?
func (app *app) createAccount(ctx context.Context, a Account) (int, error) {
	var id int
	err := app.resources.Db.
		QueryRow(ctx, createQuery, a.Email, a.password).
		Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

var ErrAccountDoesNotExist = errors.New("account does not exist")

func (app *app) RetrieveByEmail(ctx context.Context, email string) (a Account, err error) {
	baseErr := "accounts.RetrieveByEmail fails: %v"

	err = app.resources.Db.
		QueryRow(ctx, "SELECT id, email, password FROM accounts WHERE email = $1", email).
		Scan(&a.ID, &a.Email, &a.password)
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
		QueryRow(ctx, "SELECT id, email FROM accounts WHERE id = $1", id).
		Scan(&a.ID, &a.Email); err != nil {
		return a, fmt.Errorf(baseErr, err)
	}
	return a, nil
}
