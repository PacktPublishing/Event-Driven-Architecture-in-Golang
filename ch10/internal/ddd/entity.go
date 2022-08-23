package ddd

type (
	IDer interface {
		ID() string
	}

	EntityNamer interface {
		EntityName() string
	}

	Entity interface {
		IDer
		EntityNamer
		IDSetter
		NameSetter
	}

	entity struct {
		id   string
		name string
	}
)

var _ Entity = (*entity)(nil)

func NewEntity(id, name string) *entity {
	return &entity{
		id:   id,
		name: name,
	}
}

func (e entity) ID() string             { return e.id }
func (e entity) EntityName() string     { return e.name }
func (e entity) Equals(other IDer) bool { return e.id == other.ID() }

func (e *entity) SetID(id string)     { e.id = id }
func (e *entity) SetName(name string) { e.name = name }
