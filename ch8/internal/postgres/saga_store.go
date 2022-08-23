package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"eda-in-golang/internal/registry"
	"eda-in-golang/internal/sec"
)

type SagaStore struct {
	tableName string
	db        *sql.DB
	registry  registry.Registry
}

var _ sec.SagaStore = (*SagaStore)(nil)

func NewSagaStore(tableName string, db *sql.DB, registry registry.Registry) SagaStore {
	return SagaStore{
		tableName: tableName,
		db:        db,
		registry:  registry,
	}
}

func (s SagaStore) Load(ctx context.Context, sagaName, sagaID string) (*sec.SagaContext[[]byte], error) {
	const query = "SELECT data, step, done, compensating FROM %s WHERE name = $1 AND id = $2"

	sagaCtx := &sec.SagaContext[[]byte]{
		ID: sagaID,
	}
	err := s.db.QueryRowContext(ctx, s.table(query), sagaName, sagaID).Scan(&sagaCtx.Data, &sagaCtx.Step, &sagaCtx.Done, &sagaCtx.Compensating)

	return sagaCtx, err
}

func (s SagaStore) Save(ctx context.Context, sagaName string, sagaCtx *sec.SagaContext[[]byte]) error {
	const query = `INSERT INTO %s (name, id, data, step, done, compensating) 
VALUES ($1, $2, $3, $4, $5, $6) 
ON CONFLICT (name, id) DO
UPDATE SET data = EXCLUDED.data, step = EXCLUDED.step, done = EXCLUDED.done, compensating = EXCLUDED.compensating`

	_, err := s.db.ExecContext(ctx, s.table(query), sagaName, sagaCtx.ID, sagaCtx.Data, sagaCtx.Step, sagaCtx.Done, sagaCtx.Compensating)

	return err
}

func (s SagaStore) table(query string) string {
	return fmt.Sprintf(query, s.tableName)
}
