package event

import (
	"time"
)

type Event struct {
	ID          int
	Name        string
	Description string
	OwnerID     int
	Created     time.Time
	Updated     time.Time
	Timelines   []Timeline
}
