package main

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
	"github.com/pkg/errors"
	"go-transfers/api"
	"go-transfers/client"
	"go-transfers/config"
	"go-transfers/db"
	"go-transfers/sync"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// TODO add log config
// TODO embed migrations and .env file

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

	// logging
	configureLogging(configuration.Log)
	defer slog.MustClose()
	defer slog.MustFlush()

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

func configureLogging(config config.LogConfig) {
	const LogTimeFormat = "2006-01-02 15:04:05"

	// console
	logTemplate := "[{{datetime}}] [{{level}}] {{message}} {{data}} {{extra}}\n"
	formatter := slog.NewTextFormatter().WithEnableColor(true)
	formatter.SetTemplate(logTemplate)
	formatter.TimeFormat = LogTimeFormat
	slog.SetFormatter(formatter)
	logLevel := slog.LevelByName(config.Level)
	slog.SetLogLevel(logLevel)

	// error log
	h1 := handler.NewBuilder().
		WithLogfile("./error.log").
		WithLogLevels(slog.DangerLevels).
		WithRotateTime(rotatefile.EveryDay).
		WithBuffMode(handler.BuffModeLine).
		WithCompress(true).
		Build()
	errorFormatter := slog.NewTextFormatter()
	errorFormatter.TimeFormat = LogTimeFormat
	h1.SetFormatter(errorFormatter) // log complete data

	// normal application log
	h2 := handler.NewBuilder().
		WithLogfile("./application.log").
		WithLogLevels(slog.Levels{slog.InfoLevel, slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel}).
		WithRotateTime(rotatefile.EveryDay).
		WithBuffMode(handler.BuffModeLine).
		WithCompress(true).
		Build()
	infoFormatter := slog.NewTextFormatter()
	infoFormatter.TimeFormat = LogTimeFormat
	infoFormatter.SetTemplate(logTemplate)
	h2.SetFormatter(infoFormatter)

	slog.PushHandler(h1)
	slog.PushHandler(h2)

	slog.Info("Logging level set to", logLevel)
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
