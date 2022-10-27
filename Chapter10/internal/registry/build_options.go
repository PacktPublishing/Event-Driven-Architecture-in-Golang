package registry

import (
	"fmt"
	"reflect"
)

type BuildOption func(v interface{}) error

// ValidateImplements will verify that the built value implements the given interface
//
// Pass interfaces using "(*I)(nil)", like in this example:
//
//   type MyInterface interface {
//       MyMethod(int) bool
//   }
//
//   ValidateImplements((*MyInterface)(nil))
//
// You will only be able to validate interfaces with exported methods, you will not be
// able to validate any that have one or more unexported methods.
func ValidateImplements(checkV interface{}) BuildOption {
	checkT := reflect.TypeOf(checkV)

	if checkT.Kind() == reflect.Ptr {
		checkT = checkT.Elem()
	}

	if checkT.Kind() != reflect.Interface {
		panic(fmt.Sprintf("%T is not an interface", checkV))
	}

	return func(v interface{}) error {
		t := reflect.TypeOf(v)

		if !t.Implements(checkT) {
			return fmt.Errorf("%T does not implement %T", v, checkV)
		}

		return nil
	}
}
