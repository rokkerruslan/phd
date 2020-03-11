package account

import (
	"context"
	"fmt"
)

const createQuery = `
	INSERT INTO accounts (email, password, created, updated) VALUES($1, $2, NOW(), NOW()) RETURNING id
`

func (app App) createAccount(ctx context.Context, a Account) (int, error) {
	var id int
	err := app.resources.Db.QueryRow(ctx, createQuery, a.Email, a.password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

const selectQuery = `
	SELECT id, password FROM accounts WHERE email = $1
`

func (app *App) RetrieveAccountByEmail(ctx context.Context, email string) (Account, error) {
	baseErr := "retrieveAccount fails: %v"

	var acc Account
	if err := app.resources.Db.QueryRow(ctx, selectQuery, email).Scan(&acc.ID, &acc.password); err != nil {
		return acc, fmt.Errorf(baseErr, err)
	}
	acc.Email = email
	return acc, nil
}

func (app *App) RetrieveByID(ctx context.Context, id int) (Account, error) {
	baseErr := "account.RetrieveByID fails: %v"

	a := Account{ID: id}
	if err := app.resources.Db.QueryRow(ctx, "SELECT email FROM accounts WHERE id = $1", id).Scan(&a.Email); err != nil {
		return a, fmt.Errorf(baseErr, err)
	}
	return a, nil
}
