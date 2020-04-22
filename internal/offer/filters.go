package offer

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type ListFilter struct {
	AccountID int
	EventID   int
}

func NewListFilterFromQuery(values url.Values) (ListFilter, error) {
	baseErr := "NewListFilterFromQuery fails: %v"

	var f ListFilter
	var err error
	var errors []string

	accountIDRaw := values.Get("account_id")
	if accountIDRaw != "" {
		if f.AccountID, err = strconv.Atoi(accountIDRaw); err != nil {
			errors = append(errors, fmt.Sprintf("account_id parsing fails: %v", err))
		}
	}

	eventIDRaw := values.Get("event_id")
	if eventIDRaw != "" {
		if f.EventID, err = strconv.Atoi(eventIDRaw); err != nil {
			errors = append(errors, fmt.Sprintf("event_id parsing fails: %v", err))
		}
	}

	if len(errors) != 0 {
		return f, fmt.Errorf(baseErr, strings.Join(errors, ", "))
	}

	return f, nil
}
