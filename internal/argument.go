package internal

import (
	"reflect"
	"strings"
)

type ArgumentKind byte

const (
	ARG_DEFAULT ArgumentKind = iota
	ARG_OPTIONAL
	ARG_VARIADIC
)

type ParsedArgument struct {
	Name string
	Kind ArgumentKind

	Ref *reflect.Value
}

func (arg *ParsedArgument) String() string {
	var builder strings.Builder
	// longest possible variant: <name...>. i.e. len(name) + <> + ...
	builder.Grow(len(arg.Name) + 2 + 3)

	builder.WriteString("<" + arg.Name)
	switch arg.Kind {
	case ARG_OPTIONAL:
		builder.WriteString("?")
	case ARG_VARIADIC:
		builder.WriteString("...")
	}

	builder.WriteString(">")
	return builder.String()
}
