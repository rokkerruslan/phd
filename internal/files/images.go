package files

import (
	"context"
)

type Image struct {
	ID       int
	Path     string
	Name     string
	EventID  int
	AuthorID int
}

func (a *app) createImage(ctx context.Context) {

}
