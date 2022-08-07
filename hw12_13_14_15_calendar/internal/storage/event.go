package storage

import (
	"time"
)

type Event struct {
	ID          int64
	Title       string
	StartEvent  time.Time
	EndEvent    time.Time
	Description string
	UserID      int64
	EventTime   time.Time
}
