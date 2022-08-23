package serdes

import (
	"encoding/json"

	"eda-in-golang/internal/registry"
)

type JsonSerde struct {
	r registry.Registry
}

var _ registry.Serde = (*JsonSerde)(nil)

func NewJsonSerde(r registry.Registry) *JsonSerde {
	return &JsonSerde{r: r}
}

func (c JsonSerde) Register(v registry.Registrable, options ...registry.BuildOption) error {
	return registry.Register(c.r, v, c.serialize, c.deserialize, options)
}

func (c JsonSerde) RegisterKey(key string, v interface{}, options ...registry.BuildOption) error {
	return registry.RegisterKey(c.r, key, v, c.serialize, c.deserialize, options)
}

func (c JsonSerde) RegisterFactory(key string, fn func() interface{}, options ...registry.BuildOption) error {
	return registry.RegisterFactory(c.r, key, fn, c.serialize, c.deserialize, options)
}

func (JsonSerde) serialize(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (JsonSerde) deserialize(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
