package scheduler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/app"     //nolint
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/queue"   //nolint
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/storage" //nolint
)

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Notify() ([]storage.Event, error)
}

type Scheduler struct {
	period    time.Duration
	monthsOld int
	q         *queue.RabbitMq
	storage   Storage
	logger    app.Logger
	ctx       context.Context
}

func New(period time.Duration, monthsOld int, q *queue.RabbitMq, storage Storage, logger app.Logger) Scheduler {
	return Scheduler{
		period:    period,
		monthsOld: monthsOld,
		q:         q,
		storage:   storage,
		logger:    logger,
	}
}

func (s *Scheduler) Run() {
	ticker := time.NewTicker(s.period)
	for {
		select {
		case <-ticker.C:
			s.logger.Info("start scheduler send to rabbitmq")
			events, err := s.storage.Notify()
			if err != nil {
				s.logger.Error(err.Error())
			} else {
				s.Push(events)
			}
			s.logger.Info("end scheduler send to rabbitmq")
			// @todo: удалить из бд старше s.monthsOld
		case <-s.ctx.Done():
			s.logger.Info("scheduler is stopping")
			return
		}
	}
}

func (s *Scheduler) Push(events []storage.Event) {
	for _, v := range events {
		js, err := json.Marshal(v)
		if err != nil {
			s.logger.Error(err.Error())
			continue
		}

		err = s.q.Produce(js, "application/json")
		if err != nil {
			s.logger.Error(err.Error())
		}
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.ctx = ctx

	err := s.storage.Connect(ctx)
	if err != nil {
		return err
	}

	err = s.q.Start(ctx)

	return err
}

func (s *Scheduler) Stop(ctx context.Context) error {
	err := s.storage.Close(ctx)
	if err != nil {
		return err
	}

	err = s.q.Stop(ctx)

	return err
}
