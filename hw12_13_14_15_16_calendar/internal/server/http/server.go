package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

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
}

type Application interface { // TODO
}

type DefaultHandler struct{}

func (h *DefaultHandler) Info(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Hello world!"))
	//time.Sleep(time.Second)
}

func NewServer(host, port string, logger Logger, app Application) *Server {
	return &Server{
		Address: net.JoinHostPort(host, port),
		log:     logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	handler := &DefaultHandler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/info", loggingMiddleware(handler.Info, s.log))
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
