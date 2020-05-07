package files

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4"
)

func (app *App) createImage(ctx context.Context, r ImageUploadRequest) (ImageRetrieve, error) {
	baseErr := "createImage fails: %v"

	q := `
		INSERT INTO images (title, author_id, event_id, hash, created, updated)
			VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`

	row := app.assets.Db.QueryRow(
		ctx,
		q,
		r.Title,
		r.AuthorID,
		r.EventID,
		r.hash,
	)

	var image ImageRetrieve
	image.EventID = r.EventID
	image.AuthorID = r.AuthorID
	image.Title = r.Title
	image.Hash = r.hash

	if err := row.Scan(&image.ID); err != nil {
		return image, fmt.Errorf(baseErr, err)
	}

	return image, nil
}

type listFilter struct {
	eventID  int
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

func (app *App) listImage(ctx context.Context, f listFilter) ([]ImageListResponse, error) {
	baseErr := "listImage fails: %v"

	var err error
	var rows pgx.Rows
	if f.authorID != 0 {
		rows, err = app.assets.Db.Query(
			ctx, "SELECT title, author_id, event_id, hash, created FROM images WHERE author_id = $1", f.authorID)
	} else if f.eventID != 0 {
		rows, err = app.assets.Db.Query(
			ctx, "SELECT title, author_id, event_id, hash, created FROM images WHERE event_id = $1", f.eventID)
	} else {
		return nil, fmt.Errorf(baseErr, "filter is empty")
	}

	if err != nil {
		return nil, fmt.Errorf(baseErr, err)
	}
	defer rows.Close()

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
