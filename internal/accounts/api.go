package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"photo/internal/api"
	"photo/internal/tokens"
)

func (app *app) retrieveHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "accounts.retrieveHandler fails: %v"

	token, err := tokens.FromRequest(r)
	if err != nil {
		api.Error(w, err, http.StatusForbidden)
		return
	}

	id, err := tokens.Retrieve(r.Context(), app.resources.Db, token)
	if err != nil {
		switch {
		case errors.Is(err, tokens.ErrDoesNotExist):
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

func (app *app) deleteHandler(_ http.ResponseWriter, _ *http.Request) {

}

// Auth

type signInRequest struct {
	Email    string
	Password string
}

// TODO: can't sign if already have a token
func (app *app) signInHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "accounts.signInHandler fails: %v"

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

	if err := acc.CheckPassword(signData.Password, app.opts.GlobalSalt); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	// Authenticate user
	token, err := tokens.Create(r.Context(), app.resources.Db, acc.ID)
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
func (app *app) signUpHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "accounts.signUpHandler fails: %v"

	var signData signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&signData); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}
	if len(signData.Password) < app.opts.MinLenForNewPassword {
		api.Error(w, fmt.Errorf(baseErr, "`Password` length check fails"), http.StatusBadRequest)
		return
	}
	if signData.Email == "" {
		api.Error(w, fmt.Errorf(baseErr, "`Email` is empty"), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword(
		append([]byte(signData.Password), app.opts.GlobalSalt...), app.opts.BcryptWorkFactor)
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
	token, err := tokens.Create(r.Context(), app.resources.Db, acc.ID)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(signUpResponse{
		Token: token,
	})
}

func (app *app) signOutHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "accounts.signOutHandler fails: %v"

	token := r.Header.Get(api.AuthTokenHeaderName)
	if token == "" {
		api.Error(w, fmt.Errorf(baseErr,
			fmt.Sprintf("%s` isn't set", api.AuthTokenHeaderName)), http.StatusForbidden)
		return
	}
	tokens.DropToken(r.Context(), app.resources.Db, token)

	w.WriteHeader(http.StatusNoContent)
}
