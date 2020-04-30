package accounts

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"ph/internal/api"
	"ph/internal/tokens"
)

// TODO: check id from path
func (app *app) retrieveHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "retrieveHandler fails: %v"

	id, err := app.tokens.RetrieveAccountIDFromRequest(r.Context(), r)
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

	api.Response(w, account)
}

// TODO: check account id from path
func (app *app) deleteHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "deleteHandler fails: %v"

	id, err := app.tokens.RetrieveAccountIDFromRequest(r.Context(), r)
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

	if err := app.deleteAccount(r.Context(), id); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Auth

type signInRequest struct {
	Email    string
	Password string
}

type signInResponse struct {
	AccountID int
	Token     string
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
	a, err := app.RetrieveByEmail(r.Context(), signData.Email)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	if err := a.CheckPassword(signData.Password, app.opts.GlobalSalt); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	// Authenticate user
	token, err := app.tokens.Create(r.Context(), a.ID)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	api.Response(w, signInResponse{
		Token:     token,
		AccountID: a.ID,
	})
}

type signUpRequest struct {
	Name     string
	Email    string
	Password string
}

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

type signUpResponse struct {
	Token   string
	Account Account
}

// TODO: disable if already have a token.
// TODO: signUpHandler MUST return account id
func (app *app) signUpHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "accounts.signUpHandler fails: %v"

	var signData signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&signData); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	if err := signData.Validate(app.opts.MinLenForNewPassword); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword(
		append([]byte(signData.Password), app.opts.GlobalSalt...), app.opts.BcryptWorkFactor)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	a := NewAccount(signData.Name, signData.Email, string(hash))

	if a, err = app.createAccount(r.Context(), a); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrAlreadyExists) {
			status = http.StatusBadRequest
		}
		api.Error(w, fmt.Errorf(baseErr, err), status)
		return
	}

	// Authenticate
	token, err := app.tokens.Create(r.Context(), a.ID)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	api.Response(w, signUpResponse{
		Token:   token,
		Account: a,
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
	app.tokens.DropToken(r.Context(), token)

	w.WriteHeader(http.StatusNoContent)
}
