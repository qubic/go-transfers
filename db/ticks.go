package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

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
