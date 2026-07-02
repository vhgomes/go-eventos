package main

import (
	"errors"
	"vhgomes-eventos/internal/database"

	"github.com/gin-gonic/gin"
)

var ErrUnauthorized = errors.New("user not authenticated")

func (app *application) GetUserFromContext(c *gin.Context) (*database.User, error) {
	contextUser, exists := c.Get("user")

	if !exists {
		return nil, ErrUnauthorized
	}

	user, ok := contextUser.(*database.User)
	if !ok {
		return nil, ErrUnauthorized
	}

	return user, nil
}
