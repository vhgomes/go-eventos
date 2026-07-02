package handlers

import (
	"vhgomes-eventos/internal/domain"
	erro "vhgomes-eventos/internal/pkg/errors"

	"github.com/gin-gonic/gin"
)

func GetUserFromContext(c *gin.Context) (*domain.User, error) {
	val, exists := c.Get("user")
	if !exists {
		return nil, erro.ErrUnauthorized
	}
	user, ok := val.(*domain.User)
	if !ok {
		return nil, erro.ErrUnauthorized
	}
	return user, nil
}
