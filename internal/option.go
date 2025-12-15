package internal

import (
	"fmt"
	"reflect"
	"strings"
)

var HelpOption = ParsedOption{
	Name:    "help",
	Alt:     "",
	Desc:    "Print this help message and exit",
	Default: nil,
	Ref: Ref{
		Kind:    reflect.Bool,
		Pointer: nil,
	},
}

var VersionOption = ParsedOption{
	Name:    "help",
	Alt:     "",
	Desc:    "Print this help message and exit",
	Default: nil,
	Ref: Ref{
		Kind:    reflect.Bool,
		Pointer: nil,
	},
}

type ParsedOption struct {
	Name    string
	Alt     string
	Desc    string
	Default *string
	Ref
}

func (option *ParsedOption) String() string {
	var builder strings.Builder
	builder.WriteString("--" + option.Name)

	if option.Alt != "" {
		builder.WriteString(", -" + option.Alt)
	}

	if option.Kind != reflect.Bool {
		builder.WriteString(" <" + option.Kind.String() + ">")
	}

	if option.Default != nil {
		fmt.Fprintf(&builder, " (default: %v)", *option.Default)
	}

	builder.WriteString("\n")
	if option.Desc != "" {
		builder.WriteString("        # " + option.Desc + "\n")
	}

	return builder.String()
}

func (option *ParsedOption) WrapError(err error) error {
	return fmt.Errorf(
		"option '%s <%s>': %w",
		option.Name,
		option.Kind.String(),
		err,
	)
}
