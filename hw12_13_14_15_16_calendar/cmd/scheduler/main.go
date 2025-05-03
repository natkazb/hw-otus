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
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/scheduler"              //nolint
	sqlstorage "github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/storage/sql" //nolint
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	conf, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Can't parse config file, %v", err)
		os.Exit(1)
	}

	logg := logger.New(conf.Logger.Level)

	storage := sqlstorage.New(conf.Storage.SQL.Driver,
		conf.Storage.SQL.Host,
		conf.Storage.SQL.Port,
		conf.Storage.SQL.DBName,
		conf.Storage.SQL.Username,
		conf.Storage.SQL.Password)

	q := queue.New(conf.Rabbit.Host,
		conf.Rabbit.Port,
		conf.Rabbit.User,
		conf.Rabbit.Password,
		conf.Rabbit.QueueName,
		conf.Rabbit.Timeout)

	sch := scheduler.New(conf.Rabbit.Timeout, q, storage, logg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	// defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := sch.Stop(ctx); err != nil {
			logg.Error("failed to stop scheduler: " + err.Error())
		}
	}()

	logg.Info("scheduler is running...")

	err = sch.Start(ctx)
	if err != nil {
		logg.Error(err.Error())
		cancel()
		os.Exit(1)
	}
}
