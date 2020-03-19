package offer

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Offer struct {
	ID        int
	AccountID int
	EventID   int
	Created   time.Time
	Updated   time.Time
}

func (o *Offer) ValidateForCreate() error {
	baseErr := "offer.ValidateForCreate fails: %v"

	var errors []string
	if o.AccountID == 0 {
		errors = append(errors, "offer.AccountID is zero")
	}
	if o.EventID == 0 {
		errors = append(errors, "offer.EventID is zero")
	}

	if len(errors) != 0 {
		return fmt.Errorf(baseErr, strings.Join(errors, ", "))
	}

	return nil
}

func (app *app) createOffer(ctx context.Context, o Offer) error {
	_, err := app.resources.Db.Exec(
		ctx,
		"INSERT INTO offers (account_id, event_id, created, updated) VALUES ($1, $2, NOW(), NOW())",
		o.AccountID,
		o.EventID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (app *app) offerList(ctx context.Context, f Filter) ([]Offer, error) {
	defErr := "offer.List fails: %v"

	rows, err := app.resources.Db.Query(
		ctx, "SELECT id, account_id, event_id, created, updated FROM offers WHERE account_id = $1", f.AccountID)
	if err != nil {
		return nil, fmt.Errorf(defErr, err)
	}
	defer rows.Close()

	offers := []Offer{}
	for rows.Next() {
		var id int
		var accountID int
		var eventID int
		var created time.Time
		var updated time.Time
		if err := rows.Scan(&id, &accountID, &eventID, &created, &updated); err != nil {
			return nil, fmt.Errorf(defErr, err)
		}
		offers = append(offers, Offer{
			ID:        id,
			AccountID: accountID,
			EventID:   eventID,
			Created:   created,
			Updated:   updated,
		})
	}

	return offers, nil
}
