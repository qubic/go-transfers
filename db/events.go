package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

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
