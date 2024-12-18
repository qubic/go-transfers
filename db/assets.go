package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func (r *PgRepository) GetOrCreateAsset(issuer, name string) (int, error) {
	id, err := r.getAssetId(issuer, name)
	if errors.Is(err, sql.ErrNoRows) { // not found create
		return r.createAsset(issuer, name)
	}
	return id, errors.Wrap(err, "getting or creating asset")
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
