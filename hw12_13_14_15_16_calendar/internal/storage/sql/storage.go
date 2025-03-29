package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"                                                   // Add this import for Postgres driver
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/storage" //nolint
)

type Storage struct {
	Dsn    string
	Driver string
	db     *sqlx.DB
}

func New(driver, host string, port int, dbName, user, pass string) *Storage {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbName)
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

func (s *Storage) CreateEvent4(e storage.EventDB) error {
	res, err := s.db.Exec(`INSERT INTO event 
	(title, start_date, end_date, description, user_id, notify_on) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id`,
		e.Title,
		e.StartDate,
		e.EndDate,
		e.Description,
		1, // эти поля пока не реализуем
		1,
	)
	fmt.Printf("%v", res)
	return err
}

func (s *Storage) CreateEvent(e storage.EventDB) (int32, error) {
	var id int32
	err := s.db.QueryRow(`INSERT INTO event 
	(title, start_date, end_date, description, user_id, notify_on) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	RETURNING id`,
		e.Title,
		e.StartDate,
		e.EndDate,
		e.Description,
		1, // эти поля пока не реализуем
		1,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) DeleteEvent(id int32) error {
	_, err := s.db.Exec("DELETE FROM event WHERE id = $1", id)
	return err
}

func (s *Storage) UpdateEvent(e storage.Event) error {
	_, err := s.db.Exec(`UPDATE event SET 
	title = $2, 
	start_date = $3, 
	end_date = $4, 
	description = $5 
	WHERE id = $1`,
		e.ID,
		e.Title,
		e.StartDate,
		e.EndDate,
		e.Description,
	)
	return err
}

func (s *Storage) ListEvents(startData, endData time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)
	err := s.db.Select(&events, `
SELECT id, title, start_date, end_date, description
FROM event
WHERE start_date >= $1 AND start_date < $2`,
		startData, endData)
	return events, err
}
