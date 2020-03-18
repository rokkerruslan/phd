package accounts

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignUpInvalidInput(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))

	a := app{}
	a.signUpHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected: %v, got: %v", http.StatusBadRequest, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "invalid character") {
		t.Errorf("expected error with JSON, got: %v", body)
	}
}

func TestSignUpPasswordTooShort(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"Password":"TooShort!","Email":"e"}`))

	a := app{opts: Opts{
		MinLenForNewPassword: 10,
	}}
	a.signUpHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected: %v, got: %v", http.StatusBadRequest, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "`Password`") {
		t.Errorf("expected error with `Password`, got: %v", body)
	}
}

func TestSignUpEmailIsEmpty(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"Password":"VeryVeryVeryLong","Email":""}`))

	a := app{}
	a.signUpHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected: %v, got: %v", http.StatusBadRequest, rec.Code)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "`Email`") {
		t.Errorf("expected error with `Email`, got: %v", body)
	}
}
