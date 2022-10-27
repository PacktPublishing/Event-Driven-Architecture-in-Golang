package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"

	"eda-in-golang/depot/internal/domain"
	"eda-in-golang/internal/postgres"
)

type ProductCacheRepository struct {
	tableName string
	db        postgres.DB
	fallback  domain.ProductRepository
}

var _ domain.ProductCacheRepository = (*ProductCacheRepository)(nil)

func NewProductCacheRepository(tableName string, db postgres.DB, fallback domain.ProductRepository) ProductCacheRepository {
	return ProductCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r ProductCacheRepository) Add(ctx context.Context, productID, storeID, name string) error {
	const query = `INSERT INTO %s (id, store_id, name) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, r.table(query), productID, storeID, name)
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

func (r ProductCacheRepository) Rebrand(ctx context.Context, productID, name string) error {
	const query = `UPDATE %s SET name = $2 WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), productID, name)

	return err
}

func (r ProductCacheRepository) Remove(ctx context.Context, productID string) error {
	const query = `DELETE FROM %s WHERE id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), productID)

	return err
}

func (r ProductCacheRepository) Find(ctx context.Context, productID string) (*domain.Product, error) {
	const query = `SELECT store_id, name FROM %s WHERE id = $1 LIMIT 1`

	product := &domain.Product{
		ID: productID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), productID).Scan(&product.StoreID, &product.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning product")
		}
		product, err = r.fallback.Find(ctx, productID)
		if err != nil {
			return nil, errors.Wrap(err, "product fallback failed")
		}
		// attempt to add it to the cache
		return product, r.Add(ctx, product.ID, product.StoreID, product.Name)
	}

	return product, nil
}

func (r ProductCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
