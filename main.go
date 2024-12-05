package main

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"go-transfers/api"
	"go-transfers/client"
	"go-transfers/config"
	"go-transfers/db"
	"go-transfers/sync"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if runErr := run(); runErr != nil {
		log.Fatalf("main: exited with error: %s", runErr.Error())
	}
}

func run() error {
	// load config
	configuration, err := loadConfig()
	if err != nil {
		return errors.Wrap(err, "loading config")
	}

	// database
	err = migrateDatabase(&configuration.Database)
	if err != nil {
		return errors.Wrap(err, "migrating database")
	}
	pgDb, err := db.CreateDatabaseWithConfig(&configuration.Database)
	if err != nil {
		return errors.Wrap(err, "opening database")
	}
	defer pgDb.Close()
	repository := db.NewRepository(pgDb)

	// event processing
	eventProcessor := sync.NewEventProcessor(repository)
	eventClient, err := client.NewIntegrationEventClient(configuration.Client.EventApiUrl, configuration.Client.CoreApiUrl)
	if err != nil {
		return errors.Wrap(err, "creating event client")
	}
	eventService, err := sync.NewEventService(eventClient, eventProcessor, repository)
	if err != nil {
		return errors.Wrap(err, "creating event service")
	}

	if configuration.App.SyncEnabled {
		slog.Info("Starting sync...")
		go eventService.SyncInLoop()
	} else {
		slog.Info("Sync not enabled.")
	}

	if configuration.App.ApiEnabled {
		slog.Info("Starting api...")
		// api
		srv := api.NewServer(configuration.Server.GrpcHost, configuration.Server.HttpHost, repository)
		err = srv.Start()
		if err != nil {
			return errors.Wrap(err, "starting server")
		}
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-shutdown:
			slog.Info("main: shutting down...")
			return nil
		}
	}
}

func migrateDatabase(config *config.DatabaseConfig) error {
	// run migrations
	m, err := migrate.New("file://db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.User, config.Pass, config.Host, config.Port, config.Name))
	if err != nil {
		return errors.Wrap(err, "initializing migrate")
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "running migrations")
	} else {
		version, dirty, _ := m.Version() // we don't care about error here. we only log info.
		slog.Info("db migrations applied:", "version", version, "dirty", dirty,
			"changed", !errors.Is(err, migrate.ErrNoChange))
		sErr, dErr := m.Close()
		slog.Info("db migration close", "source-errors", sErr, "db-errors", dErr)
	}
	return nil
}

func loadConfig() (*config.Config, error) {
	configuration, configErr := config.GetConfig()
	if configErr == nil {
		if out, toStringErr := conf.String(configuration); toStringErr == nil {
			slog.Info(fmt.Sprintf("applied configuration properties.\n%v\n", out))
		}
	}
	return configuration, configErr
}
