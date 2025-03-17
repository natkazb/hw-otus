package sqlstorage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Add this import for Postgres driver
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
