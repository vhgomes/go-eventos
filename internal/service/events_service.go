package service

import (
	"context"
	"errors"
	"fmt"
	"vhgomes-eventos/internal/domain"
	erro "vhgomes-eventos/internal/pkg/errors"
	"vhgomes-eventos/internal/repository"
)

type EventService interface {
	Create(ctx context.Context, ownerID int, req CreateEventRequest) (*domain.Event, error)
	GetByID(ctx context.Context, id int) (*domain.Event, error)
	GetAll(ctx context.Context) ([]*domain.Event, error)
	Update(ctx context.Context, id, ownerID int, req UpdateEventRequest) error
	Delete(ctx context.Context, id, ownerID int) error
	GetByAttendeeID(ctx context.Context, userID int) ([]*domain.Event, error)
}

type CreateEventRequest struct {
	Name        string
	Description string
	Date        string
	Location    string
}

type UpdateEventRequest struct {
	Name        string
	Description string
	Date        string
	Location    string
}

type eventService struct {
	repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) EventService {
	return &eventService{repo: repo}
}

func (s *eventService) Create(ctx context.Context, ownerID int, req CreateEventRequest) (*domain.Event, error) {
	exists, err := s.repo.ExistsByNameAndDate(ctx, req.Name, req.Date)
	if err != nil {
		return nil, fmt.Errorf("failed to check event existence: %w", err)
	}
	if exists {
		return nil, erro.ErrConflict
	}

	event := &domain.Event{
		OwnerId:     ownerID,
		Name:        req.Name,
		Description: req.Description,
		Date:        req.Date,
		Location:    req.Location,
	}

	if err := s.repo.Create(ctx, event); err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return event, nil
}

func (s *eventService) GetByID(ctx context.Context, id int) (*domain.Event, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, erro.ErrNotFound) {
			return nil, erro.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get event: %w", err)
	}
	return event, nil
}

func (s *eventService) GetAll(ctx context.Context) ([]*domain.Event, error) {
	events, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all events: %w", err)
	}
	return events, nil
}

func (s *eventService) Update(ctx context.Context, id, ownerID int, req UpdateEventRequest) error {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if event.OwnerId != ownerID {
		return erro.ErrUnauthorized
	}

	event.Name = req.Name
	event.Description = req.Description
	event.Date = req.Date
	event.Location = req.Location

	if err := s.repo.Update(ctx, event); err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	return nil
}

func (s *eventService) Delete(ctx context.Context, id, ownerID int) error {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if event.OwnerId != ownerID {
		return erro.ErrUnauthorized
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}

func (s *eventService) GetByAttendeeID(ctx context.Context, userID int) ([]*domain.Event, error) {
	events, err := s.repo.GetByAttendeeID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by attendee: %w", err)
	}
	return events, nil
}
