package db

import (
	"database/sql"
	"fmt"
	"go-transfers/config"
	"log/slog"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(c *config.DatabaseConfig) (*Repository, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Pass, c.Name)
	db, err := createDatabase(connectionString)
	if err != nil {
		return nil, err
	} else {
		repo := &Repository{
			db: db,
		}
		return repo, nil
	}
}

func (r *Repository) insertEntity(identity string) (int64, error) {
	insertSQL := `INSERT INTO entities (identity) VALUES ($1);`
	result, err := r.db.Exec(insertSQL, identity)
	if err != nil {
		slog.Error("error inserting entity.", "identity", identity)
		return -1, err
	}
	id, err := result.LastInsertId()
	slog.Debug("entity inserted successfully.", "id", id, "identity", identity)
	return id, nil
}

func createDatabase(connectionString string) (*sql.DB, error) {

	// open database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// close database on exit
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			slog.Error("error closing database", "Error", err)
		}
	}(db)

	// check db
	err = db.Ping()
	if err != nil {
		return db, err
	}

	slog.Info("Connected to database!")
	return db, nil
}
