package sqlstorage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Add this import for Postgres driver
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/storage"
)

type Storage struct {
	Dsn    string
	Driver string
	db     *sqlx.DB
}

func New(driver, host string, port int, dbName, user, pass string) *Storage {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbName)
	//dsn := fmt.Sprintf("%s://%s:%s@%s:%v/%s?sslmode=disable", driver, user, pass, host, port, dbName)
	return &Storage{
		Dsn:    dsn,
		Driver: driver,
	}
}

func (s *Storage) Connect(_ context.Context) (err error) {
	s.db, err = sqlx.Connect(s.Driver, s.Dsn)
	return
}

func (s *Storage) Close(_ context.Context) error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Storage) CreateEvent(e storage.Event) error {
	_, err := s.db.Exec("INSERT INTO event (title, start_date, end_date, description, user_id, notify_on) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		e.Title,
		e.StartDate,
		e.EndDate,
		e.Description,
		1, // эти поля пока не реализуем
		1,
	)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrCreateEvent, err)
	}
	return nil
}

func (s *Storage) DeleteEvent(id int) error {
	_, err := s.db.Exec("DELETE FROM event WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("%w: %w", storage.ErrDeleteEvent, err)
	}
	return nil
}
