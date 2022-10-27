package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/stackus/errors"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/es"
	"eda-in-golang/internal/registry"
)

type (
	EventStore struct {
		tableName string
		db        DB
		registry  registry.Registry
	}

	aggregateEvent struct {
		id         string
		name       string
		payload    ddd.EventPayload
		occurredAt time.Time
		aggregate  es.EventSourcedAggregate
		version    int
	}
)

var _ es.AggregateStore = (*EventStore)(nil)

var _ ddd.AggregateEvent = (*aggregateEvent)(nil)

func NewEventStore(tableName string, db DB, registry registry.Registry) EventStore {
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
	const query = `INSERT INTO %s (stream_id, stream_name, stream_version, event_id, event_name, event_data, occurred_at) VALUES`

	aggregateID := aggregate.ID()
	aggregateName := aggregate.AggregateName()

	placeholders := make([]string, len(aggregate.Events()))
	values := make([]any, len(aggregate.Events())*7)

	for i, event := range aggregate.Events() {
		var payloadData []byte

		payloadData, err = s.registry.Serialize(event.EventName(), event.Payload())
		if err != nil {
			return err
		}

		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*7+1, i*7+2, i*7+3, i*7+4, i*7+5, i*7+6, i*7+7,
		)

		values[i*7] = aggregateID
		values[i*1+1] = aggregateName
		values[i*7+2] = event.AggregateVersion()
		values[i*7+3] = event.ID()
		values[i*7+4] = event.EventName()
		values[i*7+5] = payloadData
		values[i*7+6] = event.OccurredAt()
	}
	if _, err = s.db.ExecContext(
		ctx,
		fmt.Sprintf("%s %s", s.table(query), strings.Join(placeholders, ",")),
		values...,
	); err != nil {
		return err
	}

	return nil
}

func (s EventStore) table(query string) string {
	return fmt.Sprintf(query, s.tableName)
}

func (e aggregateEvent) ID() string                { return e.id }
func (e aggregateEvent) EventName() string         { return e.name }
func (e aggregateEvent) Payload() ddd.EventPayload { return e.payload }
func (e aggregateEvent) Metadata() ddd.Metadata    { return ddd.Metadata{} }
func (e aggregateEvent) OccurredAt() time.Time     { return e.occurredAt }
func (e aggregateEvent) AggregateName() string     { return e.aggregate.AggregateName() }
func (e aggregateEvent) AggregateID() string       { return e.aggregate.ID() }
func (e aggregateEvent) AggregateVersion() int     { return e.version }
