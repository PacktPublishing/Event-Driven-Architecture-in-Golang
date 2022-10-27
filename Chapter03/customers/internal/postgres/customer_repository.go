package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"eda-in-golang/customers/internal/domain"
)

type CustomerRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.CustomerRepository = (*CustomerRepository)(nil)

func NewCustomerRepository(tableName string, db *sql.DB) CustomerRepository {
	return CustomerRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r CustomerRepository) Save(ctx context.Context, customer *domain.Customer) error {
	const query = "INSERT INTO %s (id, name, sms_number, enabled) VALUES ($1, $2, $3, $4)"

	_, err := r.db.ExecContext(ctx, r.table(query), customer.ID, customer.Name, customer.SmsNumber, customer.Enabled)

	return err
}

func (r CustomerRepository) Find(ctx context.Context, customerID string) (*domain.Customer, error) {
	const query = "SELECT name, sms_number, enabled FROM %s WHERE id = $1 LIMIT 1"

	customer := &domain.Customer{
		ID: customerID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), customerID).Scan(&customer.Name, &customer.SmsNumber, &customer.Enabled)

	return customer, err
}

func (r CustomerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	const query = "UPDATE %s SET name = $2, sms_number = $3, enabled = $4 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), customer.ID, customer.Name, customer.SmsNumber, customer.Enabled)

	return err
}

func (r CustomerRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
