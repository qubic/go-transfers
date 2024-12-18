package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

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
