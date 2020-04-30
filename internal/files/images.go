package files

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ImageUploadRequest struct {
	EventID  int
	AuthorID int
	Name     string
	Data     string

	hash string
}

func (r *ImageUploadRequest) Validate() error {
	var e []string

	if r.EventID == 0 {
		e = append(e, "`EventID` is missing")
	}
	if r.AuthorID == 0 {
		e = append(e, "`AuthorID` is missing")
	}
	if r.Name == "" {
		e = append(e, "`Name` is missing")
	}
	if r.Data == "" {
		e = append(e, "`Data` is missing")
	}

	if len(e) == 0 {
		return nil
	}

	return fmt.Errorf("ImageUploadRequest.Validate fails: %v", strings.Join(e, ", "))
}

func (r *ImageUploadRequest) Store() error {
	baseErr := "ImageUploadRequest.Store() fails: %v"

	buf, err := base64.StdEncoding.DecodeString(r.Data)
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}

	hash := sha256.New()
	hash.Write([]byte(r.Name))

	// TODO: more clever solution (use GLOBAL_SALT)
	timeHash, err := time.Now().MarshalBinary()
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}
	r.hash = fmt.Sprintf("%x", hash.Sum(timeHash))

	// TODO: use hash sum for name
	f, err := os.OpenFile(filepath.Join("./images", r.hash), os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return fmt.Errorf(baseErr, err)
	}
	defer f.Close()

	// TODO: fill rewrite file
	if _, err := f.Write(buf); err != nil {
		return fmt.Errorf(baseErr, err)
	}

	return nil
}

func (app *App) createImage(ctx context.Context, r ImageUploadRequest) error {
	baseErr := "createImage fails: %v"

	var id int
	row := app.assets.Db.QueryRow(
		ctx,
		"INSERT INTO images (title, author_id, event_id, hash, created, updated) VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id",
		r.Name,
		r.AuthorID,
		r.EventID,
		r.hash,
	)
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf(baseErr, err)
	}

	return nil
}

type Image struct {
	ID       int
	Path     string
	Name     string
	EventID  int
	AuthorID int
	Hash     string
}
