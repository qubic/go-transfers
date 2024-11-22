package main

import (
	"fmt"
	"github.com/ardanlabs/conf"
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
	if err := run(); err != nil {
		log.Fatalf("main: exited with error: %s", err.Error())
	}
}

func run() error {
	configuration, configErr := loadConfig()
	if configErr != nil {
		return errors.Wrap(configErr, "loading config")
	}

	srv := api.NewServer(configuration.Server.GrpcHost, configuration.Server.HttpHost)
	err := srv.Start()
	if err != nil {
		return errors.Wrap(err, "starting server")
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-shutdown:
			log.Println("main: shutting down...")
			return nil
		}
	}
}

func loadConfig() (*config.Config, error) {
	configuration, configErr := config.GetConfig()
	if configErr == nil {
		if out, toStringErr := conf.String(configuration); toStringErr == nil {
			slog.Info(fmt.Sprintf("Config :\n%v\n", out))
		}
	}
	return configuration, configErr
}
