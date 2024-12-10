package db

import (
	"database/sql"
	"fmt"
	"github.com/gookit/slog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go-transfers/config"
	"go-transfers/proto"
)

type PgRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *PgRepository {
	repo := PgRepository{db: db}
	return &repo
}

// qu transfer events

func (r *PgRepository) GetQuTransferEventsForTick(tickNumber int) ([]*proto.QuTransferEvent, error) {
	selectSql := `select src.identity sourceId, 
       		dst.identity destinationId,
       		ev.amount,
       		tx.hash transactionHash,
       		ti.tick_number tick,
       		e.event_type eventType
		from qu_transfer_events ev
		join events e on ev.event_id = e.id
		join transactions tx on e.transaction_id = tx.id
		join ticks ti on tx.tick_id = ti.id
		join entities src on ev.source_entity_id = src.id
		join entities dst on ev.destination_entity_id = dst.id
		where ti.tick_number = $1 and e.event_type = 0
		order by tick_number desc;`
	var events []*proto.QuTransferEvent
	err := r.db.Select(&events, selectSql, tickNumber)
	if err != nil {
		return nil, errors.Wrap(err, "getting asset change events")
	}
	return events, nil
}

func (r *PgRepository) GetQuTransferEventsForEntity(identity string) ([]*proto.QuTransferEvent, error) {
	selectSql := `select src.identity sourceId, 
       		dst.identity destinationId,
       		ev.amount,
       		tx.hash transactionHash,
       		ti.tick_number tick,
       		e.event_type eventType
		from qu_transfer_events ev
		join events e on ev.event_id = e.id
		join transactions tx on e.transaction_id = tx.id
		join ticks ti on tx.tick_id = ti.id
		join entities src on ev.source_entity_id = src.id
		join entities dst on ev.destination_entity_id = dst.id
		where e.event_type = 0
		and (src.identity = $1 or dst.identity = $1)
		order by tick_number desc
		limit 100;`
	var events []*proto.QuTransferEvent
	err := r.db.Select(&events, selectSql, identity)
	if err != nil {
		return nil, errors.Wrap(err, "getting asset change events")
	}
	return events, nil
}

// asset change events

func (r *PgRepository) GetAssetChangeEventsForTick(tickNumber int) ([]*proto.AssetChangeEvent, error) {
	selectSql := `select src.identity sourceId, 
       		dst.identity destinationId, 
       		issuer.identity issuerId,
       		a.name, 
       		ev.number_of_shares numberOfShares,
       		tx.hash transactionHash,
       		ti.tick_number tick,
       		e.event_type eventType
		from asset_change_events ev
		join events e on ev.event_id = e.id
		join transactions tx on e.transaction_id = tx.id
		join ticks ti on tx.tick_id = ti.id
		join assets a on ev.asset_id = a.id
		join entities issuer on a.issuer_id = issuer.id
		join entities src on ev.source_entity_id = src.id
		join entities dst on ev.destination_entity_id = dst.id
		where ti.tick_number = $1 
		and e.event_type in (2, 3)
		order by tick_number desc;`
	var events []*proto.AssetChangeEvent
	err := r.db.Select(&events, selectSql, tickNumber)
	if err != nil {
		return nil, errors.Wrap(err, "getting asset change events")
	}
	return events, nil
}

func (r *PgRepository) GetAssetChangeEventsForEntity(identity string) ([]*proto.AssetChangeEvent, error) {
	selectSql := `select src.identity sourceId, 
       		dst.identity destinationId, 
       		issuer.identity issuerId,
       		a.name, 
       		ev.number_of_shares numberOfShares,
       		tx.hash transactionHash,
       		ti.tick_number tick,
       		e.event_type eventType
		from asset_change_events ev
		join events e on ev.event_id = e.id
		join transactions tx on e.transaction_id = tx.id
		join ticks ti on tx.tick_id = ti.id
		join assets a on ev.asset_id = a.id
		join entities issuer on a.issuer_id = issuer.id
		join entities src on ev.source_entity_id = src.id
		join entities dst on ev.destination_entity_id = dst.id
		where e.event_type in (2, 3)
		and (src.identity = $1 or dst.identity = $1)
		order by tick_number desc;`
	var events []*proto.AssetChangeEvent
	err := r.db.Select(&events, selectSql, identity)
	if err != nil {
		return nil, errors.Wrap(err, "getting asset change events")
	}
	return events, nil
}

