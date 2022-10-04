package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stackus/errors"

	"eda-in-golang/baskets/internal/domain"
	"eda-in-golang/internal/postgres"
)

type StoreCacheRepository struct {
	tableName string
	db        postgres.DB
	fallback  domain.StoreRepository
}

var _ domain.StoreCacheRepository = (*StoreCacheRepository)(nil)

func NewStoreCacheRepository(tableName string, db postgres.DB, fallback domain.StoreRepository) StoreCacheRepository {
	return StoreCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r StoreCacheRepository) Add(ctx context.Context, storeID, name string) error {
	const query = "INSERT INTO %s (id, NAME) VALUES ($1, $2) ON CONFLICT DO NOTHING"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name)

	return err
}

func (r StoreCacheRepository) Rename(ctx context.Context, storeID, name string) error {
	const query = "UPDATE %s SET NAME = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name)

	return err
}

func (r StoreCacheRepository) Find(ctx context.Context, storeID string) (*domain.Store, error) {
	const query = "SELECT name FROM %s WHERE id = $1 LIMIT 1"

	store := &domain.Store{
		ID: storeID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), storeID).Scan(&store.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning store")
		}
		store, err = r.fallback.Find(ctx, storeID)
		if err != nil {
			return nil, errors.Wrap(err, "store fallback failed")
		}
		// attempt to add it to the cache
		return store, r.Add(ctx, store.ID, store.Name)
	}

	return store, nil
}

func (r StoreCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
