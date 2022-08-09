package app

import (
	"context"

	"github.com/Galiks/OTUS_2022/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage Storage
}

type Storage interface {
	CreateEvent(ctx context.Context, event *storage.Event) error
	UpdateEvent(ctx context.Context, event *storage.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEventsByUserID(ctx context.Context, id int64) ([]*storage.Event, error)
	GetEventByID(ctx context.Context, id int64) (*storage.Event, error)
}

func New(storage Storage) *App {
	return &App{
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.storage.CreateEvent(ctx, &storage.Event{Title: title})
}

// TODO
