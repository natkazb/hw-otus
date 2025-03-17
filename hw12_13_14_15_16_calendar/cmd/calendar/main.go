package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/app"
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/config"
	"github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/logger"
	internalhttp "github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/server/http"
	sqlstorage "github.com/natkazb/hw-otus/hw12_13_14_15_16_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	/*if flag.Arg(0) == "version" {
		printVersion()
		return
	}*/

	conf, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Can't parse config file, %v", err)
		os.Exit(1)
	}

	logg := logger.New(conf.Logger.Level)

	//storage := memorystorage.New()
	storage := sqlstorage.New(conf.Storage.SQL.Driver, conf.Storage.SQL.Host, conf.Storage.SQL.Port, conf.Storage.SQL.DBName, conf.Storage.SQL.Username, conf.Storage.SQL.Password)
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(conf.HTTP.Host, conf.HTTP.Port, logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := storage.Close(ctx); err != nil {
			logg.Error("failed to stop storage: " + err.Error())
		}

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := storage.Connect(ctx); err != nil {
		logg.Error("failed to start storage: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
