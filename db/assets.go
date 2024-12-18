package db

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func (r *PgRepository) GetOrCreateAsset(ctx context.Context, issuer, name string) (int, error) {
	id, err := r.getAssetId(ctx, issuer, name)
	if errors.Is(err, sql.ErrNoRows) { // not found create
		return r.createAsset(ctx, issuer, name)
	}
	return id, errors.Wrap(err, "getting or creating asset")
}

func (r *PgRepository) createAsset(ctx context.Context, issuer, name string) (int, error) {
	entityId, err := r.GetOrCreateEntity(ctx, issuer)
	if err != nil {
		return 0, errors.Wrap(err, "creating asset")
	}
	return r.insertAsset(ctx, entityId, name)
}

func (r *PgRepository) getAssetId(ctx context.Context, issuer, name string) (int, error) {
	selectSql := `select a.id from assets a
    join entities e on a.issuer_id = e.id
    where e.identity=$1 and a.name=$2;`

	return getId(ctx, r.db, selectSql, issuer, name)
}

func (r *PgRepository) insertAsset(ctx context.Context, issuerId int, name string) (int, error) {
	insertSql := `insert into assets (issuer_id, name) values ($1, $2) returning id;`
	return insert(ctx, r.db, insertSql, issuerId, name)
}
