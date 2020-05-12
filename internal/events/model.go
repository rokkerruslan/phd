package events

import (
	"context"
	"fmt"

	"ph/internal/api"

	"github.com/jackc/pgx/v4"
)

const createQuery = `
	INSERT INTO events (name, description, owner_id, created, updated, is_public)
		VALUES($1, $2, $3, NOW(), NOW(), $4)
	RETURNING id, created, updated
`

func (app *App) createEvent(ctx context.Context, e Event) (Event, error) {
	baseErr := "createEvent fails: %v"

	err := app.assets.Db.QueryRow(
		ctx,
		createQuery,
		e.Name,
		e.Description,
		e.OwnerID,
		e.IsPublic,
	).Scan(&e.ID, &e.Created, &e.Updated)
	if err != nil {
		return e, fmt.Errorf(baseErr, err)
	}

	for i, timeline := range e.Timelines {
		err := app.assets.Db.QueryRow(
			ctx,
			"INSERT INTO timelines (event_id, start, \"end\", place) VALUES ($1, $2, $3, $4) RETURNING id",
			e.ID,
			timeline.Start,
			timeline.End,
			timeline.Place,
		).Scan(&e.Timelines[i].ID)
		if err != nil {
			return e, fmt.Errorf(baseErr, err)
		}
	}

	return e, nil
}

const updateQuery = `
	UPDATE events SET name = $1, updated = NOW(), is_hidden = $3, description = $4 WHERE id = $2
`

func (app *App) updateEvent(ctx context.Context, e Event) error {
	baseErr := "event.Update fails: %v"

	_, err := app.assets.Db.Exec(ctx, updateQuery, e.Name, e.ID, e.IsHidden, e.Description)
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}
	return nil
}

const retrieveQuery = `
	SELECT id, name, owner_id, created, updated FROM events WHERE id = $1
`

func (app *App) retrieveEvent(ctx context.Context, f api.RetrieveFilter) (e Event, err error) {
	defErr := "modelRetrieve fails: %v"
	err = app.assets.Db.QueryRow(ctx, retrieveQuery, f.ID).Scan(e.ID, e.Name, e.OwnerID, e.Created, e.Updated)
	if err != nil {
		return e, fmt.Errorf(defErr, err)
	}
	return e, nil
}

const deleteQuery = `
	DELETE FROM events WHERE id = $1
`

func (app *App) deleteEvent(ctx context.Context, f api.RetrieveFilter) error {
	defErr := "modelDelete fails: %v"
	_, err := app.assets.Db.Exec(ctx, deleteQuery, f.ID)
	if err != nil {
		return fmt.Errorf(defErr, err)
	}
	return nil
}

const selectQuery = `
	SELECT id, name, description, owner_id, created, updated, is_public, is_hidden FROM events LIMIT 10
`

const selectTimelinesQuery = `
	SELECT id, start, "end", place, event_id FROM timelines WHERE event_id = ANY($1)
`

func (app *App) eventList(ctx context.Context, f api.AccountAndEventFilter) ([]Event, error) {
	baseErr := "event.listHandler fails: %v"

	var err error
	var rows pgx.Rows
	if f.AccountID != 0 {
		rows, err = app.assets.Db.Query(
			ctx,
			"SELECT id, name, description, owner_id, created, updated, is_public, is_hidden FROM events WHERE owner_id = $1",
			f.AccountID,
		)
	} else {
		return []Event{}, fmt.Errorf(baseErr, "filter is empty")
	}
	defer rows.Close()

	events, err := app.constructEventList(ctx, rows)
	if err != nil {
		return nil, fmt.Errorf(baseErr, err)
	}

	return events, nil
}

func (app *App) constructEventList(ctx context.Context, rows pgx.Rows) ([]Event, error) {
	events := []Event{}
	var eventIDs []int
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.Name, &e.Description, &e.OwnerID, &e.Created, &e.Updated, &e.IsPublic, &e.IsHidden); err != nil {
			return nil, err
		}
		eventIDs = append(eventIDs, e.ID)
		events = append(events, e)
	}

	timelineRows, err := app.assets.Db.Query(ctx, selectTimelinesQuery, eventIDs)
	if err != nil {
		return nil, err
	}
	timelines := make(map[int][]Timeline)
	for timelineRows.Next() {
		var t Timeline
		var eventID int
		if err := timelineRows.Scan(&t.ID, &t.Start, &t.End, &t.Place, &eventID); err != nil {
			return nil, err
		}
		timelines[eventID] = append(timelines[eventID], t)
	}

	for i, event := range events {
		events[i].Timelines = timelines[event.ID]
	}

	return events, nil
}
