package offer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"ph/internal/api"
)

func (app *app) create(w http.ResponseWriter, r *http.Request) {
	var offer Offer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		log.Println(err)
		return
	}

	if err := offer.Validate(); err != nil {
		api.Error(w, err, http.StatusBadRequest)
		return
	}

	if err := app.insert(r.Context(), offer); err != nil {
		api.Error(w, err, http.StatusInternalServerError)
		return
	}
}

type Filter struct {
	AccountID int
}

func NewFilterFromQuery(values url.Values) (f Filter, err error) {
	defErr := "filter creating fails: %v"
	if f.AccountID, err = strconv.Atoi(values.Get("account_id")); err != nil {
		return f, fmt.Errorf(defErr, fmt.Errorf("account_id parsing fails: %v", err))
	}
	return f, nil
}

func (app *app) list(w http.ResponseWriter, r *http.Request) {
	filter, err := NewFilterFromQuery(r.URL.Query())
	if err != nil {
		api.Error(w, fmt.Errorf("account_id parsing fails: %v", err), http.StatusBadRequest)
		return
	}

	offers, err := app.offerList(r.Context(), filter)
	if err != nil {
		api.Error(w, err, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(offers)
}
