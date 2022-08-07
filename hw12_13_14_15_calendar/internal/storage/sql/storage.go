package sqlstorage

import (
	"context"

	"github.com/Galiks/OTUS_2022/hw12_13_14_15_calendar/internal/logger"
	"github.com/Galiks/OTUS_2022/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(dataSourceName string) *Storage {
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		logger.Fatal(err.Error())
	}
	return &Storage{
		db: db,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) error {
	panic("not implemented") // TODO: Implement
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	panic("not implemented") // TODO: Implement
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	panic("not implemented") // TODO: Implement
}

func (s *Storage) GetEventsByUserID(ctx context.Context, id int64) ([]*storage.Event, error) {
	panic("not implemented") // TODO: Implement
}

func (s *Storage) GetEventByID(ctx context.Context, id int64) (*storage.Event, error) {
	panic("not implemented") // TODO: Implement
}
