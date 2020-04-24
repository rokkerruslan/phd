package accounts

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func NewAccount(email string, passwordHash string) Account {
	return Account{
		Email:    strings.ToLower(email),
		password: passwordHash,
	}
}

type Account struct {
	ID    int
	Email string

	password string
}

func (a Account) CheckPassword(password string, globalSalt []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(a.password), append([]byte(password), globalSalt...))
}
