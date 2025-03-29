package internalgrpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/server/grpc/pb" //nolint
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/storage"        //nolint
)

var (
	ErrNoDay      = `{"error": "day parameter is required"}`
	ErrInvalidDay = `{"error": "invalid date format, use YYYY-MM-DD"}`
	DayListFormat = "2006-01-02"
)

type EventService struct {
	app Application
	pb.UnimplementedEventServiceServer
}

func (h *EventService) CreateEvent(_ context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	var id int32
	event := storage.EventCreateGrpc{Title: req.GetTitle(), Description: req.GetDescription(), StartDate: req.GetStartDate(), EndDate: req.GetEndDate()}
	startDateParsed, endDateParsed, err := event.ValidateCreateGrpcAndReturnParsedDates()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if id, err = h.app.CreateEvent(event.CopyToEvent(startDateParsed, endDateParsed)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateEventResponse{Id: id}, nil
}

func (h *EventService) DeleteEvent(_ context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	id := req.GetId()
	if id == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty id")
	}

	if err := h.app.DeleteEvent(id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteEventResponse{}, nil
}
