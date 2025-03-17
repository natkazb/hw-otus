package storage

import (
	"errors"
	"time"
)

type Event struct {
	ID          int
	Title       string
	StartDate   time.Time
	EndDate     time.Time
	Description string
}

var (
	ErrCreateEvent = errors.New("error in creating new event")
	ErrDeleteEvent = errors.New("error in delete event")
)
