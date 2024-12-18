package db

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func (r *PgRepository) GetOrCreateTick(ctx context.Context, tickNumber uint32) (int, error) {
	id, err := r.getTickId(ctx, tickNumber)
	if errors.Is(err, sql.ErrNoRows) {
		id, err = r.insertTick(ctx, tickNumber)
	}
	return id, errors.Wrap(err, "getting or creating tick")
}

func (r *PgRepository) getTickId(ctx context.Context, tickNumber uint32) (int, error) {
	selectSql := `select id from ticks where tick_number = $1;`
	return getId(ctx, r.db, selectSql, tickNumber)
}

func (r *PgRepository) insertTick(ctx context.Context, tickNumber uint32) (int, error) {
	insertSql := `insert into ticks (tick_number) values ($1) returning id;`
	return insert(ctx, r.db, insertSql, tickNumber)
}
