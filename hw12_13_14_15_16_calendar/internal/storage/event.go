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

type EventCreateGrpc struct {
	Title       string
	StartDate   string
	EndDate     string
	Description string
}

var (
	ErrCreateEvent      = errors.New("can't creating new event")
	ErrDeleteEvent      = errors.New("can't delete event")
	ErrUpdateEvent      = errors.New("can't update event")
	ErrDates            = errors.New("end date < start date")
	ErrList             = errors.New("can't select events")
	ErrEmptyTitle       = errors.New("empty title")
	ErrEmptyDescription = errors.New("empty description")
	ErrEmptyStartDate   = errors.New("empty start date")
	ErrEmptyEndDate     = errors.New("empty end date")
	ErrParseStartDate   = errors.New("incorrect format start date YYYY-MM-DD HH:MM")
	ErrParseEndDate     = errors.New("incorrect format end date YYYY-MM-DD HH:MM")
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

func ParseToEventDateTime(s string) (EventDateTime, error) {
	t, err := time.Parse(dateTimeLayout, s)
	if err == nil {
		t2 := EventDateTime(t)
		return t2, nil
	}
	return EventDateTime(time.Now()), err
}

func (e *Event) Validate() error {
	if e.Title == "" {
		return ErrEmptyTitle
	}
	if e.Description == "" {
		return ErrEmptyDescription
	}
	if e.StartDate.Compare(e.EndDate) > 0 {
		return ErrDates
	}
	return nil
}

func (e *EventCreateGrpc) ValidateCreateGrpcAndReturnParsedDates() (EventDateTime, EventDateTime, error) {
	defaultTime := EventDateTime(time.Now())
	if e.Title == "" {
		return defaultTime, defaultTime, ErrEmptyTitle
	}
	if e.Description == "" {
		return defaultTime, defaultTime, ErrEmptyDescription
	}
	if e.StartDate == "" {
		return defaultTime, defaultTime, ErrEmptyStartDate
	}
	if e.EndDate == "" {
		return defaultTime, defaultTime, ErrEmptyEndDate
	}
	startDateParsed, err := ParseToEventDateTime(e.StartDate)
	if err != nil {
		return defaultTime, defaultTime, ErrParseStartDate
	}
	endDateParsed, err := ParseToEventDateTime(e.EndDate)
	if err != nil {
		return defaultTime, defaultTime, ErrParseEndDate
	}
	if startDateParsed.Compare(endDateParsed) > 0 {
		return defaultTime, defaultTime, ErrDates
	}
	return startDateParsed, endDateParsed, nil
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

func (e *EventCreateGrpc) CopyToEvent(startDate, endDate EventDateTime) Event {
	eCopy := Event{
		Title:       e.Title,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: e.Description,
	}
	return eCopy
}
