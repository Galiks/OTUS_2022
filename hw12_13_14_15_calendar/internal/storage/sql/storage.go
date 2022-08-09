package sqlstorage

import (
	"context"
	"log"

	"github.com/Galiks/OTUS_2022/hw12_13_14_15_calendar/internal/logger"
	"github.com/Galiks/OTUS_2022/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib" //nolint
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(dataSourceName string) *Storage {
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		log.Fatal(err)
		logger.Fatal(err)
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
	query := `
		INSERT INTO events
		(
			title,
			start_event,
			end_event,
			description,
			user_id,
			event_time,
		)
		VALUES
		(
			:title,
			:start_event,
			:end_event,
			:description,
			:user_id,
			:event_time,
		) 
		RETURNING id
	`
	row, err := s.db.NamedQueryContext(ctx, query, event)
	if err != nil {
		return err
	}
	defer row.Close()

	if row.Next() {
		if err := row.Scan(&event.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	query := `
		UPDATE events
		SET
			title=:title,
			start_event=:start_event,
			end_event=:end_event,
			description=:description,
			user_id=:user_id,
			event_time=:event_time
		WHERE id=:id
	`
	row, err := s.db.NamedQueryContext(ctx, query, event)
	if err != nil {
		return err
	}
	if row.Next() {
		if err := row.Scan(&event.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	query := `
		DELETE FROM events
		WHERE id=$1
	`
	if _, err := s.db.NamedQueryContext(ctx, query, id); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetEventsByUserID(ctx context.Context, id int64) ([]*storage.Event, error) {
	var result []*storage.Event = make([]*storage.Event, 0)
	query := `
		SELECT *
		FROM events
		WHERE user_id=$1
	`

	rows, err := s.db.NamedQueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var event storage.Event
		if err := rows.Scan(&event); err != nil {
			return nil, err
		}
		result = append(result, &event)
	}

	return result, nil
}

func (s *Storage) GetEventByID(ctx context.Context, id int64) (*storage.Event, error) {
	query := `
		SELECT *
		FROM events
		WHERE id=$1
	`

	row, err := s.db.NamedQueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	if row.Next() {
		var event storage.Event
		if err := row.Scan(&event); err != nil {
			return nil, err
		}
		return &event, nil
	}
	return nil, storage.ErrUnknownEvent
}
