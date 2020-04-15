package offer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"
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

// Users can't create offers not from yourself.
func (o *Offer) CanBeCreated(accountID int) bool {
	return o.AccountID == accountID
}

func (app *App) createOffer(ctx context.Context, o Offer) (Offer, error) {
	row := app.assets.Db.QueryRow(
		ctx,
		"INSERT INTO offers (account_id, event_id, created, updated) VALUES ($1, $2, NOW(), NOW()) RETURNING id, created, updated",
		o.AccountID,
		o.EventID,
	)

	if err := row.Scan(&o.ID, &o.Created, &o.Updated); err != nil {
		return o, err
	}

	return o, nil
}

func (app *App) offerList(ctx context.Context, f Filter) ([]Offer, error) {
	baseErr := "accounts.offerList fails: %v"

	var err error
	var rows pgx.Rows
	if f.AccountID != 0 {
		rows, err = app.assets.Db.Query(
			ctx, "SELECT id, account_id, event_id, created, updated FROM offers WHERE account_id = $1", f.AccountID)
	} else if f.EventID != 0 {
		rows, err = app.assets.Db.Query(
			ctx, "SELECT id, account_ID, event_id, created, updated FROM offers WHERE event_id = $1", f.EventID)
	} else {
		return nil, fmt.Errorf(baseErr, "filter is empty")
	}

	if err != nil {
		return nil, fmt.Errorf(baseErr, err)
	}
	defer rows.Close()

	return constructOffersFromRows(rows)
}

func (app *App) offerFetchOne(ctx context.Context, offerID int) (Offer, error) {
	baseErr := "app.offerGetOne fails: %v"

	row := app.assets.Db.QueryRow(
		ctx,
		"SELECT id, account_id, event_id, created, updated FROM offers WHERE id = $1",
		offerID,
	)

	var o Offer
	if err := row.Scan(o.ID, o.AccountID, o.EventID, o.Created, o.Updated); err != nil {
		return o, fmt.Errorf(baseErr, err)
	}

	return o, nil
}

func constructOffersFromRows(rows pgx.Rows) ([]Offer, error) {
	baseErr := "constructOffersFromRows fails: %w"

	offers := []Offer{}
	for rows.Next() {
		var id int
		var accountID int
		var eventID int
		var created time.Time
		var updated time.Time
		if err := rows.Scan(&id, &accountID, &eventID, &created, &updated); err != nil {
			return nil, fmt.Errorf(baseErr, err)
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