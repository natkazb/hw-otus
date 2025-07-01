package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/config"                 //nolint
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/logger"                 //nolint
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/queue"                  //nolint
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/sender"              //nolint
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config_sender.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	conf, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Can't parse config file, %v", err)
		os.Exit(1)
	}

	logg := logger.New(conf.Logger.Level)

	q := queue.New(conf.Rabbit.Host,
		conf.Rabbit.Port,
		conf.Rabbit.User,
		conf.Rabbit.Password,
		conf.Rabbit.QueueName,
		conf.Rabbit.Timeout)

	send := sender.New(q, logg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	// defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := send.Stop(ctx); err != nil {
			logg.Error("failed to stop sender: " + err.Error())
		}
	}()

	err = send.Start(ctx)
	if err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1)
	}

	logg.Info("sender is running...")
	send.Run()
}
