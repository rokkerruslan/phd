package event

import (
	"context"
	"fmt"
	"time"
)

const createQuery = `
	INSERT INTO events (name, description, owner_id, created, updated, is_public) VALUES($1, $2, $3, NOW(), NOW(), $4)
`

func (app *app) CreateEvent(ctx context.Context, e Event) error {
	if _, err := app.resources.Db.Exec(ctx, createQuery, e.Name, e.Description, e.OwnerID, e.IsPublic); err != nil {
		return err
	}
	return nil
}

const updateQuery = `
	UPDATE events SET name = $1, updated = NOW() WHERE id = $2
`

func (app *app) updateEvent(ctx context.Context, e Event) error {
	baseErr := "event.Update fails: %v"
	_, err := app.resources.Db.Exec(ctx, updateQuery, e.Name, e.ID)
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}
	return nil
}

const retrieveQuery = `
	SELECT id, name, owner_id, created, updated FROM events WHERE id = $1
`

func (app *app) retrieveEvent(ctx context.Context, f filterRetrieve) (e Event, err error) {
	defErr := "modelRetrieve fails: %v"
	err = app.resources.Db.QueryRow(ctx, retrieveQuery, f.ID).Scan(e.ID, e.Name, e.OwnerID, e.Created, e.Updated)
	if err != nil {
		return e, fmt.Errorf(defErr, err)
	}
	return e, nil
}

const deleteQuery = `
	DELETE FROM events WHERE id == $1
`

func (app *app) deleteEvent(ctx context.Context, f filterRetrieve) error {
	defErr := "modelDelete fails: %v"
	_, err := app.resources.Db.Exec(ctx, deleteQuery, f.ID)
	if err != nil {
		return fmt.Errorf(defErr, err)
	}
	return nil
}

const selectQuery = `
	SELECT id, name, description, owner_id, created, updated FROM events LIMIT 10
`

const selectTimelinesQuery = `
	SELECT id, start, "end", place, event_id FROM timelines WHERE event_id = ANY($1)
`

func (app *app) eventList(ctx context.Context, _ Filter) ([]Event, error) {
	baseErr := "event.eventListHandler fails: %v"
	rows, err := app.resources.Db.Query(ctx, selectQuery)
	if err != nil {
		return nil, fmt.Errorf(baseErr, err)
	}
	defer rows.Close()

	var events []Event
	var eventIDs []int
	for rows.Next() {
		var id int
		var name string
		var description string
		var ownerID int
		var created time.Time
		var updated time.Time
		if err := rows.Scan(&id, &name, &description, &ownerID, &created, &updated); err != nil {
			return nil, fmt.Errorf(baseErr, err)
		}
		eventIDs = append(eventIDs, id)
		events = append(events, Event{
			ID:          id,
			Name:        name,
			Description: description,
			OwnerID:     ownerID,
			Created:     created,
			Updated:     updated,
		})
	}

	timelineRows, err := app.resources.Db.Query(ctx, selectTimelinesQuery, eventIDs)
	if err != nil {
		return nil, fmt.Errorf(baseErr, err)
	}
	timelines := make(map[int][]Timeline)
	for timelineRows.Next() {
		var id int
		var start time.Time
		var end time.Time
		var place string
		var eventID int
		if err := timelineRows.Scan(&id, &start, &end, &place, &eventID); err != nil {
			return nil, fmt.Errorf(baseErr, err)
		}
		timelines[eventID] = append(timelines[eventID], Timeline{
			ID:    id,
			Start: start,
			End:   end,
			Place: place,
		})
	}

	for i, event := range events {
		events[i].Timelines = timelines[event.ID]
	}

	return events, nil
}
