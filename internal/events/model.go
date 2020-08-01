package events

import (
	"context"
	"fmt"

	"ph/internal/api"

	"github.com/jackc/pgx/v4"
)

func (app *App) eventCreate(ctx context.Context, e Event) (Event, error) {
	baseErr := "eventCreate fails: %v"

	q := `
		INSERT INTO events (name, description, owner_id, created, updated, is_public)
			VALUES($1, $2, $3, NOW(), NOW(), $4)
		RETURNING id, created, updated
	`

	if err := app.assets.Db.QueryRow(
		ctx,
		q,
		e.Name,
		e.Description,
		e.OwnerID,
		e.IsPublic,
	).Scan(&e.ID, &e.Created, &e.Updated); err != nil {
		return e, fmt.Errorf(baseErr, err)
	}

	for i, timeline := range e.Timelines {
		err := app.assets.Db.QueryRow(
			ctx,
			`INSERT INTO timelines (event_id, start, "end", place) VALUES ($1, $2, $3, $4) RETURNING id`,
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

func (app *App) eventUpdate(ctx context.Context, e Event) (Event, error) {
	baseErr := "eventUpdate fails: %v"

	if err := app.assets.Db.QueryRow(
		ctx,
		`UPDATE events SET name = $1, updated = NOW(), is_hidden = $3, description = $4 WHERE id = $2 RETURNING created, updated, is_public, is_hidden`,
		e.Name,
		e.ID,
		e.IsHidden,
		e.Description,
	).Scan(&e.Created, &e.Updated, &e.IsPublic, &e.IsHidden); err != nil {
		return e, fmt.Errorf(baseErr, err)
	}

	var err error
	if e.Timelines, err = app.fetchTimelines(ctx, e.ID); err != nil {
		return e, fmt.Errorf(baseErr, err)
	}

	return e, nil
}

func (app *App) eventRetrieve(ctx context.Context, f api.RetrieveFilter) (e Event, err error) {
	baseErr := "eventRetrieve fails: %v"

	if err = app.assets.Db.QueryRow(
		ctx,
		"SELECT id, name, owner_id, created, description, updated FROM events WHERE id = $1",
		f.ID,
	).Scan(&e.ID, &e.Name, &e.OwnerID, &e.Created, &e.Description, &e.Updated); err != nil {
		return e, fmt.Errorf(baseErr, err)
	}
	return e, nil
}

func (app *App) eventDelete(ctx context.Context, eventID int) error {
	baseErr := "eventDelete fails: %v"

	if _, err := app.assets.Db.Exec(
		ctx,
		"UPDATE events SET is_deleted = TRUE WHERE id = $1",
		eventID,
	); err != nil {
		return fmt.Errorf(baseErr, err)
	}
	return nil
}

func (app *App) eventList(ctx context.Context, f api.AccountAndEventFilter) ([]Event, error) {
	baseErr := "eventList fails: %v"

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

	timelineRows, err := app.assets.Db.Query(
		ctx,
		`SELECT id, start, "end", place, event_id FROM timelines WHERE event_id = ANY($1)`,
		eventIDs,
	)
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

func (app *App) fetchTimelines(ctx context.Context, eventID int) ([]Timeline, error) {
	timelineRows, err := app.assets.Db.Query(
		ctx,
		`SELECT id, start, "end", place, event_id FROM timelines WHERE event_id = $1`,
		eventID,
	)
	if err != nil {
		return nil, err
	}
	var timelines []Timeline
	for timelineRows.Next() {
		var t Timeline
		var eventID int
		if err := timelineRows.Scan(&t.ID, &t.Start, &t.End, &t.Place, &eventID); err != nil {
			return nil, err
		}
		timelines = append(timelines, t)
	}

	return timelines, nil
}
