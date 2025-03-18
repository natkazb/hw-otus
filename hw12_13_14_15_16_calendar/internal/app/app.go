package app

import (
	"context"
	"fmt"
	"time"

	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/storage"
)

type App struct {
	log     Logger
	storage Storage
}

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	CreateEvent(e storage.Event) error
	UpdateEvent(e storage.Event) error
	DeleteEvent(id int) error
	ListEvents(startData, endData time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		log:     logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(_ context.Context, title string) error {
	now := time.Now()
	newEvent := storage.Event{Title: title, StartDate: now, EndDate: now.Add(time.Hour * 3), Description: "testing"}
	err := newEvent.Validate()
	if err != nil {
		a.log.Error(err.Error())
		return err
	}
	err = a.storage.CreateEvent(newEvent)
	if err != nil {
		a.log.Error(fmt.Errorf("%w: %w", storage.ErrCreateEvent, err).Error())
	}
	return err
}

func (a *App) DeleteEvent(_ context.Context, id int) error {
	err := a.storage.DeleteEvent(id)
	if err != nil {
		a.log.Error(fmt.Errorf("%w id=%d: %w", storage.ErrDeleteEvent, id, err).Error())
	}
	return err
}

func (a *App) UpdateEvent(_ context.Context, id int, title, description string, startDate, endDate time.Time) error {
	newEvent := storage.Event{ID: id, Title: title, StartDate: startDate, EndDate: endDate, Description: description}
	err := newEvent.Validate()
	if err != nil {
		a.log.Error(err.Error())
		return err
	}
	err = a.storage.UpdateEvent(newEvent)
	if err != nil {
		a.log.Error(fmt.Errorf("%w id=%d: %w", storage.ErrUpdateEvent, id, err).Error())
	}
	return err
}

func (a *App) ListEvents(startData, endData time.Time) ([]storage.Event, error) {
	list, err := a.storage.ListEvents(startData, endData)
	if err != nil {
		a.log.Error(fmt.Errorf("%w %v %v: %w", storage.ErrList, startData, endData, err).Error())
	}
	return list, err
}
