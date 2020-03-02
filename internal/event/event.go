package event

import (
	"context"
	"fmt"
	"log"
	"strings"
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

func (e *Event) Validate() error {
	var errors []string

	if e.Name == "" {
		errors = append(errors, "`Name` can't be empty")
	}
	if e.OwnerID == 0 {
		errors = append(errors, "`OwnerID` can't be empty")
	}
	if e.Timelines == nil || len(e.Timelines) == 0 {
		errors = append(errors, "`Timelines` can't be empty")
	}
	for _, timeline := range e.Timelines {
		if tError := timeline.Validate(); tError != nil {
			errors = append(errors, tError.Error())
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf("event.Validate fails %v", strings.Join(errors, ", "))
}

func (e *Event) Create(ctx context.Context) error {
	return e.Insert(ctx)
}

func (e *Event) Insert(ctx context.Context) error {
	if _, err := db.Exec(ctx, "insert into events (name, created, updated) values($1, now(), now())", e.Name); err != nil {
		log.Println(err)
	}
	return nil
}

type Filter struct {
}

// todo: use filter
func ModelList(ctx context.Context, _ Filter) ([]Event, error) {
	defErr := "event.List fails: %v"

	rows, err := db.Query(ctx, "select id, name, created, updated from events")
	if err != nil {
		return nil, fmt.Errorf(defErr, err)
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var id int
		var name string
		var created time.Time
		var updated time.Time
		if err := rows.Scan(&id, &name, &created, &updated); err != nil {
			return nil, fmt.Errorf(defErr, err)
		}
		events = append(events, Event{
			ID:      id,
			Name:    name,
			Created: created,
			Updated: updated,
		})
	}

	return events, nil
}
