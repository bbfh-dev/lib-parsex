package libparsex

import "strings"

type Argument struct {
	// The name of the argument as it will be displayed in --help.
	Name string
	// Whether this argument can be ommited.
	//
	// NOTE: IsOptional arguments can only be in the end.
	IsOptional bool
	// Whether this argument can repeat 1 or more times.
	// Combine with [.IsOptional] to allow for 0 repetitions.
	//
	// NOTE: There can only be one variadic argument and it must be in the end.
	IsVariadic bool
}

func (arg Argument) String() string {
	var builder strings.Builder
	// grow by name + <> + ? + ...
	builder.Grow(len(arg.Name) + 2 + 1 + 3)
	builder.WriteString("<" + arg.Name)

	if arg.IsOptional {
		builder.WriteString("?")
	}

	if arg.IsVariadic {
		builder.WriteString("...")
	}

	builder.WriteString(">")
	return builder.String()
}
