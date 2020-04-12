package tokens

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"ph/internal/api"
)

type (
	Assets struct {
		Db *pgxpool.Pool
	}
)

type App struct {
	Assets
}

func NewApp(assets Assets) *App {
	return &App{
		Assets: assets,
	}
}

// RetrieveAccountIDFromRequest fetch auth token from request headers and try
// fetch account id from tokens storage.
func (app *App) RetrieveAccountIDFromRequest(ctx context.Context, r *http.Request) (accountID int, err error) {
	baseErr := "RetrieveAccountIDFromRequest fails: %w"

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

const insertQuery = `
	INSERT INTO tokens (token, account_id, created) VALUES ($1, $2, NOW())
`

func (app *App) Create(ctx context.Context, accountID int) (session string, err error) {
	baseErr := "tokens.Create fails: %v"

	buf := make([]byte, 32)
	_, err = rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf(baseErr, err)
	}

	session = base64.StdEncoding.EncodeToString(buf)

	if _, err = app.Assets.Db.Exec(ctx, insertQuery, session, accountID); err != nil {
		return "", fmt.Errorf(baseErr, err)
	}

	return session, nil
}

var ErrDoesNotExist = errors.New("token doesn't exist")

func (app *App) retrieveAccountID(ctx context.Context, token string) (id int, err error) {
	baseErr := "retrieveAccountID fails: %w"

	if err = app.Db.QueryRow(ctx, "SELECT account_id FROM tokens WHERE token = $1", token).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = fmt.Errorf(baseErr, ErrDoesNotExist)
		} else {
			err = fmt.Errorf(baseErr, err)
		}
		return 0, err
	}

	return id, nil
}

type Session struct {
	Token     string
	AccountID int
}

func fromRequest(r *http.Request) (string, error) {
	token := r.Header.Get(api.AuthTokenHeaderName)
	if token == "" {
		return token, fmt.Errorf("`%s` isn't set", api.AuthTokenHeaderName)
	}
	return token, nil
}

func DropToken(ctx context.Context, db *pgxpool.Pool, token string) {
	baseErr := "token.DropToken fails: %v"

	_, err := db.Exec(ctx, "DELETE FROM tokens WHERE token = $1", token)
	if err != nil {
		log.Printf(baseErr, err)
	}
}
