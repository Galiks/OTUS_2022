package memorystorage

import (
	"context"
	"sync"

	"github.com/Galiks/OTUS_2022/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	storage map[int64]*storage.Event
	mu      sync.RWMutex
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(ctx context.Context, event *storage.Event) error {
	s.mu.Lock()
	event.ID = int64(len(s.storage)) + 1
	s.storage[event.ID] = event
	s.mu.Unlock()
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event *storage.Event) error {
	s.mu.Lock()
	event.ID = int64(len(s.storage)) + 1
	s.storage[event.ID] = event
	s.mu.Unlock()
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id int64) error {
	s.mu.Lock()
	delete(s.storage, id)
	s.mu.Unlock()
	return nil
}

func (s *Storage) GetEventsByUserID(ctx context.Context, id int64) ([]*storage.Event, error) {
	var result = make([]*storage.Event, 0)
	s.mu.RLock()
	for _, event := range s.storage {
		if event.UserID == id {
			result = append(result, event)
		}
	}
	s.mu.RUnlock()
	return result, nil
}

func (s *Storage) GetEventByID(ctx context.Context, id int64) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.storage[id], nil
}
