package handlers

import (
	"net/http"
	"vhgomes-eventos/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

func ResponseError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, errors.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
	case errors.Is(err, errors.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	case errors.Is(err, errors.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	case errors.Is(err, errors.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": "resource already exists"})
	case errors.Is(err, errors.ErrInvalidData):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
