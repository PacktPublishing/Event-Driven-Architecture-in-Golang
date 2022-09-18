package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/stackus/errors"

	"eda-in-golang/notifications/internal/application"
	"eda-in-golang/notifications/internal/models"
)

type CustomerCacheRepository struct {
	tableName string
	db        *sql.DB
	fallback  application.CustomerRepository
}

var _ application.CustomerCacheRepository = (*CustomerCacheRepository)(nil)

func NewCustomerCacheRepository(tableName string, db *sql.DB, fallback application.CustomerRepository) CustomerCacheRepository {
	return CustomerCacheRepository{
		tableName: tableName,
		db:        db,
		fallback:  fallback,
	}
}

func (r CustomerCacheRepository) Add(ctx context.Context, customerID, name, smsNumber string) error {
	const query = "INSERT INTO %s (id, name, sms_number) VALUES ($1, $2, $3)"

	_, err := r.db.ExecContext(ctx, r.table(query), customerID, name, smsNumber)
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

func (r CustomerCacheRepository) UpdateSmsNumber(ctx context.Context, customerID, smsNumber string) error {
	const query = `UPDATE %s SET sms_number = $2 WHERE customerID = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), customerID, smsNumber)

	return err
}

func (r CustomerCacheRepository) Find(ctx context.Context, customerID string) (*models.Customer, error) {
	const query = `SELECT name, sms_number FROM %s WHERE id = $1 LIMIT 1`

	customer := &models.Customer{
		ID: customerID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), customerID).Scan(&customer.Name, &customer.SmsNumber)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.Wrap(err, "scanning customer")
		}
		customer, err = r.fallback.Find(ctx, customerID)
		if err != nil {
			return nil, errors.Wrap(err, "customer fallback failed")
		}
		// attempt to add it to the cache
		return customer, r.Add(ctx, customer.ID, customer.Name, customer.SmsNumber)
	}

	return customer, nil
}

func (r CustomerCacheRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
