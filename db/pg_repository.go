package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go-transfers/config"
	"log/slog"
)

type PgRepository struct {
	db *sqlx.DB
}

func NewRepository(c *config.DatabaseConfig) (*PgRepository, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		c.Host, c.Port, c.User, c.Pass, c.Name)
	db, err := createDatabase(connectionString)
	if err != nil {
		return nil, err
	} else {
		return &PgRepository{db: db}, nil
	}
}

func (r *PgRepository) GetOrCreateEntity(identity string) (int64, error) {
	id, err := r.GetEntityId(identity)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return id, err
	}

	insertSql := `insert into entities (identity) values ($1) returning id;`
	err = r.db.QueryRow(insertSql, identity).Scan(&id)
	if err != nil {
		slog.Error("error inserting entity.", "identity", identity, "error", err)
		return id, err
	} else {
		slog.Debug("entity inserted successfully.", "id", id, "identity", identity)
		return id, nil
	}
}

func (r *PgRepository) GetEntityId(identity string) (int64, error) {
	selectSql := `select id from entities where identity= $1;`

	var id int64
	err := r.db.GetContext(context.Background(), &id, selectSql, identity)
	if err != nil {
		return id, err
	}

	if err := r.db.QueryRow(selectSql, identity).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Debug("no such entity.", "identity", identity)
		} else {
			slog.Error("error selecting entity.", "identity", identity, "err", err)
		}
		return id, err
	} else {
		return id, nil
	}
}

func createDatabase(connectionString string) (*sqlx.DB, error) {

	// open database
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// check db
	err = db.Ping()
	if err != nil {
		return db, err
	}

	slog.Info("Connected to database!")
	return db, nil
}

func (r *PgRepository) Close() {
	err := r.db.Close()
	if err != nil {
		slog.Error("error closing database.", "Error", err)
	} else {
		slog.Info("closed database.")
	}
}
