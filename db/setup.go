package db

import (
	"fmt"
	"github.com/gookit/slog"
	"github.com/jmoiron/sqlx"
)

func Create(user, pass, dbName, host string, port, maxOpen, maxIdle int) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		host, port, user, pass, dbName)
	pgDb, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	pgDb.SetMaxOpenConns(maxOpen)
	pgDb.SetMaxIdleConns(maxIdle)
	slog.Info("Connected to database!")
	return pgDb, nil
}
