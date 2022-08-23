package postgres

import (
	"context"
	"fmt"

	"github.com/stackus/errors"

	"eda-in-golang/internal/postgres"
	"eda-in-golang/payments/internal/application"
	"eda-in-golang/payments/internal/models"
)

type InvoiceRepository struct {
	tableName string
	db        postgres.DB
}

var _ application.InvoiceRepository = (*InvoiceRepository)(nil)

func NewInvoiceRepository(tableName string, db postgres.DB) InvoiceRepository {
	return InvoiceRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r InvoiceRepository) Find(ctx context.Context, invoiceID string) (*models.Invoice, error) {
	const query = "SELECT order_id, amount, status FROM %s WHERE id = $1 LIMIT 1"

	invoice := &models.Invoice{
		ID: invoiceID,
	}
	var status string
	err := r.db.QueryRowContext(ctx, r.table(query), invoiceID).Scan(&invoice.OrderID, &invoice.Amount, &status)
	if err != nil {
		return nil, errors.Wrap(err, "scanning invoice")
	}

	invoice.Status, err = r.statusToDomain(status)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (r InvoiceRepository) Save(ctx context.Context, invoice *models.Invoice) error {
	const query = "INSERT INTO %s (id, order_id, amount, status) VALUES ($1, $2, $3, $4)"

	_, err := r.db.ExecContext(ctx, r.table(query), invoice.ID, invoice.OrderID, invoice.Amount, invoice.Status.String())

	return err
}

func (r InvoiceRepository) Update(ctx context.Context, invoice *models.Invoice) error {
	const query = "UPDATE %s SET amount = $2, status = $3 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), invoice.ID, invoice.Amount, invoice.Status.String())

	return err
}

func (r InvoiceRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

func (r InvoiceRepository) statusToDomain(status string) (models.InvoiceStatus, error) {
	switch status {
	case models.InvoiceIsPending.String():
		return models.InvoiceIsPending, nil
	case models.InvoiceIsPaid.String():
		return models.InvoiceIsPaid, nil
	default:
		return models.InvoiceIsUnknown, fmt.Errorf("unknown invoice status: %s", status)
	}
}
