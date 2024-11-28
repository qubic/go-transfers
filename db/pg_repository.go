package db

import (
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
		db.SetMaxOpenConns(c.MaxOpen)
		db.SetMaxIdleConns(c.MaxIdle)
		return &PgRepository{db: db}, nil
	}
}

// entity

func (r *PgRepository) GetOrCreateEntity(identity string) (int, error) {
	id, err := r.getEntityId(identity)
	if errors.Is(err, sql.ErrNoRows) { // insert if not found
		id, err = r.insertEntity(identity)
	}
	return id, err
}

func (r *PgRepository) insertEntity(identity string) (int, error) {
	insertSql := `insert into entities (identity) values ($1) returning id;`
	return insert(r.db, insertSql, identity)
}

func (r *PgRepository) getEntityId(identity string) (int, error) {
	selectSql := `select id from entities where identity= $1;`
	return getId(r.db, selectSql, identity)
}

// asset

func (r *PgRepository) GetOrCreateAsset(issuer string, name string) (int, error) {
	id, err := r.getAssetId(issuer, name)
	if errors.Is(err, sql.ErrNoRows) { // not found create
		return r.createAsset(issuer, name)
	} else {
		return id, err
	}
}

func (r *PgRepository) createAsset(issuer string, name string) (int, error) {
	entityId, err := r.GetOrCreateEntity(issuer)
	if err != nil {
		return 0, err
	}
	return r.insertAsset(entityId, name)
}

func (r *PgRepository) getAssetId(issuer string, name string) (int, error) {
	selectSql := `select a.id from assets a
    join entities e on a.issuer_id = e.id
    where e.identity=$1 and a.name=$2;`

	return getId(r.db, selectSql, issuer, name)
}

func (r *PgRepository) insertAsset(issuerId int, name string) (int, error) {
	insertSql := `insert into assets (issuer_id, name) values ($1, $2) returning id;`
	return insert(r.db, insertSql, issuerId, name)
}

// event
//
//func (r *PgRepository) insertEvent(transactionId int, eventId int, eventType int, eventData string) (int, error) {
//	insertSql := `insert into events (transaction_id, event_id, event_type, event_data) values ($1, $2, $3, $4);`
//	return insert(r.db, insertSql, transactionId, eventId, eventType, eventData)
//}

func getId(db *sqlx.DB, statement string, args ...interface{}) (int, error) {
	var id int
	err := db.Get(&id, statement, args...)
	return id, err
}

func insert(db *sqlx.DB, statement string, args ...interface{}) (int, error) {
	var id int
	err := db.Get(&id, statement, args...)
	return id, err
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
