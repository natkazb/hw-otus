package app

import (
	"context"
	"fmt"
	"time"

	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/storage" //nolint
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
	CreateEvent(e storage.EventDB) error
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

func (a *App) CreateEvent(event storage.Event) error {
	err := event.Validate()
	if err != nil {
		a.log.Error(err.Error())
		return err
	}
	err = a.storage.CreateEvent(event.CopyToEventDB())
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

func (a *App) UpdateEvent(event storage.Event) error {
	err := event.Validate()
	if err != nil {
		a.log.Error(err.Error())
		return err
	}
	err = a.storage.UpdateEvent(event)
	if err != nil {
		a.log.Error(fmt.Errorf("%w id=%d: %w", storage.ErrUpdateEvent, event.ID, err).Error())
	}
	return err
}

func (a *App) ListEvents(startData, endData time.Time) ([]storage.Event, error) {
	a.log.Info(fmt.Sprintf("get ListEvents from %v to %v", startData, endData))
	list, err := a.storage.ListEvents(startData, endData)
	if err != nil {
		a.log.Error(fmt.Errorf("%w %v %v: %w", storage.ErrList, startData, endData, err).Error())
	}
	return list, err
}
