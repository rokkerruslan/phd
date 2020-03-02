package offer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
)

// todo: use pool
var db *pgx.Conn

func init() {
	var err error
	db, err = pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:10003/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

type Offer struct {
	ID        int
	AccountID int
	EventID   int
	Created   time.Time
	Updated   time.Time
}

func (o *Offer) Validate() error {
	return nil
}

func (o *Offer) Insert(ctx context.Context) error {
	_, err := db.Exec(ctx, "insert into offers (account_id, event_id, created, updated) values ($1, $2, now(), now())", o.AccountID, o.EventID)
	if err != nil {
		return err
	}
	return nil
}

type Filter struct {
	AccountID int
}

func ModelList(ctx context.Context, f Filter) ([]Offer, error) {
	defErr := "offer.List fails: %v"

	rows, err := db.Query(ctx, "select id, account_id, event_id, created, updated from offers where account_id = $1", f.AccountID)
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
