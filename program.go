package libparsex

import (
	"fmt"
	"reflect"

	"github.com/bbfh-dev/lib-parsex/v3/internal"
	"github.com/iancoleman/strcase"
)

type Program struct {
	// The name of the binary/subcommand
	Name string
	// SemVer, e.g. 0.1.2-beta.1 (no 'v' prefix).
	// Leave empty for subcommands.
	Version string
	// Brief description of the binary/subcommand.
	Description string
	// Pointer to a struct{} containing all options and flags.
	//
	// The struct{} must contain fields with the following datatypes allowed:
	//
	// int-s, uint-s, float-s, string, bool
	//
	// Each field can have the following tags:
	//
	// — alt:"" (Optional) This is a SINGLE character alternative to the option/flag.
	//
	// — desc:"" (Optional) A brief description of the option/flag.
	//
	// — default:"" (Optional) The default value of the flag.
	Options any
	// Pointer to a struct{} containing all positional arguments.
	//
	// The struct{} must contain fields with the following datatypes allowed:
	//
	// int-s, uint-s, float-s, string.
	//
	// Use slices to indicate variadic arguments, e.g. []string.
	//
	// NOTE: Only a single variadic argument is allowed and it must be in the end.
	// NOTE: Variadic arguments are ALWAYS optional.
	//
	// Use pointers to indicate that the argument is optional, e.g. *string.
	//
	// NOTE: Optional arguments must be in the end.
	Args any
	// Commands are other sub-programs that have their own nested options and arguments.
	//
	// NOTE: Parent's .Options will be parsed and set up and until the command name is mentioned in the input.
	Commands []*Program
	// The function that will be called when running the Program.
	EntryPoint func() error

	didParse        bool
	parsedOptions   []*internal.ParsedOption
	parsedOptionMap map[string]*internal.ParsedOption
	parsedArgs      []*internal.ParsedArgument
}

func (program *Program) HasVersion() bool {
	return program.Version != ""
}

func (program *Program) Parse() error {
	// Built-in values
	program.parsedOptions = []*internal.ParsedOption{
		&internal.HelpOption,
		&internal.VersionOption,
	}
	program.parsedOptionMap = map[string]*internal.ParsedOption{
		internal.HelpOption.Name:    &internal.HelpOption,
		internal.VersionOption.Name: &internal.VersionOption,
	}
	program.parsedArgs = []*internal.ParsedArgument{}

	if err := program.parseOptions(); err != nil {
		return err
	}

	if err := program.parseArgs(); err != nil {
		return err
	}

	program.didParse = true
	return nil
}

func (program *Program) expectStructPtr(name string, reflect_type reflect.Type) error {
	if reflect_type.Kind() != reflect.Pointer {
		return fmt.Errorf(
			"(parsex) [%s] expected a pointer but got %q",
			name,
			reflect_type.Kind().String(),
		)
	}

	reflect_elem := reflect_type.Elem()
	if reflect_elem.Kind() != reflect.Struct {
		return fmt.Errorf(
			"(parsex) [%s] expected a to point to struct{} but it points to %q",
			reflect_type.Name(),
			reflect_elem.Kind().String(),
		)
	}

	return nil
}

func (program *Program) parseOptions() error {
	if err := program.expectStructPtr("Options", reflect.TypeOf(program.Options)); err != nil {
		return err
	}

	reflect_type := reflect.TypeOf(program.Options).Elem()
	reflect_value := reflect.ValueOf(program.Options).Elem()

	for i := range reflect_type.NumField() {
		field := reflect_type.Field(i)
		reflect_field := reflect_value.Field(i)

		var default_value any
		if value, ok := field.Tag.Lookup("default"); ok {
			default_value = value
		}

		option := &internal.ParsedOption{
			Name:    strcase.ToKebab(field.Name),
			Type:    field.Type.Kind(),
			Alt:     field.Tag.Get("alt"),
			Desc:    field.Tag.Get("desc"),
			Default: default_value,
			Ref:     &reflect_field,
		}
		program.parsedOptions = append(program.parsedOptions, option)
		program.parsedOptionMap[option.Name] = option
	}

	return nil
}

func (program *Program) parseArgs() error {
	if err := program.expectStructPtr("Args", reflect.TypeOf(program.Args)); err != nil {
		return err
	}

	reflect_type := reflect.TypeOf(program.Args).Elem()
	reflect_value := reflect.ValueOf(program.Args).Elem()

	for i := range reflect_type.NumField() {
		field := reflect_type.Field(i)
		reflect_field := reflect_value.Field(i)

		tag := internal.ARG_DEFAULT
		switch field.Type.Kind() {
		case reflect.Pointer:
			tag = internal.ARG_OPTIONAL
		case reflect.Slice:
			tag = internal.ARG_VARIADIC
		}

		arg := &internal.ParsedArgument{
			Name: field.Name,
			Kind: field.Type.Kind(),
			Tag:  tag,
			Ref:  &reflect_field,
		}
		program.parsedArgs = append(program.parsedArgs, arg)
	}

	return nil
}
