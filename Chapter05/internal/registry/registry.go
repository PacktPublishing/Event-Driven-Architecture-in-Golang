package registry

import (
	"fmt"
	"sync"
)

type (
	Registrable interface {
		Key() string
	}

	Serializer   func(v interface{}) ([]byte, error)
	Deserializer func(d []byte, v interface{}) error

	Registry interface {
		Serialize(key string, v interface{}) ([]byte, error)
		Build(key string, options ...BuildOption) (interface{}, error)
		Deserialize(key string, data []byte, options ...BuildOption) (interface{}, error)
		register(key string, fn func() interface{}, s Serializer, d Deserializer, o []BuildOption) error
	}
)

type registered struct {
	factory      func() interface{}
	serializer   Serializer
	deserializer Deserializer
	options      []BuildOption
}

type registry struct {
	registered map[string]registered
	mu         sync.RWMutex
}

func New() *registry {
	return &registry{
		registered: make(map[string]registered),
	}
}

func (r *registry) Serialize(key string, v interface{}) ([]byte, error) {
	reg, exists := r.registered[key]
	if !exists {
		return nil, fmt.Errorf("nothing has been registered with the key `%s`", key)
	}
	return reg.serializer(v)
}

func (r *registry) Deserialize(key string, data []byte, options ...BuildOption) (interface{}, error) {
	v, err := r.Build(key, options...)
	if err != nil {
		return nil, err
	}

	err = r.registered[key].deserializer(data, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (r *registry) Build(key string, options ...BuildOption) (interface{}, error) {
	reg, exists := r.registered[key]
	if !exists {
		return nil, fmt.Errorf("nothing has been registered with the key `%s`", key)
	}

	v := reg.factory()
	uos := append(r.registered[key].options, options...)

	for _, option := range uos {
		err := option(v)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (r *registry) register(key string, fn func() interface{}, s Serializer, d Deserializer, o []BuildOption) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.registered[key]; exists {
		return fmt.Errorf("something with the key `%s` has already been registered", key)
	}

	r.registered[key] = registered{
		factory:      fn,
		serializer:   s,
		deserializer: d,
		options:      o,
	}

	return nil
}
