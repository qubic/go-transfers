package main

import (
	"fmt"
	"github.com/ardanlabs/conf"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
	"github.com/pkg/errors"
	"go-transfers/api"
	"go-transfers/client"
	"go-transfers/db"
	"go-transfers/sync"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if runErr := run(); runErr != nil {
		log.Fatalf("main: exited with error: %s", runErr.Error())
	}
}

type ServerConfig struct {
	HttpHost string `conf:"default:0.0.0.0:8000"`
	GrpcHost string `conf:"default:0.0.0.0:8001"`
}

type ClientConfig struct {
	EventApiUrl string `conf:"required"`
	CoreApiUrl  string `conf:"required"`
}

type DatabaseConfig struct {
	User    string `conf:"default:qxtr"`
	Pass    string `conf:"noprint,required"`
	Host    string `conf:"default:localhost"`
	Port    int    `conf:"default:5432"`
	Name    string `conf:"default:qxtr"`
	MaxIdle int    `conf:"default:10"`
	MaxOpen int    `conf:"default:10"`
}

type AppConfig struct {
	SyncEnabled bool `conf:"default:true"`
	ApiEnabled  bool `conf:"default:true"`
}

type LogConfig struct {
	Level     string `conf:"default:Info"`
	FileError bool   `conf:"default:false"`
	FileApp   bool   `conf:"default:false"`
}

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Client   ClientConfig
	Database DatabaseConfig
	Log      LogConfig
}

const envPrefix = "QUBIC_TRANSFERS"

func run() error {

	// load config
	var configuration Config
	if err := conf.Parse(os.Args[1:], envPrefix, &configuration); err != nil {
		switch {
		case errors.Is(err, conf.ErrHelpWanted):
			usage, err := conf.Usage(envPrefix, &configuration)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		case errors.Is(err, conf.ErrVersionWanted):
			version, err := conf.VersionString(envPrefix, &configuration)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&configuration)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config :\n%v\n", out)

	// logging
	configureLogging(configuration.Log)
	defer slog.MustClose()
	defer slog.MustFlush()

	// database
	err = migrateDatabase(&configuration.Database)
	if err != nil {
		return errors.Wrap(err, "migrating database")
	}

	dbc := configuration.Database
	pgDb, err := db.Create(dbc.User, dbc.Pass, dbc.Name, dbc.Host, dbc.Port, dbc.MaxOpen, dbc.MaxIdle)
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

func configureLogging(config LogConfig) {
	const LogTimeFormat = "2006-01-02 15:04:05"

	// console
	logTemplate := "[{{datetime}}] [{{level}}] {{message}} {{data}} {{extra}}\n"
	formatter := slog.NewTextFormatter().WithEnableColor(true)
	formatter.SetTemplate(logTemplate)
	formatter.TimeFormat = LogTimeFormat
	slog.SetFormatter(formatter)
	logLevel := slog.LevelByName(config.Level)
	slog.SetLogLevel(logLevel)

	if config.FileApp {
		// normal application log
		appLogHandler := handler.NewBuilder().
			WithLogfile("./application.log").
			WithLogLevels(slog.Levels{slog.InfoLevel, slog.PanicLevel, slog.FatalLevel, slog.ErrorLevel, slog.WarnLevel}).
			WithRotateTime(rotatefile.EveryDay).
			WithBuffMode(handler.BuffModeLine).
			WithCompress(true).
			Build()
		infoFormatter := slog.NewTextFormatter()
		infoFormatter.TimeFormat = LogTimeFormat
		infoFormatter.SetTemplate(logTemplate)
		appLogHandler.SetFormatter(infoFormatter)
		slog.PushHandler(appLogHandler)
		slog.Info("Enabled application log.")
	}

	if config.FileError {
		// error log
		errLogHandler := handler.NewBuilder().
			WithLogfile("./error.log").
			WithLogLevels(slog.DangerLevels).
			WithRotateTime(rotatefile.EveryDay).
			WithBuffMode(handler.BuffModeLine).
			WithCompress(true).
			Build()
		errorFormatter := slog.NewTextFormatter()
		errorFormatter.TimeFormat = LogTimeFormat
		errLogHandler.SetFormatter(errorFormatter) // log complete data
		slog.PushHandler(errLogHandler)
		slog.Info("Enabled error log.")
	}

	slog.Info("Log level set:", logLevel)
}

func migrateDatabase(config *DatabaseConfig) error {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", config.User, config.Pass, config.Host, config.Port, config.Name)
	return db.Migrate(connectionString)
}
