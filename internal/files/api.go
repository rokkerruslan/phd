package files

import (
	"encoding/json"
	"fmt"
	"net/http"

	"ph/internal/api"
)

func (app *App) listHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "files.listHandler fails: %v"

	f, err := newListFilter(r.URL.Query())
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	list, err := app.listImage(r.Context(), f)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	api.Response(w, list)
}

func (app *App) uploadHandler(w http.ResponseWriter, r *http.Request) {
	baseErr := "uploadHandler fails: %v"

	var data ImageUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	if err := data.Validate(); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	accountID, err := app.assets.Tokens.RetrieveAccountIDFromRequest(r.Context(), r)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusBadRequest)
		return
	}

	if data.AuthorID != accountID {
		api.Error(w, fmt.Errorf(baseErr, "you can't create image from this account"), http.StatusBadRequest)
		return
	}

	// TODO: check hidden/deleted status
	var isEventPublic bool
	row := app.assets.Db.QueryRow(r.Context(), "SELECT is_public FROM events WHERE id = $1", data.EventID)
	if err := row.Scan(&isEventPublic); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	if !isEventPublic {
		// TODO: Check creating availability
	}

	if err := data.Store(); err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	image, err := app.createImage(r.Context(), data)
	if err != nil {
		api.Error(w, fmt.Errorf(baseErr, err), http.StatusInternalServerError)
		return
	}

	api.Response(w, image)
}
