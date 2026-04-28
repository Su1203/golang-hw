package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrDateBusy      = errors.New("date is busy")
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidEvent  = errors.New("invalid event")
)

type Storage interface {
	CreateEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, id string, event Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEventByID(ctx context.Context, id string) (*Event, error)
	ListEventsForDay(ctx context.Context, date time.Time) ([]Event, error)
	ListEventsForWeek(ctx context.Context, startDate time.Time) ([]Event, error)
	ListEventsForMonth(ctx context.Context, startDate time.Time) ([]Event, error)
}
