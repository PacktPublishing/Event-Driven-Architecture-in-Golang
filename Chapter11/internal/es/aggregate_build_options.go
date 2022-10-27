package es

import (
	"fmt"

	"eda-in-golang/internal/registry"
)

type VersionSetter interface {
	SetVersion(int)
}

func SetVersion(version int) registry.BuildOption {
	return func(v interface{}) error {
		if agg, ok := v.(VersionSetter); ok {
			agg.SetVersion(version)
			return nil
		}
		return fmt.Errorf("%T does not have the method setVersion(int)", v)
	}
}
