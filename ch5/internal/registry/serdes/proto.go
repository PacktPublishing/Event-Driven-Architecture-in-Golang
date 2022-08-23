package serdes

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"

	"eda-in-golang/internal/registry"
)

type ProtoSerde struct {
	r registry.Registry
}

var _ registry.Serde = (*ProtoSerde)(nil)
var protoT = reflect.TypeOf((*proto.Message)(nil)).Elem()

func NewProtoSerde(r registry.Registry) *ProtoSerde {
	return &ProtoSerde{r: r}
}

func (c ProtoSerde) Register(v registry.Registrable, options ...registry.BuildOption) error {
	if !reflect.TypeOf(v).Implements(protoT) {
		return fmt.Errorf("%T does not implement proto.Message", v)
	}
	return registry.Register(c.r, v, c.serialize, c.deserialize, options)
}

func (c ProtoSerde) RegisterKey(key string, v interface{}, options ...registry.BuildOption) error {
	if !reflect.TypeOf(v).Implements(protoT) {
		return fmt.Errorf("%T does not implement proto.Message", v)
	}
	return registry.RegisterKey(c.r, key, v, c.serialize, c.deserialize, options)
}

func (c ProtoSerde) RegisterFactory(key string, fn func() interface{}, options ...registry.BuildOption) error {
	if v := fn(); v == nil {
		return fmt.Errorf("%s factory returns a nil value", key)
	} else if _, ok := v.(proto.Message); !ok {
		return fmt.Errorf("%s does not implement proto.Message", key)
	}
	return registry.RegisterFactory(c.r, key, fn, c.serialize, c.deserialize, options)
}

func (ProtoSerde) serialize(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func (ProtoSerde) deserialize(data []byte, v interface{}) error {
	return proto.Unmarshal(data, v.(proto.Message))
}
