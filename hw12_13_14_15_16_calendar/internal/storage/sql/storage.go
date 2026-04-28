package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func New(dsn string) (*Storage, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	query := `
		INSERT INTO events (id, title, start_time, end_time, description, user_id, notify_before)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	var overlapping int
	checkQuery := `
		SELECT COUNT(1) FROM events 
		WHERE user_id = $1 
		AND start_time < $2 
		AND end_time > $3
	`
	err := s.db.GetContext(ctx, &overlapping, checkQuery, event.UserID, event.EndTime, event.StartTime)
	if err != nil {
		return fmt.Errorf("failed to check overlapping events: %w", err)
	}
	if overlapping > 0 {
		return storage.ErrDateBusy
	}

	_, err = s.db.ExecContext(ctx, query,
		event.ID,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.UserID,
		event.NotifyBefore,
	)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, id string, event storage.Event) error {
	var overlapping int
	checkQuery := `
		SELECT COUNT(1) FROM events 
		WHERE user_id = $1 
		AND id != $2
		AND start_time < $3 
		AND end_time > $4
	`
	err := s.db.GetContext(ctx, &overlapping, checkQuery, event.UserID, id, event.EndTime, event.StartTime)
	if err != nil {
		return fmt.Errorf("failed to check overlapping events: %w", err)
	}
	if overlapping > 0 {
		return storage.ErrDateBusy
	}

	query := `
		UPDATE events 
		SET title = $1, start_time = $2, end_time = $3, description = $4, user_id = $5, notify_before = $6
		WHERE id = $7
	`

	result, err := s.db.ExecContext(ctx, query,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.UserID,
		event.NotifyBefore,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return storage.ErrEventNotFound
	}

	return nil
}

func (s *Storage) GetEventByID(ctx context.Context, id string) (*storage.Event, error) {
	query := `
		SELECT id, title, start_time, end_time, description, user_id, notify_before
		FROM events
		WHERE id = $1
	`

	var event storage.Event
	err := s.db.GetContext(ctx, &event, query, id)
	if err == sql.ErrNoRows {
		return nil, storage.ErrEventNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return &event, nil
}

func (s *Storage) ListEventsForDay(ctx context.Context, date time.Time) ([]storage.Event, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	return s.listEvents(ctx, startOfDay, endOfDay)
}

func (s *Storage) ListEventsForWeek(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	startOfWeek := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endOfWeek := startOfWeek.Add(7 * 24 * time.Hour)

	return s.listEvents(ctx, startOfWeek, endOfWeek)
}

func (s *Storage) ListEventsForMonth(ctx context.Context, startDate time.Time) ([]storage.Event, error) {
	startOfMonth := time.Date(startDate.Year(), startDate.Month(), 1, 0, 0, 0, 0, startDate.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	return s.listEvents(ctx, startOfMonth, endOfMonth)
}

func (s *Storage) listEvents(ctx context.Context, start, end time.Time) ([]storage.Event, error) {
	query := `
		SELECT id, title, start_time, end_time, description, user_id, notify_before
		FROM events
		WHERE start_time < $1 AND end_time > $2
		ORDER BY start_time
	`

	var events []storage.Event
	err := s.db.SelectContext(ctx, &events, query, end, start)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	return events, nil
}
