package storage

import (
	"time"
)

type Event struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	StartEvent  time.Time `db:"start_event"`
	EndEvent    time.Time `db:"end_event"`
	Description string    `db:"description"`
	UserID      int64     `db:"user_id"`
	EventTime   time.Time `db:"event_time"`
}
