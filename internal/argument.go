package internal

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ArgumentTag byte

const (
	ARG_DEFAULT ArgumentTag = iota
	ARG_OPTIONAL
	ARG_VARIADIC
)

type ParsedArgument struct {
	Name string
	Kind reflect.Kind
	Tag  ArgumentTag

	Ref *reflect.Value
}

func (arg *ParsedArgument) String() string {
	var builder strings.Builder
	// longest possible variant: <name...>. i.e. len(name) + <> + ...
	builder.Grow(len(arg.Name) + 2 + 3)

	builder.WriteString("<" + arg.Name)
	switch arg.Tag {
	case ARG_OPTIONAL:
		builder.WriteString("?")
	case ARG_VARIADIC:
		builder.WriteString("...")
	}

	builder.WriteString(">")
	return builder.String()
}

func (arg *ParsedArgument) wrapError(err error) error {
	return fmt.Errorf(
		"argument %q: %w",
		arg.String(),
		err,
	)
}

func (arg *ParsedArgument) Set(value string) error {
	switch arg.Kind {

	case reflect.String:
		arg.Ref.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.Atoi(value)
		if err != nil {
			return arg.wrapError(err)
		}
		arg.Ref.SetInt(int64(value))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.Atoi(value)
		if err != nil {
			return arg.wrapError(err)
		}
		arg.Ref.SetUint(uint64(value))

	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return arg.wrapError(err)
		}
		arg.Ref.SetFloat(float64(value))
	}

	return nil
}
