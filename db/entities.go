package db

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

func (r *PgRepository) GetOrCreateEntity(ctx context.Context, identity string) (int, error) {
	id, err := r.getEntityId(ctx, identity)
	if errors.Is(err, sql.ErrNoRows) { // insert if not found
		id, err = r.insertEntity(ctx, identity)
	}
	return id, errors.Wrap(err, "getting or creating entity")
}

func (r *PgRepository) insertEntity(ctx context.Context, identity string) (int, error) {
	insertSql := `insert into entities (identity) values ($1) returning id;`
	return insert(ctx, r.db, insertSql, identity)
}

func (r *PgRepository) getEntityId(ctx context.Context, identity string) (int, error) {
	selectSql := `select id from entities where identity= $1;`
	return getId(ctx, r.db, selectSql, identity)
}
