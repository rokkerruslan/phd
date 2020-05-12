package offer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"

	"ph/internal/api"
)

type Offer struct {
	ID         int
	AccountID  int
	EventID    int
	IsApproved bool
	Created    time.Time
	Updated    time.Time
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
	// TODO: check event for is_hidden, we can't create offer for hidden events.

	if len(errors) != 0 {
		return fmt.Errorf(baseErr, strings.Join(errors, ", "))
	}

	return nil
}

// Users can't create offers not from yourself.
func (o *Offer) canBeCreated(accountID int) bool {
	return o.AccountID == accountID
}

func (app *App) createOffer(ctx context.Context, o Offer) (Offer, error) {
	baseErr := "createOffer fails: %v"

	var isPublic bool
	if err := app.assets.Db.QueryRow(ctx, "SELECT is_public FROM events WHERE id = $1", o.EventID).
		Scan(&isPublic); err != nil {
		return o, fmt.Errorf(baseErr, err)
	}

	o.IsApproved = isPublic

	row := app.assets.Db.QueryRow(
		ctx,
		"INSERT INTO offers (account_id, event_id, created, updated, is_approved) VALUES ($1, $2, NOW(), NOW(), $3) RETURNING id, created, updated",
		o.AccountID,
		o.EventID,
		o.IsApproved,
	)

	if err := row.Scan(&o.ID, &o.Created, &o.Updated); err != nil {
		return o, fmt.Errorf(baseErr, err)
	}

	return o, nil
}

func (app *App) updateOffer(ctx context.Context, o Offer) (Offer, error) {
	baseErr := "updateOffer fails: %v"

	row := app.assets.Db.QueryRow(ctx, "UPDATE offers SET is_approved = $1 WHERE id = $2 RETURNING account_id, event_id, created, updated", o.IsApproved, o.ID)

	if err := row.Scan(&o.AccountID, &o.EventID, &o.Created, &o.Updated); err != nil {
		return o, fmt.Errorf(baseErr, err)
	}

	return o, nil
}

func (app *App) deleteOffer(ctx context.Context, f api.RetrieveFilter) error {
	baseErr := "deleteOffer fails: %v"

	if _, err := app.assets.Db.Exec(ctx, "DELETE FROM offers WHERE id = $1", f.ID); err != nil {
		return fmt.Errorf(baseErr, err)
	}

	return nil
}

func (app *App) offerList(ctx context.Context, f api.AccountAndEventFilter) ([]Offer, error) {
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
		var o Offer
		if err := rows.Scan(&o.ID, &o.AccountID, &o.EventID, &o.Created, &o.Updated); err != nil {
			return nil, fmt.Errorf(baseErr, err)
		}
		offers = append(offers, o)
	}

	return offers, nil
}
