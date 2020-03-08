package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const insertQuery = `
	INSERT INTO sessions (session, account_id, created) VALUES ($1, $2, NOW())
`

func Create(ctx context.Context, db *pgxpool.Pool, accountID int) (session string, err error) {
	baseErr := "session.Create fails: %v"

	buf := make([]byte, 32)
	_, err = rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf(baseErr, err)
	}

	session = base64.StdEncoding.EncodeToString(buf)

	if _, err = db.Exec(ctx, insertQuery, session, accountID); err != nil {
		return "", fmt.Errorf(baseErr, err)
	}

	return session, nil
}

var ErrDoesNotExist = errors.New("token doesn't exist")

func Retrieve(ctx context.Context, db *pgxpool.Pool, session string) (id int, err error) {
	baseErr := "session.Retrieve fails: %w"

	if err = db.QueryRow(ctx, "SELECT account_id FROM sessions WHERE session = $1", session).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = fmt.Errorf(baseErr, ErrDoesNotExist)
		} else {
			err = fmt.Errorf(baseErr, err)
		}
		return 0, err
	}

	return id, nil
}

func DropSession(ctx context.Context, db *pgxpool.Pool, session string) {
	baseErr := "session.DropSession fails: %v"

	_, err := db.Exec(ctx, "DELETE FROM sessions WHERE session = $1", session)
	if err != nil {
		log.Printf(baseErr, err)
	}
}
