package internalgrpc

import (
	"context"
	"net"
	"time"

	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/server/grpc/pb" //nolint
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/storage"        //nolint
	"github.com/pkg/errors"
	"google.golang.org/grpc" //nolint
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
	server  *grpc.Server
	app     Application
}

type Application interface {
	CreateEvent(e storage.Event) (int32, error)
	UpdateEvent(e storage.Event) error
	DeleteEvent(id int32) error
	ListEvents(startData, endData time.Time) ([]storage.Event, error)
}

func NewServer(host, port string, logger Logger, app Application) *Server {
	newServ := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging(logger),
		),
	)
	return &Server{
		Address: net.JoinHostPort(host, port),
		log:     logger,
		app:     app,
		server:  newServ,
	}
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.Address)
	if err != nil {
		errWrapped := errors.Wrap(err, "start server error")
		s.log.Error(errWrapped.Error())
		return errWrapped
	}
	pb.RegisterEventServiceServer(s.server, &EventService{app: s.app})
	if err := s.server.Serve(lis); err != nil {
		errWrapped := errors.Wrap(err, "serve error")
		s.log.Error(errWrapped.Error())
		return errWrapped
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.server.GracefulStop()
	s.log.Info("server grpc stopped")
	return nil
}
