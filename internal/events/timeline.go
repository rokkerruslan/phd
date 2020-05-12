package events

import (
	"fmt"
	"strings"
	"time"
)

const (
	TimelineStartOffset = time.Hour + time.Minute
)

type Timeline struct {
	ID    int
	Start time.Time
	End   time.Time
	Place string // Oblast'/Kray
}

func (t *Timeline) Validate() error {
	var errors []string

	if t.Start.IsZero() {
		errors = append(errors, "`Start` is zero")
	}
	if t.End.IsZero() {
		errors = append(errors, "`End` is zero")
	}
	if t.Start.Add(TimelineStartOffset).After(t.End) {
		errors = append(errors, "`Start` after end")
	}
	if t.Place == "" {
		errors = append(errors, "`Place` is empty")
	}

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf("timeline.validateForCreate fails %v", strings.Join(errors, ", "))
}
