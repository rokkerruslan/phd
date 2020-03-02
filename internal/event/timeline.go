package event

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
}

func (t *Timeline) Validate() error {
	var errors []string

	if t.Start.IsZero() {
		errors = append(errors, "start is zero")
	}
	if t.End.IsZero() {
		errors = append(errors, "end is zero")
	}
	if t.Start.Add(TimelineStartOffset).After(t.End) {
		errors = append(errors, "start after end")
	}

	if len(errors) == 0 {
		return nil
	}

	return fmt.Errorf("timeline.validateForCreate fails %v", strings.Join(errors, ", "))
}
