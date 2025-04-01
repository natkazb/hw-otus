package storage

import (
	"errors"
	"time"
)

type EventDateTime time.Time

// type EventDateTime struct { time.Time }

const dateTimeLayout = "2006-01-02 15:04"

type Event struct {
	ID          int32
	Title       string        `json:"title"`
	StartDate   EventDateTime `db:"start_date" json:"startDate"`
	EndDate     EventDateTime `db:"end_date" json:"endDate"`
	Description string        `json:"description"`
}

type EventDB struct {
	ID          int32
	Title       string
	StartDate   time.Time `db:"start_date"`
	EndDate     time.Time `db:"end_date"`
	Description string
}

type EventModifyGrpc struct {
	ID          int32
	Title       string
	StartDate   string
	EndDate     string
	Description string
}

type DeleteEvent struct {
	ID int32 `json:"id"`
}

var (
	ErrCreateEvent      = errors.New("can't creating new event")
	ErrDeleteEvent      = errors.New("can't delete event")
	ErrUpdateEvent      = errors.New("can't update event")
	ErrDates            = errors.New("end date < start date")
	ErrList             = errors.New("can't select events")
	ErrEmptyID          = errors.New("empty id")
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

func (e *Event) ValidateUpdate() error {
	if e.ID == 0 {
		return ErrEmptyID
	}
	return e.Validate()
}

func (e *EventModifyGrpc) ValidateCreateGrpcAndReturnParsedDates() (EventDateTime, EventDateTime, error) {
	defaultTime := EventDateTime(time.Now())
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

func (ct EventDateTime) String() string {
	timeDay := time.Time(ct)
	return timeDay.String()
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

func (e *EventModifyGrpc) CopyToEvent(startDate, endDate EventDateTime) Event {
	eCopy := Event{
		ID:          e.ID,
		Title:       e.Title,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: e.Description,
	}
	return eCopy
}
