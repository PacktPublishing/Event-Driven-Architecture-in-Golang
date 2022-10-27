package postgres

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stackus/errors"

	"eda-in-golang/internal/postgres"
	"eda-in-golang/search/internal/application"
	"eda-in-golang/search/internal/models"
)

type OrderRepository struct {
	tableName string
	db        postgres.DB
}

var _ application.OrderRepository = (*OrderRepository)(nil)

func NewOrderRepository(tableName string, db postgres.DB) OrderRepository {
	return OrderRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r OrderRepository) Add(ctx context.Context, order *models.Order) error {
	const query = `INSERT INTO %s (
order_id, customer_id, customer_name,
items, status, product_ids, store_ids,
created_at) VALUES (
$1, $2, $3,
$4, $5, $6, $7,
$8)`

	items, err := json.Marshal(order.Items)
	if err != nil {
		return err
	}

	productIDs := make(IDArray, len(order.Items))
	storeMap := make(map[string]struct{})
	for i, item := range order.Items {
		productIDs[i] = item.ProductID
		storeMap[item.StoreID] = struct{}{}
	}
	storeIDs := make(IDArray, 0, len(storeMap))
	for storeID, _ := range storeMap {
		storeIDs = append(storeIDs, storeID)
	}

	_, err = r.db.ExecContext(ctx, r.table(query),
		order.OrderID, order.CustomerID, order.CustomerName,
		items, order.Status, productIDs, storeIDs,
		order.CreatedAt,
	)
	return err
}

func (r OrderRepository) UpdateStatus(ctx context.Context, orderID, status string) error {
	const query = `UPDATE %s SET status = $2 WHERE order_id = $1`

	_, err := r.db.ExecContext(ctx, r.table(query), orderID, status)
	return err
}

func (r OrderRepository) Search(ctx context.Context, search application.SearchOrders) ([]*models.Order, error) {
	// TODO implement me
	panic("implement me")
}

func (r OrderRepository) Get(ctx context.Context, orderID string) (*models.Order, error) {
	const query = `SELECT customer_id, customer_name, items, status, created_at FROM %s WHERE order_id = $1`

	order := &models.Order{
		OrderID: orderID,
	}

	var itemData []byte
	err := r.db.QueryRowContext(ctx, r.table(query)).Scan(&order.CustomerID, &order.CustomerName, &itemData, &order.Status, &order.CreatedAt)
	if err != nil {
		return nil, err
	}

	var items []models.Item
	err = json.Unmarshal(itemData, &items)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (r OrderRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}

type IDArray []string

func (a *IDArray) Scan(src any) error {
	var sep = []byte(",")

	var data []byte
	switch v := src.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return errors.ErrInvalidArgument.Msgf("IDArray: unsupported type: %T", src)
	}

	ids := make([]string, bytes.Count(data, sep))
	for i, id := range bytes.Split(bytes.Trim(data, "{}"), sep) {
		ids[i] = string(id)
	}

	*a = ids

	return nil
}

func (a IDArray) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	if len(a) == 0 {
		return "{}", nil
	}
	// unsafe way to do this; assumption is all ids are UUIDs
	return fmt.Sprintf("{%s}", strings.Join(a, ",")), nil
}
