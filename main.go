package main

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"go-transfers/api"
	"go-transfers/config"
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

	err = migrateDatabase(configuration)
	if err != nil {
		return errors.Wrap(err, "migrating database")
	}

	srv := api.NewServer(configuration.Server.GrpcHost, configuration.Server.HttpHost)
	err = srv.Start()
	if err != nil {
		return errors.Wrap(err, "starting server")
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

func migrateDatabase(configuration *config.Config) error {
	dbConfig := configuration.Database
	// run migrations
	m, err := migrate.New("file://db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s", dbConfig.User, dbConfig.Pass, dbConfig.Url))
	if err != nil {
		return errors.Wrap(err, "initializing migrate")
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "running migrations")
	} else {
		version, dirty, _ := m.Version() // we don't care about error here. we only log info.
		slog.Info("db migrations applied.", "Version", version, "Dirty", dirty)
	}
	return nil
}

func loadConfig() (*config.Config, error) {
	configuration, configErr := config.GetConfig()
	if configErr == nil {
		if out, toStringErr := conf.String(configuration); toStringErr == nil {
			slog.Info(fmt.Sprintf("Applied configuration properties:\n%v\n", out))
		}
	}
	return configuration, configErr
}
