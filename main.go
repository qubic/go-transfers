package main

import (
	"fmt"
	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"go-transfers/api"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const prefix = "QUBIC_TRANSFERS"

func main() {
	if err := run(); err != nil {
		log.Fatalf("main: exited with error: %s", err.Error())
	}
}

func run() error {

	var config struct {
		Server struct {
			HttpHost string `conf:"default:0.0.0.0:8000"`
			GrpcHost string `conf:"default:0.0.0.0:8001"`
		}
	}

	if err := conf.Parse(os.Args[1:], prefix, &config); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			usage, err := conf.Usage(prefix, &config)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		if errors.Is(err, conf.ErrVersionWanted) {
			version, err := conf.VersionString(prefix, &config)
			if err != nil {
				return errors.Wrap(err, "generating config version")
			}
			fmt.Println(version)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	out, err := conf.String(&config)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("main: Config :\n%v\n", out)
	srv := api.NewServer(config.Server.GrpcHost, config.Server.HttpHost)
	err = srv.Start()
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
