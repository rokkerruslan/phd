package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// deprecated
type errorResponse struct {
	Error string
}

type errorResponse2 struct {
	Error struct{
		Code string
		Text string
	}
}

type apiError struct {
	StatusCode int
	ApiCode    string
	Text       string
}

const (
	// Common
	DatabaseError = iota
	NotExistError
	InvalidJsonError
	TemporaryError
	AuthHeaderError

	// Account
	PasswordLenError
	AccountAlreadyExist
)

var codes = map[int]apiError{
	DatabaseError: {
		ApiCode:    "PH100",
		StatusCode: http.StatusInternalServerError,
		Text:       "database error",
	},
	NotExistError: {
		ApiCode: "PH101",
		StatusCode: http.StatusBadRequest,
		Text: "entity does not exist",
	},
	InvalidJsonError: {
		ApiCode:    "PH102",
		StatusCode: http.StatusBadRequest,
		Text:       "invalid json",
	},
	AuthHeaderError: {
		ApiCode: "PH103",
		StatusCode: http.StatusForbidden,
		Text: "need auth",
	},
	TemporaryError: {
		ApiCode:    "PH142",
		StatusCode: http.StatusInternalServerError,
		Text:       "something went wrong",
	},

	PasswordLenError: {
		ApiCode:    "PH201",
		StatusCode: http.StatusBadRequest,
		Text:       "Password too small",
	},
	AccountAlreadyExist: {
		ApiCode:    "PH202",
		StatusCode: http.StatusBadRequest,
		Text:       "Account already exist",
	},
}

// deprecated
// TODO: logging here?
func Error(w http.ResponseWriter, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
}

func Error2(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	err, ok := codes[code]
	if !ok {
		panic(fmt.Sprintf("code %d does not exist", code))
	}
	w.WriteHeader(err.StatusCode)
	_ = json.NewEncoder(w).Encode(errorResponse2{Error: struct {
		Code string
		Text string
	}{Code: err.ApiCode, Text: err.Text}})
}
