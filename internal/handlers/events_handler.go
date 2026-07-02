package handlers

import (
	"net/http"
	"strconv"
	"vhgomes-eventos/internal/service"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	eventService service.EventService
}

func NewEventHandler(eventService service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

// Create godoc
// @Summary Create a new event
// @Tags events
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param event body service.CreateEventRequest true "Event data"
// @Success 201 {object} domain.Event
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/events [post]
func (h *EventHandler) Create(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		ResponseError(c, err)
		return
	}

	var req service.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := h.eventService.Create(c.Request.Context(), user.Id, req)
	if err != nil {
		ResponseError(c, err)
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetAll godoc
// @Summary Get all events
// @Tags events
// @Produce json
// @Success 200 {array} domain.Event
// @Router /api/v1/events [get]
func (h *EventHandler) GetAll(c *gin.Context) {
	events, err := h.eventService.GetAll(c.Request.Context())
	if err != nil {
		ResponseError(c, err)
		return
	}
	c.JSON(http.StatusOK, events)
}

// GetByID godoc
// @Summary Get event by ID
// @Tags events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} domain.Event
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/events/{id} [get]
func (h *EventHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	event, err := h.eventService.GetByID(c.Request.Context(), id)
	if err != nil {
		ResponseError(c, err)
		return
	}
	c.JSON(http.StatusOK, event)
}

// Update godoc
// @Summary Update an event
// @Tags events
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param event body service.UpdateEventRequest true "Event data"
// @Success 200 {object} domain.Event
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/events/{id} [put]
func (h *EventHandler) Update(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		ResponseError(c, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	var req service.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.eventService.Update(c.Request.Context(), id, user.Id, req); err != nil {
		ResponseError(c, err)
		return
	}

	// Busca o evento atualizado para retornar
	event, _ := h.eventService.GetByID(c.Request.Context(), id)
	c.JSON(http.StatusOK, event)
}

// Delete godoc
// @Summary Delete an event
// @Tags events
// @Security BearerAuth
// @Produce json
// @Param id path int true "Event ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/events/{id} [delete]
func (h *EventHandler) Delete(c *gin.Context) {
	user, err := GetUserFromContext(c)
	if err != nil {
		ResponseError(c, err)
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := h.eventService.Delete(c.Request.Context(), id, user.Id); err != nil {
		ResponseError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