// key value

func (r *PgRepository) GetLatestTick() (int, error) {
	return r.getNumericValue("tick")
}

func (r *PgRepository) UpdateLatestTick(tickNumber int) error {
	return r.updateNumericValue("tick", tickNumber)
}

func (r *PgRepository) getNumericValue(key string) (int, error) {
	selectSql := `select numeric_value from key_values where key = $1`
	var value int
	err := r.db.Get(&value, selectSql, key)
	return value, errors.Wrap(err, "getting numeric value")
}

func (r *PgRepository) updateNumericValue(key string, value int) error {
	updateSql := `update key_values set numeric_value = $1 where key = $2`
	_, err := r.db.Exec(updateSql, value, key)
	return errors.Wrap(err, "updating numeric value")
}

// entity

func (r *PgRepository) GetOrCreateEntity(identity string) (int, error) {
	id, err := r.getEntityId(identity)
	if errors.Is(err, sql.ErrNoRows) { // insert if not found
		id, err = r.insertEntity(identity)
	}
	return id, errors.Wrap(err, "getting or creating entity")
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
		return id, errors.Wrap(err, "getting or creating asset")
	}
}

func (r *PgRepository) createAsset(issuer, name string) (int, error) {
	entityId, err := r.GetOrCreateEntity(issuer)
	if err != nil {
		return 0, errors.Wrap(err, "creating asset")
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
	return id, errors.Wrap(err, "getting or creating tick")
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
	return id, errors.Wrap(err, "getting or creating transaction")
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
	return id, errors.Wrap(err, "getting or creating event")
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
	return id, errors.Wrap(err, "getting or creating qu transfer event")
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
	return id, errors.Wrap(err, "getting or creating asset change event")
}

func (r *PgRepository) insertAssetChangeEvent(eventId, assetId, sourceEntityId, destinationEntityId int, numberOfShares int64) (int, error) {
	insertSql := `insert into asset_change_events (event_id, asset_id, source_entity_id, destination_entity_id, number_of_shares) values ($1, $2, $3, $4, $5) returning id;`
	return insert(r.db, insertSql, eventId, assetId, sourceEntityId, destinationEntityId, numberOfShares)
}

func (r *PgRepository) getAssetChangeEventId(eventId int) (int, error) {
	selectSql := `select id from asset_change_events where event_id = $1;`
	return getId(r.db, selectSql, eventId)
}

// asset issuance events

func (r *PgRepository) GetOrCreateAssetIssuanceEvent(eventId int, assetId int, numberOfShares int64, unitOfMeasurement []byte, numberOfDecimalPlaces uint32) (int, error) {
	id, err := r.getAssetIssuanceEventId(eventId)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = r.insertAssetIssuanceEvent(eventId, assetId, numberOfShares, unitOfMeasurement, numberOfDecimalPlaces)
	}
	return id, errors.Wrap(err, "getting or creating asset issuance event")
}

func (r *PgRepository) insertAssetIssuanceEvent(eventId int, assetId int, numberOfShares int64, unitOfMeasurement []byte, numberOfDecimalPlaces uint32) (int, error) {
	insertSql := `insert into asset_issuance_events (event_id, asset_id, number_of_shares, unit_of_measurement, number_of_decimal_places) VALUES ($1, $2, $3, $4, $5) returning id;`
	return insert(r.db, insertSql, eventId, assetId, numberOfShares, unitOfMeasurement, numberOfDecimalPlaces)
}

func (r *PgRepository) getAssetIssuanceEventId(eventId int) (int, error) {
	selectSql := `select id from asset_issuance_events where event_id = $1;`
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

func (r *PgRepository) Close() {
	err := r.db.Close()
	if err != nil {
		slog.Error("error closing database.", "Error", err)
	} else {
		slog.Info("closed database.")
	}
}

func CreateDatabaseWithConfig(c *config.DatabaseConfig) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		c.Host, c.Port, c.User, c.Pass, c.Name)
	pgDb, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	pgDb.SetMaxOpenConns(c.MaxOpen)
	pgDb.SetMaxIdleConns(c.MaxIdle)
	slog.Info("Connected to database!")
	return pgDb, nil
}
