package storage

import "time"

type Event struct {
	ID          string
	Title       string
	StartDate   time.Time
	EndDate     time.Time
	Description string
}
