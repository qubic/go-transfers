package db

import (
	"embed"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/gookit/slog"
	"github.com/pkg/errors"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(connectionString string) error {
	dir, err := iofs.New(migrations, "migrations")
	if err != nil {
		return errors.Wrap(err, "loading embedded migrations")
	}
	migs, err := migrate.NewWithSourceInstance("iofs", dir, connectionString)
	if err != nil {
		return errors.Wrap(err, "creating migrations")
	}
	if err = migs.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "running migrations")
	}
	version, dirty, _ := migs.Version() // we don't care about error here. we only log info.
	slog.Info("Db migrations applied:", "version", version, "dirty", dirty,
		"changed", !errors.Is(err, migrate.ErrNoChange))
	sErr, dErr := migs.Close()
	slog.Info("Db migrations closed:", "source-errors", sErr, "db-errors", dErr)
	return nil
}
