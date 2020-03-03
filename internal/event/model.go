package event

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

const createQuery = `
	INSERT INTO events (name, description owner_id, created, updated) VALUES($1, $2, $3, NOW(), NOW())
`

func (e *Event) Create(ctx context.Context) error {
	if _, err := db.Exec(ctx, createQuery, e.Name, e.Description, e.OwnerID); err != nil {
		return err
	}
	return nil
}

const updateQuery = `
	UPDATE events SET name = $1, updated = NOW() WHERE id = $2
`

func (e *Event) Update(ctx context.Context) error {
	baseErr := "event.Update fails: %v"
	_, err := db.Exec(ctx, updateQuery, e.Name, e.ID)
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}
	return nil
}

const retrieveQuery = `
	SELECT id, name, owner_id, created, updated FROM events WHERE id = $1
`

func modelRetrieve(ctx context.Context, f filterRetrieve) (e Event, err error) {
	defErr := "modelRetrieve fails: %v"
	err = db.QueryRow(ctx, retrieveQuery, f.ID).Scan(e.ID, e.Name, e.OwnerID, e.Created, e.Updated)
	if err != nil {
		return e, fmt.Errorf(defErr, err)
	}
	return e, nil
}

func ModelList(ctx context.Context, _ Filter) ([]Event, error) {
	defErr := "event.List fails: %v"
	rows, err := db.Query(ctx, "SELECT id, name, owner_id, created, updated FROM events")
	if err != nil {
		return nil, fmt.Errorf(defErr, err)
	}
	defer rows.Close()
	return construct(rows)
}

func construct(rows pgx.Rows) (events []Event, err error) {
	baseErr := "construct fails: %v"

	for rows.Next() {
		var id int
		var name string
		var ownerID int
		var created time.Time
		var updated time.Time
		if err := rows.Scan(&id, &name, &ownerID, &created, &updated); err != nil {
			return nil, fmt.Errorf(baseErr, err)
		}
		events = append(events, Event{
			ID:      id,
			Name:    name,
			OwnerID: ownerID,
			Created: created,
			Updated: updated,
		})
	}

	return events, nil
}
