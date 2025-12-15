package internal

import (
	"fmt"
	"reflect"
	"strconv"
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
	Default *string

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
		fmt.Fprintf(&builder, " (default: %v)", *option.Default)
	}

	builder.WriteString("\n")
	if option.Desc != "" {
		builder.WriteString("        # " + option.Desc + "\n")
	}

	return builder.String()
}

func (option *ParsedOption) wrapError(err error) error {
	return fmt.Errorf(
		"option '%s <%s>': %w",
		option.Name,
		option.Type.String(),
		err,
	)
}

func (option *ParsedOption) Set(value string) error {
	switch option.Type {

	case reflect.String:
		option.Ref.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.Atoi(value)
		if err != nil {
			return option.wrapError(err)
		}
		option.Ref.SetInt(int64(value))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.Atoi(value)
		if err != nil {
			return option.wrapError(err)
		}
		option.Ref.SetUint(uint64(value))

	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return option.wrapError(err)
		}
		option.Ref.SetFloat(float64(value))
	}

	return nil
}
