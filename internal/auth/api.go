package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"photo/internal/account"
	"photo/internal/api"
	"photo/internal/session"
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
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	// TODO: hash password anyway
	acc, err := account.RetrieveAccountByEmail(r.Context(), signData.Email)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	if err := acc.CheckPassword(signData.Password); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	// Authenticate user
	token, err := session.Create(r.Context(), db, acc.ID)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(signUpResponse{
		Token: token,
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
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}
	if len(signData.Password) < minimalPasswordLength {
		api.Error(w, fmt.Errorf(baseErr, "`Password` length check fails"), http.StatusBadRequest)
		return
	}
	if signData.Email == "" {
		api.Error(w, fmt.Errorf(baseErr, "`Email` is empty"), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword(append([]byte(signData.Password), options.GlobalSalt...), options.BcryptWorkFactor)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	acc, err := account.New(r.Context(), signData.Email, string(hash))
	if err != nil {
		// TODO: dispatch by type (for example duplicate email it's BadRequest)
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	// Authenticate
	token, err := session.Create(r.Context(), db, acc.ID)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(signUpResponse{
		Token: token,
	})
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	baseErr := "account.Retrieve fails: %v"

	token := r.Header.Get("X-Auth-Token")
	if token == "" {
		api.Error(w, fmt.Errorf(baseErr, fmt.Sprintf("`X-Auth-token` isn't set")), http.StatusForbidden)
		return
	}
	fmt.Println(token)
	session.DropSession(r.Context(), db, token)

	w.WriteHeader(http.StatusNoContent)
}
