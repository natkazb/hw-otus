package app

import (
	"context"
	"time"
	"fmt"

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
	//UpdateEvent(e storage.Event) error
	//DeleteEvent(e storage.Event) error
	//ListEvents(startData, endData time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		log:     logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(_ context.Context, title string) error {
	now := time.Now()
	err := a.storage.CreateEvent(storage.Event{Title: title, StartDate: now, EndDate: now.Add(1000), Description: "testing"})
	if (err != nil) {
		a.log.Error(fmt.Sprintf("error in CreateEvent: %v", err))
	}
	return err
}