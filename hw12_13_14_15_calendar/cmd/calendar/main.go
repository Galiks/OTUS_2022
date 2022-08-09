package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Galiks/OTUS_2022/hw12_13_14_15_calendar/internal/app"
	"github.com/Galiks/OTUS_2022/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/Galiks/OTUS_2022/hw12_13_14_15_calendar/internal/server/http"

	// memorystorage "github.com/Galiks/OTUS_2022/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Galiks/OTUS_2022/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}
	if err := logger.InitLog(config.Logger.Level, config.Logger.PrintStackTrace, config.Logger.PathToFile); err != nil {
		log.Fatal(err)
	}
	storage := sqlstorage.New(config.PostgreSQL.ConnectString)
	if err := storage.Connect(context.Background()); err != nil {
		logger.Fatal(err)
	}
	calendar := app.New(storage)

	server := internalhttp.NewServer(calendar, config.Server.Host, config.Server.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	if err := server.Start(); err != nil {
		logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
