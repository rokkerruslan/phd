package offer

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"photo/internal/errors"
)

func Create(w http.ResponseWriter, r *http.Request) {
	var offer Offer
	if err := json.NewDecoder(r.Body).Decode(&offer); err != nil {
		log.Println(err)
		return
	}

	if err := offer.Validate(); err != nil {
		errors.APIError(w, err, http.StatusBadRequest)
		return
	}

	if err := offer.Insert(r.Context()); err != nil {
		errors.APIError(w, err, http.StatusInternalServerError)
		return
	}
}

func List(w http.ResponseWriter, r *http.Request) {
	var err error
	var filter Filter
	if filter.AccountID, err = strconv.Atoi(r.URL.Query().Get("account_id")); err != nil {
		errors.APIError(w, fmt.Errorf("account_id parsing fails: %v", err), http.StatusBadRequest)
		return
	}

	offers, err := ModelList(r.Context(), filter)
	if err != nil {
		errors.APIError(w, err, http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(offers)
}
