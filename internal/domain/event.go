package domain

import (
	"errors"
	"time"
)

type Event struct {
	Id          int    `json:"id"`
	Name        string `json:"name" binding:"required,min=3"`
	Description string `json:"description" binding:"required,min=10"`
	Date        string `json:"date" binding:"required,datetime=2006-01-02"`
	Location    string `json:"location" binding:"required,min=3"`
	OwnerId     int    `json:"ownerId"`
}

func parseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

func validDate(date string) bool {
	_, err := parseDate(date)
	return err == nil
}

func (e *Event) Validate() error {
	if e.Name == "" || e.Description == "" || e.Location == "" {
		return errors.New("invalid event")
	}

	if e.OwnerId <= 0 {
		return errors.New("invalid owner")
	}

	if !validDate(e.Date) {
		return errors.New("invalid date")
	}
	return nil
}
