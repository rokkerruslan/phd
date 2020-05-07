package tokens

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"ph/internal/api"
)

type (
	Assets struct {
		Db *pgxpool.Pool
	}
	Opts struct {
		TokenTTL time.Duration
	}
)

type App struct {
	Assets
	Opts
}

func NewApp(assets Assets, opts Opts) *App {
	return &App{
		Assets: assets,
		Opts:   opts,
	}
}

var ErrDoesNotExist = errors.New("token doesn't exist")

// RetrieveAccountIDFromRequest fetch auth token from request headers and try
// fetch account id from tokens storage.
func (app *App) RetrieveAccountIDFromRequest(ctx context.Context, r *http.Request) (accountID int, err error) {
	baseErr := "tokens.RetrieveAccountIDFromRequest fails: %w"

	token, err := fromRequest(r)
	if err != nil {
		return 0, fmt.Errorf(baseErr, err)
	}
	accountID, err = app.retrieveAccountID(ctx, token)
	if err != nil {
		return 0, fmt.Errorf(baseErr, err)
	}

	return accountID, nil
}

// Create instantiate new token for submitted account id.
func (app *App) Create(ctx context.Context, accountID int) (token string, err error) {
	baseErr := "tokens.Create fails: %v"

	buf := make([]byte, 32)
	_, err = rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf(baseErr, err)
	}

	token = base64.StdEncoding.EncodeToString(buf)

	if _, err = app.Assets.Db.Exec(
		ctx,
		"INSERT INTO tokens (token, account_id, created) VALUES ($1, $2, NOW())",
		token,
		accountID,
	); err != nil {
		return "", fmt.Errorf(baseErr, err)
	}

	return token, nil
}

func (app *App) DropToken(ctx context.Context, token string) {
	baseErr := "tokens.DropToken fails: %v"

	_, err := app.Db.Exec(ctx, "DELETE FROM tokens WHERE token = $1", token)
	if err != nil {
		log.Printf(baseErr, err)
	}
}

func (app *App) retrieveAccountID(ctx context.Context, token string) (id int, err error) {
	baseErr := "tokens.retrieveAccountID fails: %w"

	if err = app.Db.QueryRow(
		ctx,
		"SELECT account_id FROM tokens WHERE token = $1 AND created > $2",
		token, time.Now().Add(-app.TokenTTL)).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = fmt.Errorf(baseErr, ErrDoesNotExist)
		} else {
			err = fmt.Errorf(baseErr, err)
		}
		return 0, err
	}

	return id, nil
}

func fromRequest(r *http.Request) (string, error) {
	token := r.Header.Get(api.AuthTokenHeaderName)
	if token == "" {
		return token, fmt.Errorf("`%s` isn't set", api.AuthTokenHeaderName)
	}
	return token, nil
}
