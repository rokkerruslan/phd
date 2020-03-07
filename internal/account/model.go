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

func RetrieveAccount(ctx context.Context, email string) (Account, error) {
	baseErr := "retrieveAccount fails: %v"

	var acc Account
	if err := db.QueryRow(ctx, selectQuery, email).Scan(&acc.ID, &acc.password); err != nil {
		return acc, fmt.Errorf(baseErr, err)
	}
	acc.Email = email
	return acc, nil
}
