package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/stackus/errors"

	"eda-in-golang/ordering/internal/domain"
)

type OrderRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(tableName string, db *sql.DB) OrderRepository {
	return OrderRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r OrderRepository) Find(ctx context.Context, orderID string) (*domain.Order, error) {
	const query = "SELECT customer_id, payment_id, shopping_id, invoice_id, items, status FROM %s WHERE id = $1 LIMIT 1"

	order := &domain.Order{
		ID: orderID,
	}

	var items []byte
	var status string

	err := r.db.QueryRowContext(ctx, r.table(query), orderID).Scan(&order.CustomerID, &order.PaymentID, &order.ShoppingID, &order.InvoiceID, &items, &status)
	if err != nil {
		return nil, errors.Wrap(err, "scanning order")
	}

	order.Status = domain.ToOrderStatus(status)

	err = json.Unmarshal(items, &order.Items)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshalling items")
	}

	return order, nil
}

func (r OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	const query = "INSERT INTO %s (id, customer_id, payment_id, shopping_id, invoice_id, items, status) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	items, err := json.Marshal(order.Items)
	if err != nil {
		return errors.Wrap(err, "marshalling items")
	}

	_, err = r.db.ExecContext(ctx, r.table(query), order.ID, order.CustomerID, order.PaymentID, order.ShoppingID, order.InvoiceID, items, order.Status.String())
	if err != nil {
		return errors.Wrap(err, "inserting order")
	}

	return nil
}

func (r OrderRepository) Update(ctx context.Context, order *domain.Order) error {
	const query = "UPDATE %s SET customer_id = $2, payment_id = $3, shopping_id = $4, invoice_id = $5, items = $6, status = $7 WHERE id = $1"

	items, err := json.Marshal(order.Items)
	if err != nil {
		return errors.Wrap(err, "marshalling items")
	}

	_, err = r.db.ExecContext(ctx, r.table(query), order.ID, order.CustomerID, order.PaymentID, order.ShoppingID, order.InvoiceID, items, order.Status.String())
	if err != nil {
		return errors.Wrap(err, "updating order")
	}

	return nil
}

func (r OrderRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
