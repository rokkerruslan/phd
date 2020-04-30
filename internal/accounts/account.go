package accounts

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func NewAccount(name, email string, passwordHash string) Account {
	return Account{
		Name:     name,
		Email:    strings.ToLower(email),
		password: passwordHash,
	}
}

type Account struct {
	ID    int
	Name  string
	Email string

	Created time.Time
	Updated time.Time

	password string
}

func (a Account) CheckPassword(password string, globalSalt []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(a.password), append([]byte(password), globalSalt...))
}
