package account

import (
	"context"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID    int
	Email string

	password string
}

func New(ctx context.Context, email, password string) (Account, error) {
	a := Account{
		Email:    email,
		password: password,
	}

	return a, nil
}

func (a Account) CheckPassword(password string, globalSalt []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(a.password), append([]byte(password), globalSalt...))
}
