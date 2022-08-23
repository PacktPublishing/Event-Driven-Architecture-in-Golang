package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/stores/internal/domain"
)

type ProductRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(tableName string, db *sql.DB) ProductRepository {
	return ProductRepository{tableName: tableName, db: db}
}

func (r ProductRepository) Find(ctx context.Context, id string) (*domain.Product, error) {
	const query = "SELECT store_id, name, description, sku, price FROM %s WHERE id = $1 LIMIT 1"

	product := &domain.Product{
		AggregateBase: ddd.AggregateBase{ID: id},
	}

	err := r.db.QueryRowContext(ctx, r.table(query), id).Scan(&product.StoreID, &product.Name, &product.Description, &product.SKU, &product.Price)
	if err != nil {
		return nil, errors.Wrap(err, "scanning product")
	}

	return product, nil
}

func (r ProductRepository) Save(ctx context.Context, product *domain.Product) error {
	const query = "INSERT INTO %s (id, store_id, name, description, sku, price) VALUES ($1, $2, $3, $4, $5, $6)"

	_, err := r.db.ExecContext(ctx, r.table(query), product.ID, product.StoreID, product.Name, product.Description, product.SKU, product.Price)

	return errors.Wrap(err, "inserting product")
}

func (r ProductRepository) Delete(ctx context.Context, id string) error {
	const query = "DELETE FROM %s WHERE id = $1 LIMIT 1"

	_, err := r.db.ExecContext(ctx, r.table(query), id)

	return errors.Wrap(err, "deleting product")
}

func (r ProductRepository) GetCatalog(ctx context.Context, storeID string) ([]*domain.Product, error) {
	const query = "SELECT id, name, description, sku, price FROM %s WHERE store_id = $1"

	products := make([]*domain.Product, 0)

	rows, err := r.db.QueryContext(ctx, r.table(query), storeID)
	if err != nil {
		return nil, errors.Wrap(err, "querying products")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing product rows")
		}
	}(rows)

	for rows.Next() {
		product := &domain.Product{
			StoreID: storeID,
		}
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.SKU, &product.Price)
		if err != nil {
			return nil, errors.Wrap(err, "scanning product")
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing product rows")
	}

	return products, nil
}

func (r ProductRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
