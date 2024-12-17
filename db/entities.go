package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// entity

func (r *PgRepository) GetOrCreateEntity(identity string) (int, error) {
	id, err := r.getEntityId(identity)
	if errors.Is(err, sql.ErrNoRows) { // insert if not found
		id, err = r.insertEntity(identity)
	}
	return id, errors.Wrap(err, "getting or creating entity")
}

func (r *PgRepository) insertEntity(identity string) (int, error) {
	insertSql := `insert into entities (identity) values ($1) returning id;`
	return insert(r.db, insertSql, identity)
}

func (r *PgRepository) getEntityId(identity string) (int, error) {
	selectSql := `select id from entities where identity= $1;`
	return getId(r.db, selectSql, identity)
}
