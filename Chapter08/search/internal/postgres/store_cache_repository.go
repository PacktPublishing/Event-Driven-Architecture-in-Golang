package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"

	"eda-in-golang/search/internal/application"
	"eda-in-golang/search/internal/models"
)

type StoreCacheRepository struct {
	tableName string
	db        *sql.DB
	fallback  application.StoreRepository
}

var _ application.StoreCacheRepository = (*StoreCacheRepository)(nil)

func NewStoreCacheRepository(tableName string, db *sql.DB, fallback application.StoreRepository) StoreCacheRepository {
	return StoreCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r StoreCacheRepository) Add(ctx context.Context, storeID, name string) error {
	const query = "INSERT INTO %s (id, name) VALUES ($1, $2)"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return nil
			}
		}
	}

	return err
}

func (r StoreCacheRepository) Rename(ctx context.Context, storeID, name string) error {
	const query = "UPDATE %s SET name = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name)

	return err
}

func (r StoreCacheRepository) Find(ctx context.Context, storeID string) (*models.Store, error) {
	const query = "SELECT name FROM %s WHERE id = $1 LIMIT 1"

	store := &models.Store{
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
