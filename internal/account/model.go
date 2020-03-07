package account

import (
	"context"
	"fmt"
)

const createQuery = `
	INSERT INTO accounts (email, password, created, updated) VALUES($1, $2, NOW(), NOW()) RETURNING id
`

func (a *Account) create(ctx context.Context) error {
	err := db.QueryRow(ctx, createQuery, a.Email, a.password).Scan(&a.ID)
	if err != nil {
		return err
	}
	return nil
}

const selectQuery = `
	SELECT id, password FROM accounts WHERE email = $1
`

func RetrieveAccountByEmail(ctx context.Context, email string) (Account, error) {
	baseErr := "retrieveAccount fails: %v"

	var acc Account
	if err := db.QueryRow(ctx, selectQuery, email).Scan(&acc.ID, &acc.password); err != nil {
		return acc, fmt.Errorf(baseErr, err)
	}
	acc.Email = email
	return acc, nil
}

func RetrieveByID(ctx context.Context, id int) (Account, error) {
	baseErr := "account.RetrieveByID fails: %v"

	a := Account{ID: id}
	if err := db.QueryRow(ctx, "SELECT email FROM accounts WHERE id = $1", id).Scan(&a.Email); err != nil {
		return a, fmt.Errorf(baseErr, err)
	}
	return a, nil
}
