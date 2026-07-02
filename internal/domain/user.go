package domain

import (
	"errors"
	"regexp"
)

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

var validator = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (u *User) Validate() error {
	if u.Email == "" || !validator.MatchString(u.Email) {
		return errors.New("email is required and must be a valid email address")
	}

	if u.Name == "" || len(u.Name) > 128 {
		return errors.New("name is required and must be at most 128 characters")
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	if len(u.Password) < 8 || len(u.Password) > 24 {
		return errors.New("password must be at least 8 characters and at most 24")
	}

	return nil
}
