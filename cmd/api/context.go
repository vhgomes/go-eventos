package main

import (
	"errors"
	"vhgomes-eventos/internal/database"

	"github.com/gin-gonic/gin"
)

var ErrUserNotFound = errors.New("user not found")

func (app *application) GetUserFromContext(c *gin.Context) (*database.User, error) {
	contextUser, exists := c.Get("user")

	if !exists {
		return nil, ErrUserNotFound
	}

	user, ok := contextUser.(*database.User)
	if !ok {
		return nil, ErrUserNotFound
	}

	return user, nil
}
