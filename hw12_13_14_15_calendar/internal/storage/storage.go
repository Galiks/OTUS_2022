package storage

import "context"

type Storage interface {
	CreateEvent(ctx context.Context, event *Event) error
	UpdateEvent(ctx context.Context, event *Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEventsByUserID(ctx context.Context, id int64) ([]*Event, error)
	GetEventByID(ctx context.Context, id int64) (*Event, error)
}
