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

var _ repository.AttendeeRepository = (*attendeeRepo)(nil)

type attendeeRepo struct {
	db *sql.DB
}

func NewAttendeeRepo(db *sql.DB) repository.AttendeeRepository {
	return &attendeeRepo{db: db}
}

func (r *attendeeRepo) Create(ctx context.Context, attendee *domain.Attendee) error {
	query := `INSERT INTO attendees (event_id, user_id) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, attendee.EventId, attendee.UserId).Scan(&attendee.Id)
	if err != nil {
		return fmt.Errorf("failed to create attendee: %w", err)
	}
	return nil
}

func (r *attendeeRepo) GetByEventAndUser(ctx context.Context, eventID, userID int) (*domain.Attendee, error) {
	query := `SELECT id, event_id, user_id FROM attendees WHERE event_id = $1 AND user_id = $2`
	var a domain.Attendee
	err := r.db.QueryRowContext(ctx, query, eventID, userID).Scan(&a.Id, &a.EventId, &a.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, erro.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get attendee: %w", err)
	}
	return &a, nil
}

func (r *attendeeRepo) GetUsersByEvent(ctx context.Context, eventID int) ([]*domain.User, error) {
	query := `
		SELECT u.id, u.name, u.email
		FROM users u
		JOIN attendees a ON u.id = a.user_id
		WHERE a.event_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by event: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.Id, &u.Name, &u.Email); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &u)
	}
	return users, nil
}

func (r *attendeeRepo) Delete(ctx context.Context, eventID, userID int) error {
	query := `DELETE FROM attendees WHERE event_id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, eventID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete attendee: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return erro.ErrNotFound
	}
	return nil
}
