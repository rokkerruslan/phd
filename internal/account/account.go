package account

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID    int
	Email string

	password string
}

func New(ctx context.Context, email, password string) (*Account, error) {
	baseErr := "account.New fails: %v"

	a := Account{
		Email:    email,
		password: password,
	}

	err := a.create(ctx)
	if err != nil {
		return nil, fmt.Errorf(baseErr, err)
	}

	return &a, nil
}

func (a Account) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.password), append([]byte(password), options.GlobalSalt...))
}
