package es

import (
	"eda-in-golang/internal/ddd"
)

type (
	Versioner interface {
		Version() int
		PendingVersion() int
	}

	Aggregate struct {
		ddd.Aggregate
		version int
	}
)

var _ interface {
	EventCommitter
	Versioner
	VersionSetter
} = (*Aggregate)(nil)

func NewAggregate(id, name string) Aggregate {
	return Aggregate{
		Aggregate: ddd.NewAggregate(id, name),
		version:   0,
	}
}

func (a *Aggregate) AddEvent(name string, payload ddd.EventPayload, options ...ddd.EventOption) {
	options = append(
		options,
		ddd.Metadata{
			ddd.AggregateVersionKey: a.PendingVersion() + 1,
		},
	)
	a.Aggregate.AddEvent(name, payload, options...)
}

func (a *Aggregate) CommitEvents() {
	a.version += len(a.Events())
	a.ClearEvents()
}

func (a Aggregate) Version() int        { return a.version }
func (a Aggregate) PendingVersion() int { return a.version + len(a.Events()) }

func (a *Aggregate) setVersion(version int) { a.version = version }
