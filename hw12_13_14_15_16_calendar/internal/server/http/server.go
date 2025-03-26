package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/storage" //nolint
	"github.com/pkg/errors"
)

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warn(msg string)
	Error(msg string)
}

type Server struct {
	Address string
	log     Logger
	server  http.Server
	app     Application
}

type Application interface {
	CreateEvent(e storage.Event) error
	//UpdateEvent(e storage.Event) error
	//DeleteEvent(id int) error
	ListEvents(startData, endData time.Time) ([]storage.Event, error)
}

func NewServer(host, port string, logger Logger, app Application) *Server {
	return &Server{
		Address: net.JoinHostPort(host, port),
		log:     logger,
		app:     app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	eventHandler := &EventHandler{app: s.app}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /events/list-day", loggingMiddleware(eventHandler.ListDay, s.log))
	mux.HandleFunc("GET /events/list-week", loggingMiddleware(eventHandler.ListWeek, s.log))
	mux.HandleFunc("GET /events/list-month", loggingMiddleware(eventHandler.ListMonth, s.log))
	mux.HandleFunc("POST /event/create", loggingMiddleware(eventHandler.Create, s.log))
	s.server = http.Server{
		Addr:         s.Address,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "start server error")
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error(fmt.Sprintf("server shutdown error: %v", err))
		return err
	}
	s.log.Info("server stopped")
	return nil
}
