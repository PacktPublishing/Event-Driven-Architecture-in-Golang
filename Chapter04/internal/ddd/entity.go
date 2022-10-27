package ddd

type Entity interface {
	GetID() string
}

type EntityBase struct {
	ID string
}

func (e EntityBase) GetID() string {
	return e.ID
}
