package accounts

import (
	"context"
	"fmt"
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

const selectQuery = `
	SELECT id, password FROM accounts WHERE email = $1
`

func (app *app) RetrieveByEmail(ctx context.Context, email string) (Account, error) {
	baseErr := "accounts.RetrieveByEmail fails: %v"

	var a Account
	if err := app.resources.Db.
		QueryRow(ctx, selectQuery, email).
		Scan(&a.ID, &a.password); err != nil {
		return a, fmt.Errorf(baseErr, err)
	}
	a.Email = email
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
