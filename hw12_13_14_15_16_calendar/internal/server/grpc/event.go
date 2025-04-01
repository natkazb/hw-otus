package internalgrpc

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc/codes"  //nolint
	"google.golang.org/grpc/status" //nolint

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
	event := storage.EventModifyGrpc{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		StartDate:   req.GetStartDate(),
		EndDate:     req.GetEndDate(),
	}
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

func (h *EventService) UpdateEvent(_ context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	event := storage.EventModifyGrpc{
		ID:          req.GetId(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		StartDate:   req.GetStartDate(),
		EndDate:     req.GetEndDate(),
	}
	startDateParsed, endDateParsed, err := event.ValidateCreateGrpcAndReturnParsedDates()
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err = h.app.UpdateEvent(event.CopyToEvent(startDateParsed, endDateParsed)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateEventResponse{}, nil
}

func (h *EventService) ListDay(_ context.Context, req *pb.ListDayEventRequest) (*pb.ListDayEventResponse, error) {
	dayParsed, err := validateAndReturnParsedDate(req.GetDay())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	items, err := h.app.ListEvents(dayParsed, dayParsed.Add(24*time.Hour))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	responseItems := make([]*pb.EventList, len(items))
	for i, v := range items {
		responseItems[i] = &pb.EventList{
			Id:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			StartDate:   v.StartDate.String(),
			EndDate:     v.EndDate.String(),
		}
	}
	return &pb.ListDayEventResponse{Events: responseItems}, nil
}

func (h *EventService) ListWeek(_ context.Context, req *pb.ListWeekEventRequest) (*pb.ListWeekEventResponse, error) {
	dayParsed, err := validateAndReturnParsedDate(req.GetDay())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	items, err := h.app.ListEvents(dayParsed, dayParsed.Add(7*24*time.Hour))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	responseItems := make([]*pb.EventList, len(items))
	for i, v := range items {
		responseItems[i] = &pb.EventList{
			Id:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			StartDate:   v.StartDate.String(),
			EndDate:     v.EndDate.String(),
		}
	}
	return &pb.ListWeekEventResponse{Events: responseItems}, nil
}

func (h *EventService) ListMonth(_ context.Context, req *pb.ListMonthEventRequest) (*pb.ListMonthEventResponse, error) {
	dayParsed, err := validateAndReturnParsedDate(req.GetDay())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	items, err := h.app.ListEvents(dayParsed, dayParsed.AddDate(0, 1, 0))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	responseItems := make([]*pb.EventList, len(items))
	for i, v := range items {
		responseItems[i] = &pb.EventList{
			Id:          v.ID,
			Title:       v.Title,
			Description: v.Description,
			StartDate:   v.StartDate.String(),
			EndDate:     v.EndDate.String(),
		}
	}
	return &pb.ListMonthEventResponse{Events: responseItems}, nil
}

func validateAndReturnParsedDate(day string) (time.Time, error) {
	defaultDay := time.Now()
	if day == "" {
		return defaultDay, errors.New("empty day")
	}

	dayParsed, err := time.Parse(DayListFormat, day)
	if err != nil {
		return defaultDay, errors.New("incorrect format day YYYY-MM-DDM")
	}
	return dayParsed, nil
}
