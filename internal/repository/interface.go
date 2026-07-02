package repository

import (
	"context"
	"vhgomes-eventos/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	GetByID(ctx context.Context, id int) (*domain.Event, error)
	GetAll(ctx context.Context) ([]*domain.Event, error)
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id int) error
	GetByAttendeeID(ctx context.Context, userID int) ([]*domain.Event, error)
	ExistsByNameAndDate(ctx context.Context, name, date string) (bool, error)
}

type AttendeeRepository interface {
	Create(ctx context.Context, attendee *domain.Attendee) error
	GetByEventAndUser(ctx context.Context, eventID, userID int) (*domain.Attendee, error)
	GetUsersByEvent(ctx context.Context, eventID int) ([]*domain.User, error)
	Delete(ctx context.Context, eventID, userID int) error
}
