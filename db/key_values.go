package db

import (
	"context"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

// key value

func (r *PgRepository) GetLatestTick(ctx context.Context) (int, error) {
	return r.getNumericValue(ctx, "tick")
}

func (r *PgRepository) UpdateLatestTick(ctx context.Context, tickNumber int) error {
	return r.updateNumericValue(ctx, "tick", tickNumber)
}

func (r *PgRepository) getNumericValue(ctx context.Context, key string) (int, error) {
	selectSql := `select numeric_value from key_values where key = $1`
	var value int
	err := r.db.GetContext(ctx, &value, selectSql, key)
	return value, errors.Wrap(err, "getting numeric value")
}

func (r *PgRepository) updateNumericValue(ctx context.Context, key string, value int) error {
	updateSql := `update key_values set numeric_value = $1 where key = $2`
	_, err := r.db.ExecContext(ctx, updateSql, value, key)
	return errors.Wrap(err, "updating numeric value")
}
