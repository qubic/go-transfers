package db

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
)

func MigrateDatabase(sourceUrl string, connectionString string) error {
	m, err := migrate.New(sourceUrl, connectionString)
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
