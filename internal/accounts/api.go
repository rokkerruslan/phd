package accounts

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"ph/internal/api"
	"ph/internal/tokens"

	"golang.org/x/crypto/bcrypt"
)

// TODO: check id from path
func (app *App) retrieveHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "retrieveHandler fails: %v"

	id, err := app.assets.Tokens.RetrieveAccountIDFromRequest(r.Context(), r)
	if err != nil {
		log.Printf(baseErr, err)
		switch {
		case errors.Is(err, tokens.AuthHeaderError):
			api.Error2(w, api.AuthHeaderError)
		case errors.Is(err, tokens.ErrDoesNotExist):
			api.Error2(w, api.NotExistError)
		default:
			api.Error2(w, api.DatabaseError)
		}
		return
	}

	// todo (rr): we need one query for retrieve account info
	account, err := app.RetrieveByID(r.Context(), id)
	if err != nil {
		log.Printf(baseErr, err)
		api.Error2(w, api.DatabaseError)
		return
	}

	api.Response(w, account)
}

// TODO: check account id from path
func (app *App) deleteHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "deleteHandler fails: %v"

	id, err := app.assets.Tokens.RetrieveAccountIDFromRequest(r.Context(), r)
	if err != nil {
		log.Printf(baseErr, err)
		switch {
		case errors.Is(err, tokens.AuthHeaderError):
			api.Error2(w, api.AuthHeaderError)
		case errors.Is(err, tokens.ErrDoesNotExist):
			api.Error2(w, api.NotExistError)
		default:
			api.Error2(w, api.DatabaseError)
		}
		return
	}

	if err := app.deleteAccount(r.Context(), id); err != nil {
		api.Error2(w, api.DatabaseError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type signInRequest struct {
	Email    string
	Password string
}

// TODO: can't sign if already have a token
func (app *App) signInHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "accounts.signInHandler fails: %v"

	var signData signInRequest
	if err := json.NewDecoder(r.Body).Decode(&signData); err != nil {
		log.Printf(baseErr, err)
		api.Error2(w, api.InvalidJsonError)
		return
	}

	// TODO: hash password anyway
	a, err := app.RetrieveByEmail(r.Context(), signData.Email)
	if err != nil {
		log.Printf(baseErr, err)
		switch {
		case errors.Is(err, ErrAccountDoesNotExist):
			api.Error2(w, api.NotExistError)
		default:
			api.Error2(w, api.DatabaseError)
		}
		return
	}

	if err := a.CheckPassword(signData.Password, app.opts.GlobalSalt); err != nil {
		log.Printf(baseErr, err)
		api.Error2(w, api.NotExistError)
		return
	}

	// Authenticate user
	token, err := app.assets.Tokens.Create(r.Context(), a.ID)
	if err != nil {
		log.Printf(baseErr, err)
		api.Error2(w, api.DatabaseError)
		return
	}

	api.Response(w, signUpResponse{
		Token:   token,
		Account: a,
	})
}

type signUpRequest struct {
	Name     string
	Email    string
	Password string
}

type signUpResponse struct {
	Token   string
	Account Account
}

// TODO: disable if already have a token.
// TODO: signUpHandler MUST return account id
func (app *App) signUpHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "accounts.signUpHandler fails: %v"

	var signData signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&signData); err != nil {
		log.Printf(baseErr, err)
		api.Error2(w, api.InvalidJsonError)
		return
	}

	if err := signData.Validate(app.opts.MinLenForNewPassword); err != nil {
		log.Printf(baseErr, err)
		api.Error2(w, api.PasswordLenError)
		return
	}

	hash, err := bcrypt.GenerateFromPassword(
		append([]byte(signData.Password), app.opts.GlobalSalt...), app.opts.BcryptWorkFactor)
	if err != nil {
		log.Printf(baseErr, err)
		api.Error2(w, api.TemporaryError)
		return
	}

	a := NewAccount(signData.Name, signData.Email, string(hash))

	if a, err = app.createAccount(r.Context(), a); err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			api.Error2(w, api.AccountAlreadyExist)
		} else {
			api.Error2(w, api.DatabaseError)
		}
		return
	}

	// Authenticate
	token, err := app.assets.Tokens.Create(r.Context(), a.ID)
	if err != nil {
		log.Printf(baseErr, err)
		api.Error2(w, api.DatabaseError)
		return
	}

	api.Response(w, signUpResponse{
		Token:   token,
		Account: a,
	})
}

func (app *App) signOutHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "accounts.signOutHandler fails: %v"

	token := r.Header.Get(api.AuthTokenHeaderName)
	if token == "" {
		log.Printf(baseErr, "token does not exist")
		api.Error2(w, api.AuthHeaderError)
		return
	}
	app.assets.Tokens.DropToken(r.Context(), token)

	w.WriteHeader(http.StatusNoContent)
}
