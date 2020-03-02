package event

import (
	"context"
	"fmt"
	"time"

	"photo/internal/geo"
)

type Event struct {
	ID        int
	Name      string
	OwnerID   int
	Created   time.Time
	Updated   time.Time
	Timelines []Timeline
	Point     geo.Point
}

func (e *Event) Create(ctx context.Context) error {
	return e.Insert(ctx)
}

const InsertQuery = `insert into events (name, owner_id, created, updated) values($1, $2, now(), now())`

func (e *Event) Insert(ctx context.Context) error {
	if _, err := db.Exec(ctx, InsertQuery, e.Name, e.OwnerID); err != nil {
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


func ModelList(ctx context.Context, _ Filter) ([]Event, error) {
	defErr := "event.List fails: %v"

	rows, err := db.Query(ctx, "SELECT id, name, owner_id, created, updated FROM events")
	if err != nil {
		return nil, fmt.Errorf(defErr, err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var id int
		var name string
		var ownerID int
		var created time.Time
		var updated time.Time
		if err := rows.Scan(&id, &name, &ownerID, &created, &updated); err != nil {
			return nil, fmt.Errorf(defErr, err)
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
