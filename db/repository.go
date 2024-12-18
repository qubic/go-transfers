package db

import (
	"context"
	"github.com/gookit/slog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PgRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *PgRepository {
	repo := PgRepository{db: db}
	return &repo
}

// helper methods

func getId(ctx context.Context, db *sqlx.DB, statement string, args ...interface{}) (int, error) {
	var id int
	err := db.GetContext(ctx, &id, statement, args...)
	return id, err
}

func insert(ctx context.Context, db *sqlx.DB, statement string, args ...interface{}) (int, error) {
	var id int
	err := db.GetContext(ctx, &id, statement, args...)
	return id, err
}

func (r *PgRepository) Close() {
	err := r.db.Close()
	if err != nil {
		slog.Error("error closing database.", "Error", err)
	} else {
		slog.Info("closed database.")
	}
}
