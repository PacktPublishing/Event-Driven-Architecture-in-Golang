package postgres

import (
	"time"

	"eda-in-golang/internal/ddd"
	"eda-in-golang/internal/es"
)

type aggregateEvent struct {
	id         string
	name       string
	payload    ddd.EventPayload
	occurredAt time.Time
	aggregate  es.EventSourcedAggregate
	version    int
}

var _ ddd.AggregateEvent = (*aggregateEvent)(nil)

func (e aggregateEvent) ID() string                { return e.id }
func (e aggregateEvent) EventName() string         { return e.name }
func (e aggregateEvent) Payload() ddd.EventPayload { return e.payload }
func (e aggregateEvent) Metadata() ddd.Metadata    { return ddd.Metadata{} }
func (e aggregateEvent) OccurredAt() time.Time     { return e.occurredAt }
func (e aggregateEvent) AggregateName() string     { return e.aggregate.AggregateName() }
func (e aggregateEvent) AggregateID() string       { return e.aggregate.ID() }
func (e aggregateEvent) AggregateVersion() int     { return e.version }
