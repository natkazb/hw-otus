package storage

import (
	"errors"
	"time"
)

type Event struct {
	ID          int
	Title       string
	StartDate   time.Time `db:"start_date"`
	EndDate     time.Time `db:"end_date"`
	Description string
}

var (
	ErrCreateEvent = errors.New("error in creating new event")
	ErrDeleteEvent = errors.New("error in delete event")
	ErrUpdateEvent = errors.New("error in update event")
	ErrDates       = errors.New("error date end < date start")
	ErrList        = errors.New("error in select data")
)

func (e *Event) Validate() error {
	if e.StartDate.Compare(e.EndDate) > 0 {
		return ErrDates
	}
	return nil
}
