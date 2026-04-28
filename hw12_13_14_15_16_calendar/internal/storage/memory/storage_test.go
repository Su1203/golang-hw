package memorystorage

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/Su1203/golang-hw/hw12_13_14_15_16_calendar/internal/storage"
)

func TestCreateEvent(t *testing.T) {
	s := New()
	ctx := context.Background()

	event := storage.Event{
		ID:        "1",
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		UserID:    "user1",
	}

	err := s.CreateEvent(ctx, event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	retrieved, err := s.GetEventByID(ctx, "1")
	if err != nil {
		t.Fatalf("Failed to get event: %v", err)
	}

	if retrieved.Title != event.Title {
		t.Errorf("Expected title %s, got %s", event.Title, retrieved.Title)
	}
}

func TestCreateEventInvalidData(t *testing.T) {
	s := New()
	ctx := context.Background()

	tests := []struct {
		name  string
		event storage.Event
	}{
		{
			name: "empty ID",
			event: storage.Event{
				Title:     "Test",
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
			},
		},
		{
			name: "empty title",
			event: storage.Event{
				ID:        "1",
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.CreateEvent(ctx, tt.event)
			if !errors.Is(err, storage.ErrInvalidEvent) {
				t.Errorf("Expected ErrInvalidEvent, got %v", err)
			}
		})
	}
}

func TestCreateEventDateBusy(t *testing.T) {
	s := New()
	ctx := context.Background()

	event1 := storage.Event{
		ID:        "1",
		Title:     "Event 1",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		UserID:    "user1",
	}

	err := s.CreateEvent(ctx, event1)
	if err != nil {
		t.Fatalf("Failed to create first event: %v", err)
	}

	overlappingEvent := storage.Event{
		ID:        "2",
		Title:     "Event 2",
		StartTime: time.Now().Add(30 * time.Minute),
		EndTime:   time.Now().Add(90 * time.Minute),
		UserID:    "user1",
	}

	err = s.CreateEvent(ctx, overlappingEvent)
	if !errors.Is(err, storage.ErrDateBusy) {
		t.Errorf("Expected ErrDateBusy, got %v", err)
	}
}

func TestUpdateEvent(t *testing.T) {
	s := New()
	ctx := context.Background()

	event := storage.Event{
		ID:        "1",
		Title:     "Original Title",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		UserID:    "user1",
	}

	err := s.CreateEvent(ctx, event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	updatedEvent := storage.Event{
		Title:     "Updated Title",
		StartTime: time.Now().Add(2 * time.Hour),
		EndTime:   time.Now().Add(3 * time.Hour),
		UserID:    "user1",
	}

	err = s.UpdateEvent(ctx, "1", updatedEvent)
	if err != nil {
		t.Fatalf("Failed to update event: %v", err)
	}

	retrieved, err := s.GetEventByID(ctx, "1")
	if err != nil {
		t.Fatalf("Failed to get event: %v", err)
	}

	if retrieved.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %s", retrieved.Title)
	}
}

func TestUpdateEventNotFound(t *testing.T) {
	s := New()
	ctx := context.Background()

	event := storage.Event{
		Title:     "Test",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		UserID:    "user1",
	}

	err := s.UpdateEvent(ctx, "nonexistent", event)
	if !errors.Is(err, storage.ErrEventNotFound) {
		t.Errorf("Expected ErrEventNotFound, got %v", err)
	}
}

func TestDeleteEvent(t *testing.T) {
	s := New()
	ctx := context.Background()

	event := storage.Event{
		ID:        "1",
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		UserID:    "user1",
	}

	err := s.CreateEvent(ctx, event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	err = s.DeleteEvent(ctx, "1")
	if err != nil {
		t.Fatalf("Failed to delete event: %v", err)
	}

	_, err = s.GetEventByID(ctx, "1")
	if err != storage.ErrEventNotFound {
		t.Errorf("Expected ErrEventNotFound after deletion, got %v", err)
	}
}

func TestListEventsForDay(t *testing.T) {
	s := New()
	ctx := context.Background()

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 10, 0, 0, 0, now.Location())

	event1 := storage.Event{
		ID:        "1",
		Title:     "Today Event",
		StartTime: today,
		EndTime:   today.Add(time.Hour),
		UserID:    "user1",
	}

	event2 := storage.Event{
		ID:        "2",
		Title:     "Tomorrow Event",
		StartTime: today.Add(25 * time.Hour),
		EndTime:   today.Add(26 * time.Hour),
		UserID:    "user1",
	}

	if err := s.CreateEvent(ctx, event1); err != nil {
		t.Fatalf("Failed to create event1: %v", err)
	}
	if err := s.CreateEvent(ctx, event2); err != nil {
		t.Fatalf("Failed to create event2: %v", err)
	}

	events, err := s.ListEventsForDay(ctx, today)
	if err != nil {
		t.Fatalf("Failed to list events: %v", err)
	}

	if len(events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(events))
	}

	if len(events) > 0 && events[0].ID != "1" {
		t.Errorf("Expected event ID '1', got %s", events[0].ID)
	}
}

func TestConcurrentAccess(t *testing.T) {
	s := New()
	ctx := context.Background()

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			event := storage.Event{
				ID:        string(rune('a' + id)),
				Title:     "Concurrent Event",
				StartTime: time.Now().Add(time.Duration(id) * time.Hour),
				EndTime:   time.Now().Add(time.Duration(id)*time.Hour + 30*time.Minute),
				UserID:    "user1",
			}

			err := s.CreateEvent(ctx, event)
			if err != nil && err != storage.ErrDateBusy && err != storage.ErrInvalidEvent {
				t.Errorf("Unexpected error: %v", err)
			}
		}(i)
	}

	wg.Wait()
}

func TestConcurrentReadWrite(t *testing.T) {
	s := New()
	ctx := context.Background()

	event := storage.Event{
		ID:        "1",
		Title:     "Test Event",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		UserID:    "user1",
	}

	err := s.CreateEvent(ctx, event)
	if err != nil {
		t.Fatalf("Failed to create event: %v", err)
	}

	var wg sync.WaitGroup
	numReaders := 50
	numWriters := 10

	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = s.GetEventByID(ctx, "1")
		}()
	}

	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			updatedEvent := storage.Event{
				Title:     "Updated",
				StartTime: time.Now().Add(time.Duration(id) * time.Hour),
				EndTime:   time.Now().Add(time.Duration(id)*time.Hour + time.Hour),
				UserID:    "user1",
			}
			_ = s.UpdateEvent(ctx, "1", updatedEvent)
		}(i)
	}

	wg.Wait()
}
