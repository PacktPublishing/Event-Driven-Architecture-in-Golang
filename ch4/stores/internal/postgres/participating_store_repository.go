package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/stackus/errors"

	"eda-in-golang/stores/internal/domain"
)

type ParticipatingStoreRepository struct {
	tableName string
	db        *sql.DB
}

var _ domain.ParticipatingStoreRepository = (*ParticipatingStoreRepository)(nil)

func NewParticipatingStoreRepository(tableName string, db *sql.DB) ParticipatingStoreRepository {
	return ParticipatingStoreRepository{tableName: tableName, db: db}
}

func (r ParticipatingStoreRepository) FindAll(ctx context.Context) (stores []*domain.Store, err error) {
	const query = "SELECT id, name, location, participating FROM %s WHERE participating is true"

	rows, err := r.db.QueryContext(ctx, r.table(query))
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
		store := &domain.Store{}
		err := rows.Scan(&store.ID, &store.Name, &store.Location, &store.Participating)
		if err != nil {
			return nil, errors.Wrap(err, "scanning participating store")
		}

		stores = append(stores, store)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "finishing participating store rows")
	}

	return stores, nil
}

func (r ParticipatingStoreRepository) table(query string) string {
	return fmt.Sprintf(query, r.tableName)
}
