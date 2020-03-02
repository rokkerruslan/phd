package event

import (
	"photo/internal/geo"
)

var events = []Event{
	{
		ID:        1,
		Name:      "Event 1",
		Timelines: nil,
		Point:     geo.Point{Lt: 1, Ln: 2},
	},
	{
		ID:        2,
		Name:      "Event 2",
		Timelines: nil,
		Point:     geo.Point{Lt: 13.1, Ln: 12.1},
	},
	{
		ID:        3,
		Name:      "Event 3",
		Timelines: nil,
		Point:     geo.Point{Lt: 81.1, Ln: 21.2},
	},
	{
		ID:        4,
		Name:      "Event 4",
		Timelines: nil,
		Point:     geo.Point{Lt: 4.1, Ln: 51.1},
	},
}
