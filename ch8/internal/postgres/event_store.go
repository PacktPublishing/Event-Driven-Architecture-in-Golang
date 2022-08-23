package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/stackus/errors"

	"eda-in-golang/internal/es"
	"eda-in-golang/internal/registry"
)

type EventStore struct {
	tableName string
	db        *sql.DB
	registry  registry.Registry
}

var _ es.AggregateStore = (*EventStore)(nil)

func NewEventStore(tableName string, db *sql.DB, registry registry.Registry) EventStore {
	return EventStore{
		tableName: tableName,
		db:        db,
		registry:  registry,
	}
}

func (s EventStore) Load(ctx context.Context, aggregate es.EventSourcedAggregate) (err error) {
	const query = `SELECT stream_version, event_id, event_name, event_data, occurred_at FROM %s WHERE stream_id = $1 AND stream_name = $2 AND stream_version > $3 ORDER BY stream_version ASC`

	aggregateID := aggregate.ID()
	aggregateName := aggregate.AggregateName()

	var rows *sql.Rows

	rows, err = s.db.QueryContext(ctx, s.table(query), aggregateID, aggregateName, aggregate.Version())
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			err = errors.Wrap(err, "closing event rows")
		}
	}(rows)

	for rows.Next() {
		var eventID, eventName string
		var payloadData []byte
		var aggregateVersion int
		var occurredAt time.Time
		err := rows.Scan(&aggregateVersion, &eventID, &eventName, &payloadData, &occurredAt)
		if err != nil {
			return err
		}

		var payload interface{}
		payload, err = s.registry.Deserialize(eventName, payloadData)
		if err != nil {
			return err
		}

		event := aggregateEvent{
			id:         eventID,
			name:       eventName,
			payload:    payload,
			aggregate:  aggregate,
			version:    aggregateVersion,
			occurredAt: occurredAt,
		}

		if err = es.LoadEvent(aggregate, event); err != nil {
			return err
		}
	}

	return nil
}

func (s EventStore) Save(ctx context.Context, aggregate es.EventSourcedAggregate) (err error) {
	const query = `INSERT INTO %s (stream_id, stream_name, stream_version, event_id, event_name, event_data, occurred_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	var tx *sql.Tx
	tx, err = s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		p := recover()
		switch {
		case p != nil:
			_ = tx.Rollback()
			panic(p)
		case err != nil:
			rErr := tx.Rollback()
			if rErr != nil {
				err = errors.Wrap(err, rErr.Error())
			}
		default:
			err = tx.Commit()
		}
	}()

	aggregateID := aggregate.ID()
	aggregateName := aggregate.AggregateName()

	for _, event := range aggregate.Events() {
		var payloadData []byte

		payloadData, err = s.registry.Serialize(event.EventName(), event.Payload())
		if err != nil {
			return err
		}
		if _, err = tx.ExecContext(
			ctx, s.table(query),
			aggregateID, aggregateName, event.AggregateVersion(), event.ID(), event.EventName(), payloadData, event.OccurredAt(),
		); err != nil {
			return err
		}
	}

	return nil
}

func (s EventStore) table(query string) string {
	return fmt.Sprintf(query, s.tableName)
}
