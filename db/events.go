package db

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// events

func (r *PgRepository) GetOrCreateEvent(ctx context.Context, transactionId int, eventEventId uint64, eventType uint32, eventData string) (int, error) {
	id, err := r.getEventId(ctx, transactionId, eventEventId)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = r.insertEvent(ctx, transactionId, eventEventId, eventType, eventData)
	}
	return id, errors.Wrapf(err, "getting or creating event for transaction id [%d] and events event id [%d]", transactionId, eventEventId)
}

func (r *PgRepository) insertEvent(ctx context.Context, transactionId int, eventEventId uint64, eventType uint32, eventData string) (int, error) {
	insertSql := `insert into events (transaction_id, event_id, event_type, event_data) values ($1, $2, $3, $4) returning id;`
	return insert(ctx, r.db, insertSql, transactionId, eventEventId, eventType, eventData)
}

func (r *PgRepository) getEventId(ctx context.Context, transactionId int, eventEventId uint64) (int, error) {
	selectSql := `select id from events where transaction_id = $1 and event_id = $2;`
	return getId(ctx, r.db, selectSql, transactionId, eventEventId)
}

// qu transfer events

func (r *PgRepository) GetOrCreateQuTransferEvent(ctx context.Context, eventId int, sourceEntityId int, destinationEntityId int, amount uint64) (int, error) {
	id, err := r.getQuTransferEventId(ctx, eventId)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = r.insertQuTransferEvent(ctx, eventId, sourceEntityId, destinationEntityId, amount)
	}
	return id, errors.Wrapf(err, "getting or creating qu transfer for event [%d]", eventId)
}

func (r *PgRepository) insertQuTransferEvent(ctx context.Context, eventId int, sourceEntityId int, destinationEntityId int, amount uint64) (int, error) {
	insertSql := `insert into qu_transfer_events (event_id, source_entity_id, destination_entity_id, amount) values ($1, $2, $3, $4) returning id;`
	return insert(ctx, r.db, insertSql, eventId, sourceEntityId, destinationEntityId, amount)
}

func (r *PgRepository) getQuTransferEventId(ctx context.Context, eventId int) (int, error) {
	selectSql := `select id from qu_transfer_events where event_id = $1;`
	return getId(ctx, r.db, selectSql, eventId)
}

// asset change events

func (r *PgRepository) GetOrCreateAssetChangeEvent(ctx context.Context, eventId, assetId, sourceEntityId, destinationEntityId int, numberOfShares int64) (int, error) {
	id, err := r.getAssetChangeEventId(ctx, eventId)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = r.insertAssetChangeEvent(ctx, eventId, assetId, sourceEntityId, destinationEntityId, numberOfShares)
	}
	return id, errors.Wrapf(err, "getting or creating asset change for event [%d]", eventId)
}

func (r *PgRepository) insertAssetChangeEvent(ctx context.Context, eventId, assetId, sourceEntityId, destinationEntityId int, numberOfShares int64) (int, error) {
	insertSql := `insert into asset_change_events (event_id, asset_id, source_entity_id, destination_entity_id, number_of_shares) values ($1, $2, $3, $4, $5) returning id;`
	return insert(ctx, r.db, insertSql, eventId, assetId, sourceEntityId, destinationEntityId, numberOfShares)
}

func (r *PgRepository) getAssetChangeEventId(ctx context.Context, eventId int) (int, error) {
	selectSql := `select id from asset_change_events where event_id = $1;`
	return getId(ctx, r.db, selectSql, eventId)
}

// asset issuance events

func (r *PgRepository) GetOrCreateAssetIssuanceEvent(ctx context.Context, eventId int, assetId int, numberOfShares int64, unitOfMeasurement string, numberOfDecimalPlaces uint32) (int, error) {
	id, err := r.getAssetIssuanceEventId(ctx, eventId)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = r.insertAssetIssuanceEvent(ctx, eventId, assetId, numberOfShares, unitOfMeasurement, numberOfDecimalPlaces)
	}
	return id, errors.Wrapf(err, "getting or creating asset issuance for event [%d]", eventId)
}

func (r *PgRepository) insertAssetIssuanceEvent(ctx context.Context, eventId int, assetId int, numberOfShares int64, unitOfMeasurement string, numberOfDecimalPlaces uint32) (int, error) {
	insertSql := `insert into asset_issuance_events (event_id, asset_id, number_of_shares, unit_of_measurement, number_of_decimal_places) VALUES ($1, $2, $3, $4, $5) returning id;`
	return insert(ctx, r.db, insertSql, eventId, assetId, numberOfShares, unitOfMeasurement, numberOfDecimalPlaces)
}

func (r *PgRepository) getAssetIssuanceEventId(ctx context.Context, eventId int) (int, error) {
	selectSql := `select id from asset_issuance_events where event_id = $1;`
	return getId(ctx, r.db, selectSql, eventId)
}
