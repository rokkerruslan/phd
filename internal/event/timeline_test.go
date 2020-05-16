package event

import (
	"strings"
	"testing"
	"time"
)

func TestTimeline_Validate_Positive(t *testing.T) {
	src := newValidTimeline()
	got := src.Validate()

	if got != nil {
		t.Errorf("expected nil, got: %v", got)
	}
}

func TestTimeline_Validate_NegativeWrongOffset(t *testing.T) {
	src := Timeline{
		Start: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
	}
	got := src.Validate()

	if got == nil {
		t.Error("expected error, got: nil")
	}
}

func TestTimeline_Validate_EndBeforeStart(t *testing.T) {
	src := Timeline{
		Start: time.Date(2020, 1, 1, 1, 0, 0, 0, time.UTC),
		End:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	got := src.Validate()

	if got == nil {
		t.Error("expected error, got: nil")
	}
}

func TestTimeline_Validate_StartZero(t *testing.T) {
	src := Timeline{
		End: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	got := src.Validate()

	if got == nil {
		t.Error("expected error, got: nil")
	}
}

func TestTimeline_Validate_EndZero(t *testing.T) {
	src := Timeline{
		Start: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	got := src.Validate()

	if got == nil {
		t.Error("expected error, got: nil")
	}
}

func TestTimeline_Validate_EmptyPlace(t *testing.T) {
	src := newValidTimeline()
	src.Place = ""
	got := src.Validate()

	if got == nil || !strings.Contains(got.Error(), "`Place`") {
		t.Errorf("expected error with `Place`, got: %v", got)
	}
}

func newValidTimeline() Timeline {
	return Timeline{
		Start: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		End:   time.Date(2020, 1, 1, 2, 0, 0, 0, time.UTC),
		Place: "San Andreas",
	}
}