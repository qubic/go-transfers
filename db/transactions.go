package db

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func (r *PgRepository) GetOrCreateTransaction(ctx context.Context, hash string, tickId int) (int, error) {
	id, err := r.getTransactionId(ctx, hash, tickId)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = r.insertTransaction(ctx, hash, tickId)
	}
	return id, errors.Wrapf(err, "getting or creating transaction [%s]", hash)
}

func (r *PgRepository) getTransactionId(ctx context.Context, hash string, tickId int) (int, error) {
	selectSql := `select id from transactions where hash = $1 and tick_id = $2;`
	return getId(ctx, r.db, selectSql, hash, tickId)
}

func (r *PgRepository) insertTransaction(ctx context.Context, hash string, tickId int) (int, error) {
	insertSql := `insert into transactions (hash, tick_id) values ($1, $2) returning id;`
	return insert(ctx, r.db, insertSql, hash, tickId)
}
