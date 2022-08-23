package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stackus/errors"

	"eda-in-golang/internal/postgres"
	"eda-in-golang/stores/internal/domain"
)

type MallRepository struct {
	tableName string
	db        postgres.DB
}

var _ domain.MallRepository = (*MallRepository)(nil)

func NewMallRepository(tableName string, db postgres.DB) MallRepository {
	return MallRepository{
		tableName: tableName,
		db:        db,
	}
}

func (r MallRepository) AddStore(ctx context.Context, storeID, name, location string) error {
	const query = "INSERT INTO %s (id, name, location, participating) VALUES ($1, $2, $3, $4)"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name, location, false)

	return err
}

func (r MallRepository) SetStoreParticipation(ctx context.Context, storeID string, participating bool) error {
	const query = "UPDATE %s SET participating = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, participating)

	return err
}

func (r MallRepository) RenameStore(ctx context.Context, storeID, name string) error {
	const query = "UPDATE %s SET name = $2 WHERE id = $1"

	_, err := r.db.ExecContext(ctx, r.table(query), storeID, name)

	return err
}

func (r MallRepository) Find(ctx context.Context, storeID string) (*domain.MallStore, error) {
	const query = "SELECT name, location, participating FROM %s WHERE id = $1 LIMIT 1"

	store := &domain.MallStore{
		ID: storeID,
	}

	err := r.db.QueryRowContext(ctx, r.table(query), storeID).Scan(&store.Name, &store.Location, &store.Participating)
	if err != nil {
		return nil, errors.Wrap(err, "scanning store")
	}

	return store, nil
}

func (r MallRepository) All(ctx context.Context) (stores []*domain.MallStore, err error) {
	const query = "SELECT id, name, location, participating FROM %s"

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, errors.Wrap(err, "querying stores")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing store rows")
		}
	}(rows)

	for rows.Next() {
		store := new(domain.MallStore)
		err := rows.Scan(&store.ID, &store.Name, &store.Location, &store.Participating)
		if err != nil {
			return nil, errors.Wrap(err, "scanning store")
		}

		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing store rows")
	}

	return stores, nil
}

func (r MallRepository) AllParticipating(ctx context.Context) (stores []*domain.MallStore, err error) {
	const query = "SELECT id, name, location, participating FROM %s WHERE participating is true"

	var rows *sql.Rows
	rows, err = r.db.QueryContext(ctx, r.table(query))
	if err != nil {
		return nil, errors.Wrap(err, "querying participating stores")
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing participating store rows")
		}
	}(rows)

	for rows.Next() {
		store := new(domain.MallStore)
		err := rows.Scan(&store.ID, &store.Name, &store.Location, &store.Participating)
		if err != nil {
			return nil, errors.Wrap(err, "scanning store")
		}

		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing participating store rows")
	}

	return stores, nil
}

func (r MallRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
