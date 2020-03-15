package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"photo/internal/api"
	"photo/internal/session"
)

func (app *app) retrieveHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "account.retrieveHandler fails: %v"

	token := r.Header.Get(api.AuthTokenHeaderName)
	if token == "" {
		api.Error(w, fmt.Errorf(baseErr,
			fmt.Sprintf("`%s` isn't set", api.AuthTokenHeaderName)), http.StatusForbidden)
		return
	}

	id, err := session.Retrieve(r.Context(), app.resources.Db, token)
	if err != nil {
		switch {
		case errors.Is(err, session.ErrDoesNotExist):
			api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
			return
		default:
			api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
			return
		}
	}

	account, err := app.RetrieveByID(r.Context(), id)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(account)
}

func (app *app) Delete(_ http.ResponseWriter, _ *http.Request) {

}

// Auth

const minimalPasswordLength = 10

type signInRequest struct {
	Email    string
	Password string
}

// TODO: can't sign if already have a token
func (app *app) SignIn(w http.ResponseWriter, r *http.Request) {
	baseErr := "auth.SignIn fails: %v"

	var signData signInRequest
	if err := json.NewDecoder(r.Body).Decode(&signData); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	// TODO: hash password anyway
	acc, err := app.RetrieveByEmail(r.Context(), signData.Email)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	if err := acc.CheckPassword(signData.Password, app.options.GlobalSalt); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	// Authenticate user
	token, err := session.Create(r.Context(), app.resources.Db, acc.ID)
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
func (app *app) SignUp(w http.ResponseWriter, r *http.Request) {
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

	hash, err := bcrypt.GenerateFromPassword(
		append([]byte(signData.Password), app.options.GlobalSalt...), app.options.BcryptWorkFactor)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	acc := Account{
		Email:    signData.Email,
		password: string(hash),
	}

	if acc.ID, err = app.createAccount(r.Context(), acc); err != nil {
		// TODO: dispatch by type (for example duplicate email it's BadRequest)
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	// Authenticate
	token, err := session.Create(r.Context(), app.resources.Db, acc.ID)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(signUpResponse{
		Token: token,
	})
}

func (app *app) SignOut(w http.ResponseWriter, r *http.Request) {
	baseErr := "account.retrieve fails: %v"

	token := r.Header.Get(api.AuthTokenHeaderName)
	if token == "" {
		api.Error(w, fmt.Errorf(baseErr,
			fmt.Sprintf("%s` isn't set", api.AuthTokenHeaderName)), http.StatusForbidden)
		return
	}
	session.DropSession(r.Context(), app.resources.Db, token)

	w.WriteHeader(http.StatusNoContent)
}