package sender

import (
	"context"

	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/app"   //nolint
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/queue" //nolint
)

type Sender struct {
	q      *queue.RabbitMq
	ctx    context.Context
	logger app.Logger
}

func New(q *queue.RabbitMq, logger app.Logger) Sender {
	return Sender{
		q:      q,
		logger: logger,
	}
}

func (s *Sender) Run() {
	ch, err := s.q.Consume()
	if err != nil {
		s.logger.Error(err.Error())
		return
	}
	for {
		select {
		case msg := <-ch:
			s.logger.Info("received event" + string(msg.Body))
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Sender) Start(ctx context.Context) error {
	s.ctx = ctx

	return s.q.Start(ctx)
}

func (s *Sender) Stop(ctx context.Context) error {
	return s.q.Stop(ctx)
}
