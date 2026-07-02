package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"vhgomes-eventos/internal/domain"
	erro "vhgomes-eventos/internal/pkg/errors"
	"vhgomes-eventos/internal/repository"
)

var _ repository.EventRepository = (*eventRepo)(nil)

type eventRepo struct {
	db *sql.DB
}

func NewEventRepo(db *sql.DB) repository.EventRepository {
	return &eventRepo{db: db}
}

func (r *eventRepo) Create(ctx context.Context, event *domain.Event) error {
	query := `INSERT INTO events (owner_id, name, description, date, location) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, event.OwnerId, event.Name, event.Description, event.Date, event.Location).Scan(&event.Id)
	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}
	return nil
}

func (r *eventRepo) GetByID(ctx context.Context, id int) (*domain.Event, error) {
	query := `SELECT id, owner_id, name, description, date, location FROM events WHERE id = $1`
	var event domain.Event
	err := r.db.QueryRowContext(ctx, query, id).Scan(&event.Id, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, erro.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get event by id: %w", err)
	}
	return &event, nil
}

func (r *eventRepo) GetAll(ctx context.Context) ([]*domain.Event, error) {
	query := `SELECT id, owner_id, name, description, date, location FROM events`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all events: %w", err)
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(&e.Id, &e.OwnerId, &e.Name, &e.Description, &e.Date, &e.Location); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return events, nil
}

func (r *eventRepo) Update(ctx context.Context, event *domain.Event) error {
	query := `UPDATE events SET name = $1, description = $2, date = $3, location = $4 WHERE id = $5`
	result, err := r.db.ExecContext(ctx, query, event.Name, event.Description, event.Date, event.Location, event.Id)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return erro.ErrNotFound
	}
	return nil
}

func (r *eventRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM events WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return erro.ErrNotFound
	}
	return nil
}

func (r *eventRepo) GetByAttendeeID(ctx context.Context, userID int) ([]*domain.Event, error) {
	query := `
		SELECT e.id, e.owner_id, e.name, e.description, e.date, e.location
		FROM events e
		JOIN attendees a ON e.id = a.event_id
		WHERE a.user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by attendee: %w", err)
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		var e domain.Event
		if err := rows.Scan(&e.Id, &e.OwnerId, &e.Name, &e.Description, &e.Date, &e.Location); err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, &e)
	}
	return events, nil
}

func (r *eventRepo) ExistsByNameAndDate(ctx context.Context, name, date string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM events WHERE name = $1 AND date = $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, name, date).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check event existence: %w", err)
	}
	return exists, nil
}
