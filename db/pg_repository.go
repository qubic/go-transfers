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

func (r *PgRepository) GetOrCreateAsset(issuer, name string) (int, error) {
	id, err := r.getAssetId(issuer, name)
	if errors.Is(err, sql.ErrNoRows) { // not found create
		return r.createAsset(issuer, name)
	} else {
		return id, err
	}
}

func (r *PgRepository) createAsset(issuer, name string) (int, error) {
	entityId, err := r.GetOrCreateEntity(issuer)
	if err != nil {
		return 0, err
	}
	return r.insertAsset(entityId, name)
}

func (r *PgRepository) getAssetId(issuer, name string) (int, error) {
	selectSql := `select a.id from assets a
    join entities e on a.issuer_id = e.id
    where e.identity=$1 and a.name=$2;`

	return getId(r.db, selectSql, issuer, name)
}

func (r *PgRepository) insertAsset(issuerId int, name string) (int, error) {
	insertSql := `insert into assets (issuer_id, name) values ($1, $2) returning id;`
	return insert(r.db, insertSql, issuerId, name)
}

// tick

func (r *PgRepository) GetOrCreateTick(tickNumber uint32) (int, error) {
	id, err := r.getTickId(tickNumber)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = r.insertTick(tickNumber)
	}
	return id, err
}

func (r *PgRepository) getTickId(tickNumber uint32) (int, error) {
	selectSql := `select id from ticks where tick_number = $1;`
	return getId(r.db, selectSql, tickNumber)
}

func (r *PgRepository) insertTick(tickNumber uint32) (int, error) {
	insertSql := `insert into ticks (tick_number) values ($1) returning id;`
	return insert(r.db, insertSql, tickNumber)
}

// transactions

func (r *PgRepository) GetOrCreateTransaction(hash string, tickId int) (int, error) {
	id, err := r.getTransactionId(hash, tickId)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = r.insertTransaction(hash, tickId)
	}
	return id, err
}

func (r *PgRepository) getTransactionId(hash string, tickId int) (int, error) {
	selectSql := `select id from transactions where hash = $1 and tick_id = $2;`
	return getId(r.db, selectSql, hash, tickId)
}

func (r *PgRepository) insertTransaction(hash string, tickId int) (int, error) {
	insertSql := `insert into transactions (hash, tick_id) values ($1, $2) returning id;`
	return insert(r.db, insertSql, hash, tickId)
}

// events

func (r *PgRepository) GetOrCreateEvent(transactionId int, eventEventId uint64, eventType uint32, eventData string) (int, error) {
	id, err := r.getEventId(transactionId, eventEventId)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = r.insertEvent(transactionId, eventEventId, eventType, eventData)
	}
	return id, err
}

func (r *PgRepository) insertEvent(transactionId int, eventEventId uint64, eventType uint32, eventData string) (int, error) {
	insertSql := `insert into events (transaction_id, event_id, event_type, event_data) values ($1, $2, $3, $4) returning id;`
	return insert(r.db, insertSql, transactionId, eventEventId, eventType, eventData)
}

func (r *PgRepository) getEventId(transactionId int, eventEventId uint64) (int, error) {
	selectSql := `select id from events where transaction_id = $1 and event_id = $2;`
	return getId(r.db, selectSql, transactionId, eventEventId)
}

// qu transfer events

func (r *PgRepository) GetOrCreateQuTransferEvent(eventId int, sourceEntityId int, destinationEntityId int, amount uint64) (int, error) {
	id, err := r.getQuTransferEventId(eventId)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = r.insertQuTransferEvent(eventId, sourceEntityId, destinationEntityId, amount)
	}
	return id, err
}

func (r *PgRepository) insertQuTransferEvent(eventId int, sourceEntityId int, destinationEntityId int, amount uint64) (int, error) {
	insertSql := `insert into qu_transfer_events (event_id, source_entity_id, destination_entity_id, amount) values ($1, $2, $3, $4) returning id;`
	return insert(r.db, insertSql, eventId, sourceEntityId, destinationEntityId, amount)
}

func (r *PgRepository) getQuTransferEventId(eventId int) (int, error) {
	selectSql := `select id from qu_transfer_events where event_id = $1;`
	return getId(r.db, selectSql, eventId)
}

// asset change events

func (r *PgRepository) GetOrCreateAssetChangeEvent(eventId, assetId, sourceEntityId, destinationEntityId int, numberOfShares int64) (int, error) {
	id, err := r.getAssetChangeEventId(eventId)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = r.insertAssetChangeEvent(eventId, assetId, sourceEntityId, destinationEntityId, numberOfShares)
	}
	return id, err
}

func (r *PgRepository) insertAssetChangeEvent(eventId, assetId, sourceEntityId, destinationEntityId int, numberOfShares int64) (int, error) {
	insertSql := `insert into asset_change_events (event_id, asset_id, source_entity_id, destination_entity_id, number_of_shares) values ($1, $2, $3, $4, $5) returning id;`
	return insert(r.db, insertSql, eventId, assetId, sourceEntityId, destinationEntityId, numberOfShares)
}

func (r *PgRepository) getAssetChangeEventId(eventId int) (int, error) {
	selectSql := `select id from asset_change_events where event_id = $1;`
	return getId(r.db, selectSql, eventId)
}

// helper methods

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
