package accounts

import (
	"fmt"
	"strings"
)

func (r *signUpRequest) Validate(passwordMinLen int) error {
	var e []string

	if r.Name == "" {
		e = append(e, "`Name` is empty")
	}
	if r.Email == "" {
		e = append(e, "`Email` is empty")
	}
	if len(r.Password) < passwordMinLen {
		e = append(e, "`Password` length check fails")
	}

	if len(e) != 0 {
		return fmt.Errorf("signUpRequest.Validate fails: %v", strings.Join(e, ", "))
	}

	return nil
}
