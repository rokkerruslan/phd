package files

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func (app *App) createImage(ctx context.Context, r ImageUploadRequest) error {
	baseErr := "createImage fails: %v"

	var id int
	row := app.assets.Db.QueryRow(
		ctx,
		"INSERT INTO images (title, author_id, event_id, hash, created, updated) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id",
		r.Title,
		r.AuthorID,
		r.EventID,
		r.hash,
	)
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf(baseErr, err)
	}

	return nil
}

type listFilter struct {
	eventID int
	authorID int
}

func newListFilter(values url.Values) (listFilter, error) {
	baseErr := "listFilter fails: %v"

	var f listFilter
	var err error
	var errors []string

	authorID := values.Get("author_id")
	if authorID != "" {
		if f.authorID, err = strconv.Atoi(authorID); err != nil {
			errors = append(errors, fmt.Sprintf("author_id parsing fails: %v", err))
		}
	}

	eventIDRaw := values.Get("event_id")
	if eventIDRaw != "" {
		if f.eventID, err = strconv.Atoi(eventIDRaw); err != nil {
			errors = append(errors, fmt.Sprintf("event_id parsing fails: %v", err))
		}
	}

	if f.eventID == 0 && f.authorID == 0 {
		errors = append(errors, "filter is empty")
	}
	if f.eventID != 0 && f.authorID != 0 {
		errors = append(errors, "only one parameter available")
	}

	if len(errors) != 0 {
		return f, fmt.Errorf(baseErr, strings.Join(errors, ", "))
	}

	return f, nil
}

func (app *App) listImage(ctx context.Context, _ listFilter) ([]ImageListResponse, error) {
	baseErr := "listImage fails: %v"

	rows, err := app.assets.Db.Query(
		ctx,
		"SELECT title, author_id, event_id, hash, created FROM images",
	)
	if err != nil {
		return nil, fmt.Errorf(baseErr, err)
	}

	list := make([]ImageListResponse, 0)
	for rows.Next() {
		var el ImageListResponse
		if err := rows.Scan(&el.Title, &el.AuthorID, &el.EventID, &el.Hash, &el.Created); err != nil {
			return nil, fmt.Errorf(baseErr, err)
		}
		list = append(list, el)
	}

	return list, nil
}
