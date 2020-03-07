package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"photo/internal/account"
	"photo/internal/errors"
)

const minimalPasswordLength = 10

type signInRequest struct {
	Email    string
	Password string
}

// TODO: can't sign if already have a token
func SignIn(w http.ResponseWriter, r *http.Request) {
	baseErr := "auth.SignIn fails: %v"

	var signData signInRequest
	if err := json.NewDecoder(r.Body).Decode(&signData); err != nil {
		errors.APIError(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	// TODO: hash password anyway
	acc, err := account.RetrieveAccount(r.Context(), signData.Email)
	if err != nil {
		errors.APIError(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	if err := acc.CheckPassword(signData.Password); err != nil {
		errors.APIError(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	// Authenticate user
	session, err := createSession(r.Context(), acc.ID)
	if err != nil {
		errors.APIError(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(signUpResponse{
		Token: session,
	})
}

type signUpRequest struct {
	Email    string
	Password string
}

type signUpResponse struct {
	Token string
}

// TODO: disable if already have a token.
func SignUp(w http.ResponseWriter, r *http.Request) {
	baseErr := "auth.SignUp fails: %v"

	var signData signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&signData); err != nil {
		errors.APIError(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}
	if len(signData.Password) < minimalPasswordLength {
		errors.APIError(w, fmt.Errorf(baseErr, "`Password` length check fails"), http.StatusBadRequest)
		return
	}
	if signData.Email == "" {
		errors.APIError(w, fmt.Errorf(baseErr, "`Email` is empty"), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword(append([]byte(signData.Password), options.GlobalSalt...), options.BcryptWorkFactor)
	if err != nil {
		errors.APIError(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	acc, err := account.New(r.Context(), signData.Email, string(hash))
	if err != nil {
		// TODO: dispatch by type (for example duplicate email it's BadRequest)
		errors.APIError(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	// Authenticate
	session, err := createSession(r.Context(), acc.ID)
	if err != nil {
		errors.APIError(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(signUpResponse{
		Token: session,
	})
}

func SignOut(_ http.ResponseWriter, _ *http.Request) {

}
