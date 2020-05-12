package events

import (
	"context"
	"fmt"
)

func (app *App) suggestedEvents(ctx context.Context, accountID int) ([]Event, error) {
	baseErr := "suggestedEvents fails: %v"

	rows, err := app.assets.Db.Query(
		ctx,
		"SELECT id, name, description, owner_id, created, updated, is_public, is_hidden FROM events WHERE owner_id != $1",
		accountID,
	)
	if err != nil {
		return nil, fmt.Errorf(baseErr, err)
	}

	events, err := app.constructEventList(ctx, rows)
	if err != nil {
		return nil, fmt.Errorf(baseErr, err)
	}

	return events, nil
}
