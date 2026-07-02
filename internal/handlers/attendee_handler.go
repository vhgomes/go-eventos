package handlers

import (
	"net/http"
	"strconv"
	"vhgomes-eventos/internal/pkg/logger"
	"vhgomes-eventos/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AttendeeHandler struct {
	attendeeService service.AttendeeService
}

func NewAttendeeHandler(attendeeService service.AttendeeService) *AttendeeHandler {
	return &AttendeeHandler{attendeeService: attendeeService}
}

// AddAttendee godoc
// @Summary Add attendee to event
// @Tags attendees
// @Security BearerAuth
// @Produce json
// @Param id path int true "Event ID"
// @Param userId path int true "User ID"
// @Success 201 {object} domain.Attendee
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/events/{id}/attendees/{userId} [post]
func (h *AttendeeHandler) AddAttendee(c *gin.Context) {
	currentUser, err := GetUserFromContext(c)
	if err != nil {
		ResponseError(c, err)
		return
	}

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.attendeeService.AddAttendee(c.Request.Context(), eventID, userID, currentUser.Id); err != nil {
		logger.Error("failed to add attendee", err, zap.Int("event_id", eventID), zap.Int("user_id", userID))
		ResponseError(c, err)
		return
	}

	logger.Info("attendee added successfully", zap.Int("event_id", eventID), zap.Int("user_id", userID))
	c.JSON(http.StatusCreated, gin.H{"message": "attendee added successfully"})
}

// RemoveAttendee godoc
// @Summary Remove attendee from event
// @Tags attendees
// @Security BearerAuth
// @Produce json
// @Param id path int true "Event ID"
// @Param userId path int true "User ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/events/{id}/attendees/{userId} [delete]
func (h *AttendeeHandler) RemoveAttendee(c *gin.Context) {
	currentUser, err := GetUserFromContext(c)
	if err != nil {
		ResponseError(c, err)
		return
	}

	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	userID, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.attendeeService.RemoveAttendee(c.Request.Context(), eventID, userID, currentUser.Id); err != nil {
		logger.Error("failed to remove attendee", err, zap.Int("event_id", eventID), zap.Int("user_id", userID))
		ResponseError(c, err)
		return
	}

	logger.Info("attendee removed successfully", zap.Int("event_id", eventID), zap.Int("user_id", userID))
	c.JSON(http.StatusNoContent, nil)
}

// GetAttendees godoc
// @Summary Get attendees of an event
// @Tags attendees
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {array} domain.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/events/{id}/attendees [get]
func (h *AttendeeHandler) GetAttendees(c *gin.Context) {
	eventID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	logger.Info("fetching attendees for event", zap.Int("event_id", eventID))
	users, err := h.attendeeService.GetAttendeesByEvent(c.Request.Context(), eventID)
	if err != nil {
		logger.Error("failed to fetch attendees", err, zap.Int("event_id", eventID))
		ResponseError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetEventsByAttendee godoc
// @Summary Get events by attendee
// @Tags attendees
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {array} domain.Event
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/attendees/{id}/events [get]
func (h *AttendeeHandler) GetEventsByAttendee(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	logger.Info("fetching events by attendee", zap.Int("user_id", userID))
	events, err := h.attendeeService.GetEventsByAttendee(c.Request.Context(), userID)
	if err != nil {
		logger.Error("failed to fetch events by attendee", err, zap.Int("user_id", userID))
		ResponseError(c, err)
		return
	}
	c.JSON(http.StatusOK, events)
}
