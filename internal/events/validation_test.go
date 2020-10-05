package events

import (
	"testing"
	"time"
)

func TestTimelinesValidation(t *testing.T) {
	t.Run("one", func(t *testing.T) {
		base := time.Now()
		timelines := []Timeline{
			{Start: base, End: base.Add(1*time.Hour)},
		}

		if !validateTimelinesCrossing(timelines) {
			t.Error("timelines crossing")
		}
	})

	t.Run("multiple positive", func(t *testing.T) {
		base := time.Now()
		timelines := []Timeline{
			{Start: base.Add(time.Minute), End: base.Add(time.Hour - time.Minute)},
			{Start: base.Add(time.Hour + time.Minute), End: base.Add(2*time.Hour)},
		}

		if !validateTimelinesCrossing(timelines) {
			t.Error("timelines crossing")
		}
	})

	t.Run("multiple positive (unsorted)", func(t *testing.T) {
		base := time.Now()
		timelines := []Timeline{
			{Start: base.Add(time.Hour + time.Minute), End: base.Add(2*time.Hour)},
			{Start: base.Add(time.Minute), End: base.Add(time.Hour - time.Minute)},
		}

		if !validateTimelinesCrossing(timelines) {
			t.Error("timelines crossing")
		}
	})

	t.Run("multiple negative", func(t *testing.T) {
		base := time.Now()
		timelines := []Timeline{
			{Start: base.Add(time.Minute), End: base.Add(time.Hour + time.Minute)},
			{Start: base.Add(time.Hour - time.Minute), End: base.Add(2*time.Hour)},
		}

		if validateTimelinesCrossing(timelines) {
			t.Error("timelines crossing")
		}
	})

	t.Run("multiple negative", func(t *testing.T) {
		base := time.Now()
		timelines := []Timeline{
			{Start: base.Add(time.Minute), End: base.Add(time.Hour - time.Minute)},
			{Start: base.Add(time.Hour + time.Minute), End: base.Add(2*time.Hour)},
			{Start: base.Add(2*time.Hour-time.Minute), End: base.Add(3*time.Hour)},
		}

		if validateTimelinesCrossing(timelines) {
			t.Error("timelines crossing")
		}
	})
}
