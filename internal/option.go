package internal

import (
	"fmt"
	"reflect"
	"strings"
)

var HelpOption = ParsedOption{
	Name:    "help",
	Type:    reflect.Bool,
	Alt:     "",
	Desc:    "Print this help message and exit",
	Default: nil,
	Ref:     nil,
}

var VersionOption = ParsedOption{
	Name:    "help",
	Type:    reflect.Bool,
	Alt:     "",
	Desc:    "Print this help message and exit",
	Default: nil,
	Ref:     nil,
}

type ParsedOption struct {
	Name    string
	Type    reflect.Kind
	Alt     string
	Desc    string
	Default any

	Ref *reflect.Value
}

func (option *ParsedOption) String() string {
	var builder strings.Builder
	builder.WriteString("--" + option.Name)

	if option.Alt != "" {
		builder.WriteString(", -" + option.Alt)
	}

	if option.Type != reflect.Bool {
		builder.WriteString(" <" + option.Type.String() + ">")
	}

	if option.Default != nil {
		switch default_value := option.Default.(type) {
		case string:
			fmt.Fprintf(&builder, " (default: %q)", default_value)
		default:
			fmt.Fprintf(&builder, " (default: %v)", default_value)
		}
	}

	builder.WriteString("\n")
	if option.Desc != "" {
		builder.WriteString("        # " + option.Desc)
	}

	return builder.String()
}
