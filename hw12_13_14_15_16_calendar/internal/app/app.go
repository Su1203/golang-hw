package app

import (
	"context"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warn(msg string)
	Debug(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEventByID(ctx context.Context, id string) (*storage.Event, error)
	ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	ListEventsForWeek(ctx context.Context, startDate time.Time) ([]storage.Event, error)
	ListEventsForMonth(ctx context.Context, startDate time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	a.logger.Debug("Creating event: " + event.ID)
	if err := a.storage.CreateEvent(ctx, event); err != nil {
		a.logger.Error("Failed to create event: " + err.Error())
		return err
	}
	a.logger.Info("Event created successfully: " + event.ID)
	return nil
}

func (a *App) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	a.logger.Debug("Updating event: " + id)
	if err := a.storage.UpdateEvent(ctx, id, event); err != nil {
		a.logger.Error("Failed to update event: " + err.Error())
		return err
	}
	a.logger.Info("Event updated successfully: " + id)
	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	a.logger.Debug("Deleting event: " + id)
	if err := a.storage.DeleteEvent(ctx, id); err != nil {
		a.logger.Error("Failed to delete event: " + err.Error())
		return err
	}
	a.logger.Info("Event deleted successfully: " + id)
	return nil
}

func (a *App) GetEventByID(ctx context.Context, id string) (*storage.Event, error) {
	a.logger.Debug("Getting event: " + id)
	event, err := a.storage.GetEventByID(ctx, id)
	if err != nil {
		a.logger.Error("Failed to get event: " + err.Error())
		return nil, err
	}
	return event, nil
}

func (a *App) ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	a.logger.Debug("Listing events for day: " + date.Format("2006-01-02"))
	events, err := a.storage.ListEventsForDay(ctx, date)
	if err != nil {
		a.logger.Error("Failed to list events for day: " + err.Error())
		return nil, err
	}
	return events, nil
}

func (a *App) ListEventsForWeek(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	a.logger.Debug("Listing events for week starting: " + startDate.Format("2006-01-02"))
	events, err := a.storage.ListEventsForWeek(ctx, startDate)
	if err != nil {
		a.logger.Error("Failed to list events for week: " + err.Error())
		return nil, err
	}
	return events, nil
}

func (a *App) ListEventsForMonth(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	a.logger.Debug("Listing events for month starting: " + startDate.Format("2006-01-02"))
	events, err := a.storage.ListEventsForMonth(ctx, startDate)
	if err != nil {
		a.logger.Error("Failed to list events for month: " + err.Error())
		return nil, err
	}
	return events, nil
}
