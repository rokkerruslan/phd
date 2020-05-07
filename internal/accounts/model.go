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

func (app *App) createAccount(ctx context.Context, a Account) (Account, error) {
	baseErr := "createAccount fails: %w"

	err := app.assets.Db.
		QueryRow(
			ctx,
			"INSERT INTO accounts (name, email, password, is_deleted, created, updated) VALUES ($1, $2, $3, FALSE, NOW(), NOW()) RETURNING id, created, updated",
			a.Name,
			a.Email,
			a.password,
		).Scan(&a.ID, &a.Created, &a.Updated)

	// TODO: write helpers? 23505 - UNIQUE
	if err != nil {
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.Code == "23505" {
				return a, fmt.Errorf(baseErr, ErrAlreadyExists)
			}
		}
		return a, fmt.Errorf(baseErr, err)
	}

	return a, nil
}

func (app *App) deleteAccount(ctx context.Context, id int) error {
	baseErr := "deleteAccount fails: %v"

	if _, err := app.assets.Db.
		Exec(ctx, "UPDATE accounts SET is_deleted = TRUE WHERE id = $1", id); err != nil {
		return fmt.Errorf(baseErr, err)
	}

	if _, err := app.assets.Db.
		Exec(ctx, "DELETE FROM tokens WHERE account_id = $1", id); err != nil {
		return fmt.Errorf(baseErr, err)
	}

	return nil
}

var ErrAccountDoesNotExist = errors.New("account does not exist")

func (app *App) RetrieveByEmail(ctx context.Context, email string) (a Account, err error) {
	baseErr := "accounts.RetrieveByEmail fails: %v"

	err = app.assets.Db.
		QueryRow(
			ctx,
			"SELECT id, name, email, password FROM accounts WHERE email = $1 AND is_deleted = FALSE",
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
func (app *App) RetrieveByID(ctx context.Context, id int) (Account, error) {
	baseErr := "accounts.RetrieveByID fails: %v"

	var a Account
	if err := app.assets.Db.
		QueryRow(
			ctx,
			"SELECT id, name, email FROM accounts WHERE id = $1 AND is_deleted = FALSE",
			id,
		).Scan(&a.ID, &a.Name, &a.Email); err != nil {
		return a, fmt.Errorf(baseErr, err)
	}
	return a, nil
}
