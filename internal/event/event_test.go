package event

import (
	"strings"
	"testing"
	"time"

	"photo/internal/geo"
)

func TestEvent_Validate_Positive(t *testing.T) {
	event := Event{
		Name:      "Event Positive",
		Timelines: []Timeline{{Start: time.Now(), End: time.Now().Add(2 * time.Hour)}},
		Point:     geo.Point{Lt: 1.1, Ln: 2.2},
	}
	got := event.Validate()

	if got != nil {
		t.Errorf("expected nil, got: %v", got)
	}
}

func TestEvent_Validate_Negative_EmptyName(t *testing.T) {
	event := Event{
		Name:      "",
		Timelines: []Timeline{{Start: time.Now(), End: time.Now().Add(2 * time.Hour)}},
		Point:     geo.Point{Lt: 1.1, Ln: 2.2},
	}
	got := event.Validate()

	if got == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(got.Error(), "`Name`") {
		t.Errorf("error must contains reason, got: %v", got)
	}
}

func TestEvent_Validate_Negative_EmptyTimelines(t *testing.T) {
	event := Event{
		Name:      "Event With Empty Timelines",
		Timelines: []Timeline{},
		Point:     geo.Point{Lt: 1.1, Ln: 2.2},
	}
	got := event.Validate()

	if got == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(got.Error(), "`Timelines`") {
		t.Errorf("error must contains reason, got: %v", got)
	}
}
