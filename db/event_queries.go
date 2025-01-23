package db

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"go-transfers/proto"
)

// qu transfer events

func (r *PgRepository) GetQuTransferEventsForTick(ctx context.Context, tickNumber int) ([]*proto.QuTransferEvent, error) {
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
		order by e.event_id;`
	var events []*proto.QuTransferEvent
	err := r.db.SelectContext(ctx, &events, selectSql, tickNumber)
	if err != nil {
		return nil, errors.Wrap(err, "getting asset change events")
	}
	return events, nil
}

func (r *PgRepository) GetQuTransferEventsForEntity(ctx context.Context, identity string) ([]*proto.QuTransferEvent, error) {
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
	err := r.db.SelectContext(ctx, &events, selectSql, identity)
	if err != nil {
		return nil, errors.Wrap(err, "getting asset change events")
	}
	return events, nil
}

// asset change events

func (r *PgRepository) GetAssetChangeEventsForTick(ctx context.Context, tickNumber int) ([]*proto.AssetChangeEvent, error) {
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
		order by e.event_id;`
	var events []*proto.AssetChangeEvent
	err := r.db.SelectContext(ctx, &events, selectSql, tickNumber)
	if err != nil {
		return nil, errors.Wrap(err, "getting asset change events")
	}
	return events, nil
}

func (r *PgRepository) GetAssetIssuanceEventsForTick(ctx context.Context, tickNumber int) ([]*proto.AssetIssuanceEvent, error) {
	selectSql := `select issuer.identity issuerId,
       		a.name, 
       		aie.number_of_shares numberOfShares,
       		aie.unit_of_measurement unitOfMeasurement,
       		aie.number_of_decimal_places numberOfDecimalPlaces,
       		tx.hash transactionHash,
       		ti.tick_number tick,
       		e.event_type eventType
		from asset_issuance_events aie
		join events e on aie.event_id = e.id
		join transactions tx on e.transaction_id = tx.id
		join ticks ti on tx.tick_id = ti.id
		join assets a on aie.asset_id = a.id
		join entities issuer on a.issuer_id = issuer.id
		where ti.tick_number = $1
		and e.event_type = 1
		order by e.event_id;`
	var events []*proto.AssetIssuanceEvent
	err := r.db.SelectContext(ctx, &events, selectSql, tickNumber)
	if err != nil {
		return nil, errors.Wrap(err, "getting asset issuance events")
	}
	return events, nil
}

func (r *PgRepository) GetAssetChangeEventsForEntity(ctx context.Context, identity string) ([]*proto.AssetChangeEvent, error) {
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
	err := r.db.SelectContext(ctx, &events, selectSql, identity)
	if err != nil {
		return nil, errors.Wrap(err, "getting asset change events")
	}
	return events, nil
}
