package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Su1203/golang-hw/hw12_13_14_15_16_calendar/internal/app"
	"github.com/Su1203/golang-hw/hw12_13_14_15_16_calendar/internal/logger"
	internalhttp "github.com/Su1203/golang-hw/hw12_13_14_15_16_calendar/internal/server/http"
	memorystorage "github.com/Su1203/golang-hw/hw12_13_14_15_16_calendar/internal/storage/memory"
	sqlstorage "github.com/Su1203/golang-hw/hw12_13_14_15_16_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	logg := logger.New(config.Logger.Level)
	logg.Info("Starting calendar service...")
	logg.Info("Config loaded from: " + configFile)

	var storage app.Storage
	switch config.Storage.Type {
	case "sql":
		logg.Info("Using SQL storage")
		sqlStorage, err := sqlstorage.New(config.Database.DSN)
		if err != nil {
			logg.Error("Failed to initialize SQL storage: " + err.Error())
			os.Exit(1)
		}
		defer func() {
			if err := sqlStorage.Close(); err != nil {
				logg.Error("Failed to close SQL storage: " + err.Error())
			}
		}()
		storage = sqlStorage
	case "memory":
		logg.Info("Using in-memory storage")
		storage = memorystorage.New()
	default:
		logg.Error("Unknown storage type: " + config.Storage.Type)
		os.Exit(1)
	}

	calendar := app.New(logg, storage)
	server := internalhttp.NewServer(logg, calendar, config.HTTP.Host, config.HTTP.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()
		logg.Info("Shutdown signal received")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*3)
		defer shutdownCancel()

		if err := server.Stop(shutdownCtx); err != nil {
			logg.Error("Failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("Calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("Failed to start http server: " + err.Error())
		cancel()
		os.Exit(1)
	}

	logg.Info("Calendar service stopped gracefully")
}
