package domain

import "errors"

type Attendee struct {
	Id      int `json:"id"`
	UserId  int `json:"userId"`
	EventId int `json:"eventId"`
}

func (a *Attendee) Validate() error {
	if a.UserId == 0 || a.EventId == 0 {
		return errors.New("userId and eventId are required")
	}

	return nil
}
