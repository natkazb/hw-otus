package storage

import (
	"errors"
	"time"
)

type EventDateTime time.Time

const dateTimeLayout = "2006-01-02 15:04"

type Event struct {
	ID          int
	Title       string        `json:"title"`
	StartDate   EventDateTime `db:"start_date" json:"startDate"`
	EndDate     EventDateTime `db:"end_date" json:"endDate"`
	Description string        `json:"description"`
}

type EventDB struct {
	ID          int
	Title       string
	StartDate   time.Time `db:"start_date"`
	EndDate     time.Time `db:"end_date"`
	Description string
}

var (
	ErrCreateEvent = errors.New("can't creating new event")
	ErrDeleteEvent = errors.New("can't delete event")
	ErrUpdateEvent = errors.New("can't update event")
	ErrDates       = errors.New("date end < date start")
	ErrList        = errors.New("can't select events")
)

func (ct *EventDateTime) UnmarshalJSON(b []byte) error {
	str := string(b)
	str = str[1 : len(str)-1] // Remove the quotes
	t, err := time.Parse(dateTimeLayout, str)
	if err != nil {
		return err
	}
	*ct = EventDateTime(t)
	return nil
}

func (e *Event) Validate() error {
	if e.StartDate.Compare(e.EndDate) > 0 {
		return ErrDates
	}
	return nil
}

func (ct EventDateTime) Compare(other EventDateTime) int {
	t1 := time.Time(ct)
	t2 := time.Time(other)
	if t1.Before(t2) {
		return -1
	} else if t1.After(t2) {
		return 1
	}
	return 0
}

func (e *Event) CopyToEventDB() EventDB {
	eDB := EventDB{
		ID:          e.ID,
		Title:       e.Title,
		StartDate:   time.Time(e.StartDate),
		EndDate:     time.Time(e.EndDate),
		Description: e.Description,
	}
	return eDB
}
