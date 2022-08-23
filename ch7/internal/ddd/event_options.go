package ddd

type EventOption interface {
	configureEvent(*event)
}
