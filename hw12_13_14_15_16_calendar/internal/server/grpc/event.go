package internalgrpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/server/grpc/pb" //nolint
)

var (
	ErrNoDay      = `{"error": "day parameter is required"}`
	ErrInvalidDay = `{"error": "invalid date format, use YYYY-MM-DD"}`
	DayListFormat = "2006-01-02"
)

type EventService struct {
	app     Application
	pb.UnimplementedEventServiceServer
}

func (h *EventService) CreateEvent(_ context.Context, req *pb.CreateEventRequest) (*pb.CreateEventResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "no data")
	}
	//Info(fmt.Sprintf("%v", req))
	//todo: дописать вставку в бд

	return &pb.CreateEventResponse{Id: 1257}, nil
}
