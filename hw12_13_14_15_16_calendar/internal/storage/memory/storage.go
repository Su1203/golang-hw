package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/Su1203/golang-hw/hw12_13_14_15_16_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) CreateEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event.ID == "" || event.Title == "" {
		return storage.ErrInvalidEvent
	}

	if _, exists := s.events[event.ID]; exists {
		return storage.ErrDateBusy
	}

	for _, e := range s.events {
		if e.UserID == event.UserID && s.isOverlapping(e, event) {
			return storage.ErrDateBusy
		}
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, id string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return storage.ErrEventNotFound
	}

	if event.Title == "" {
		return storage.ErrInvalidEvent
	}

	for eventID, e := range s.events {
		if eventID != id && e.UserID == event.UserID && s.isOverlapping(e, event) {
			return storage.ErrDateBusy
		}
	}

	event.ID = id
	s.events[id] = event
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[id]; !exists {
		return storage.ErrEventNotFound
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) GetEventByID(ctx context.Context, id string) (*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.events[id]
	if !exists {
		return nil, storage.ErrEventNotFound
	}

	return &event, nil
}

func (s *Storage) ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return s.filterEvents(startOfDay, endOfDay), nil
}

func (s *Storage) ListEventsForWeek(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startOfWeek := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endOfWeek := startOfWeek.Add(7 * 24 * time.Hour)

	return s.filterEvents(startOfWeek, endOfWeek), nil
}

func (s *Storage) ListEventsForMonth(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	startOfMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return s.filterEvents(startOfMonth, endOfMonth), nil
}

func (s *Storage) isOverlapping(e1, e2 storage.Event) bool {
	return e1.StartTime.Before(e2.EndTime) && e2.StartTime.Before(e1.EndTime)
}

func (s *Storage) filterEvents(start, end time.Time) []storage.Event {
	var result []storage.Event
	for _, event := range s.events {
		if event.StartTime.Before(end) && event.EndTime.After(start) {
			result = append(result, event)
		}
	}
	return result
}
