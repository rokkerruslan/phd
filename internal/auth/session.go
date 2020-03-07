package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const insertQuery = `
	INSERT INTO sessions (session, account_id, created) VALUES ($1, $2, NOW())
`

func createSession(ctx context.Context, accountID int) (session string, err error) {
	baseErr := "createSession fails: %v"

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
