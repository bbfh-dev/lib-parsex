package internal

import (
	"reflect"
	"strconv"
)

type Ref struct {
	Kind    reflect.Kind
	Pointer *reflect.Value
}

func (ref Ref) Set(value string) error {
	switch ref.Kind {

	case reflect.Bool:
		ref.Pointer.SetBool(true)

	case reflect.String:
		ref.Pointer.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		ref.Pointer.SetInt(int64(value))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		ref.Pointer.SetUint(uint64(value))

	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		ref.Pointer.SetFloat(float64(value))
	}

	return nil
}
