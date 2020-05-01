package event

import (
	"strings"
	"testing"
)

func TestEvent_Validate_Positive(t *testing.T) {
	event := Event{
		Name:      "Event Positive",
		OwnerID:   1,
		Timelines: []Timeline{newValidTimeline()},
	}
	got := validateForCreate(event)

	if got != nil {
		t.Errorf("expected nil, got: %v", got)
	}
}

func TestEvent_Validate_Negative_EmptyName(t *testing.T) {
	event := Event{
		Name:      "",
		OwnerID:   1,
		Timelines: []Timeline{newValidTimeline()},
	}
	got := validateForCreate(event)

	if got == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(got.Error(), "`Title`") {
		t.Errorf("error must contains reason, got: %v", got)
	}
}

func TestEvent_Validate_Negative_EmptyTimelines(t *testing.T) {
	event := Event{
		Name:      "Event With Empty Timelines",
		OwnerID:   1,
		Timelines: []Timeline{},
	}
	got := validateForCreate(event)

	if got == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(got.Error(), "`Timelines`") {
		t.Errorf("error must contains reason, got: %v", got)
	}
}
