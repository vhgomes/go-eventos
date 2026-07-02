package service

import (
	"context"
	"errors"
	"fmt"
	"vhgomes-eventos/internal/domain"
	erro "vhgomes-eventos/internal/pkg/errors"
	"vhgomes-eventos/internal/repository"
)

type AttendeeService interface {
	AddAttendee(ctx context.Context, eventID, userID, currentUserID int) error
	RemoveAttendee(ctx context.Context, eventID, userID, currentUserID int) error
	GetAttendeesByEvent(ctx context.Context, eventID int) ([]*domain.User, error)
	GetEventsByAttendee(ctx context.Context, userID int) ([]*domain.Event, error)
}

type attendeeService struct {
	attendeeRepo repository.AttendeeRepository
	eventRepo    repository.EventRepository
	userRepo     repository.UserRepository
}

func NewAttendeeService(
	attendeeRepo repository.AttendeeRepository,
	eventRepo repository.EventRepository,
	userRepo repository.UserRepository,
) AttendeeService {
	return &attendeeService{
		attendeeRepo: attendeeRepo,
		eventRepo:    eventRepo,
		userRepo:     userRepo,
	}
}

func (s *attendeeService) AddAttendee(ctx context.Context, eventID, userID, currentUserID int) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.OwnerId != currentUserID {
		return erro.ErrUnauthorized
	}

	_, err = s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	existing, err := s.attendeeRepo.GetByEventAndUser(ctx, eventID, userID)
	if err != nil && !errors.Is(err, erro.ErrNotFound) {
		return fmt.Errorf("failed to check existing attendee: %w", err)
	}
	if existing != nil {
		return erro.ErrConflict
	}

	attendee := &domain.Attendee{
		EventId: eventID,
		UserId:  userID,
	}
	if err := s.attendeeRepo.Create(ctx, attendee); err != nil {
		return fmt.Errorf("failed to add attendee: %w", err)
	}
	return nil
}

func (s *attendeeService) RemoveAttendee(ctx context.Context, eventID, userID, currentUserID int) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.OwnerId != currentUserID {
		return erro.ErrUnauthorized
	}

	if err := s.attendeeRepo.Delete(ctx, eventID, userID); err != nil {
		if errors.Is(err, erro.ErrNotFound) {
			return erro.ErrNotFound
		}
		return fmt.Errorf("failed to remove attendee: %w", err)
	}
	return nil
}

func (s *attendeeService) GetAttendeesByEvent(ctx context.Context, eventID int) ([]*domain.User, error) {
	_, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	users, err := s.attendeeRepo.GetUsersByEvent(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attendees: %w", err)
	}
	return users, nil
}

func (s *attendeeService) GetEventsByAttendee(ctx context.Context, userID int) ([]*domain.Event, error) {
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	events, err := s.eventRepo.GetByAttendeeID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get events by attendee: %w", err)
	}
	return events, nil
}
