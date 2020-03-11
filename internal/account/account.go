package account

import (
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID    int
	Email string

	password string
}

func (a Account) CheckPassword(password string, globalSalt []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(a.password), append([]byte(password), globalSalt...))
}
