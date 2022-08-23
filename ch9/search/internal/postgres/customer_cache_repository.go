package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"

	"eda-in-golang/internal/postgres"
	"eda-in-golang/search/internal/application"
	"eda-in-golang/search/internal/models"
)

type CustomerCacheRepository struct {
	tableName string
	db        postgres.DB
	fallback  application.CustomerRepository
}

var _ application.CustomerCacheRepository = (*CustomerCacheRepository)(nil)

func NewCustomerCacheRepository(tableName string, db postgres.DB, fallback application.CustomerRepository) CustomerCacheRepository {
	return CustomerCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r CustomerCacheRepository) Add(ctx context.Context, customerID, name string) error {
	const query = "INSERT INTO %s (id, name) VALUES ($1, $2)"

	_, err := r.db.ExecContext(ctx, r.table(query), customerID, name)
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

func (r CustomerCacheRepository) Find(ctx context.Context, customerID string) (*models.Customer, error) {
	const query = `SELECT name FROM %s WHERE id = $1 LIMIT 1`

	customer := &models.Customer{
		ID: customerID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), customerID).Scan(&customer.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning customer")
		}
		customer, err = r.fallback.Find(ctx, customerID)
		if err != nil {
			return nil, errors.Wrap(err, "customer fallback failed")
		}
		// attempt to add it to the cache
		return customer, r.Add(ctx, customer.ID, customer.Name)
	}

	return customer, nil
}

func (r CustomerCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
