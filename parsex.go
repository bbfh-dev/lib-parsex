package libparsex

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/bbfh-dev/lib-parsex/v3/internal"
)

// Return this error in [EntryPoint] to print help message and exit.
var PrintHelpErr = errors.New("Refer to --help for usage information")

// Return this error in [EntryPoint] to print version message and exit.
var PrintVersionErr = errors.New("--version")

func Run(program *Program, args []string) error {
	if !program.didParse {
		if err := program.Parse(); err != nil {
			return err
		}
	}

	var err error
	program, err = parseInput(program, args)
	if err != nil {
		goto handle_error
	}

	err = program.EntryPoint(program.rawArgs)
handle_error:
	switch err {
	case PrintHelpErr:
		program.PrintHelp(os.Stdout)
	case PrintVersionErr:
		program.PrintVersion(os.Stdout)
	default:
		return err
	}

	return nil
}

func parseInput(program *Program, args []string) (*Program, error) {
	if len(args) == 0 {
		return program, nil
	}

	modifiedOptions := make(map[string]bool)
	positionalArgs := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if !strings.HasPrefix(arg, "-") {
			if cmd, ok := program.parsedCommands[arg]; ok {
				return parseCommand(cmd, args[i+1:])
			}
			positionalArgs = append(positionalArgs, arg)
			continue
		}

		if arg == "--" {
			positionalArgs = append(positionalArgs, args[i+1:]...)
			break
		}

		if strings.HasPrefix(arg, "--") {
			if err := handleLongOption(program, arg[2:], args, &i, modifiedOptions); err != nil {
				return program, err
			}
			continue
		}

		if err := handleShortOption(program, arg[1:], args, &i, modifiedOptions); err != nil {
			return program, err
		}
	}

	if err := applyDefaultOptions(program, modifiedOptions); err != nil {
		return program, err
	}

	if err := assignPositionalArgs(program, positionalArgs); err != nil {
		return program, err
	}
	program.rawArgs = positionalArgs

	return program, nil
}

func parseCommand(command *Program, args []string) (*Program, error) {
	if err := command.Parse(); err != nil {
		return command, err
	}
	return parseInput(command, args)
}

func handleLongOption(
	program *Program,
	raw string,
	args []string,
	index *int,
	modified map[string]bool,
) error {
	parts := strings.SplitN(raw, "=", 2)
	name := parts[0]

	option, ok := program.parsedOptionMap[name]
	if !ok {
		return fmt.Errorf("unknown option %q. %s", "--"+raw, PrintHelpErr.Error())
	}

	modified[option.Name] = true
	return parseOptionValue(option, parts, args, index)
}

func handleShortOption(
	program *Program,
	raw string,
	args []string,
	index *int,
	modified map[string]bool,
) error {
	if raw == "" {
		return nil
	}

	parts := strings.SplitN(raw, "=", 2)
	name := parts[0]

	if option, ok := program.parsedOptionMap[name]; ok {
		modified[option.Name] = true
		return parseOptionValue(option, parts, args, index)
	}

	if len(name) == 1 {
		option, ok := program.parsedAltMap[name]
		if !ok {
			return fmt.Errorf("unknown option %q. %s", "-"+raw, PrintHelpErr.Error())
		}
		modified[option.Name] = true

		if option.Kind == reflect.Bool {
			option.Ref.Pointer.SetBool(true)
			return nil
		}

		return parseOptionValue(option, parts, args, index)
	}

	for _, char := range name {
		option, ok := program.parsedAltMap[string(char)]
		if !ok {
			return fmt.Errorf(
				"provided option %q is neither a long or combined option. %s",
				"-"+raw,
				PrintHelpErr.Error(),
			)
		}
		if option.Kind != reflect.Bool {
			return fmt.Errorf(
				"provided combined flag %q contains a non-flag option %q. %s",
				"-"+raw,
				option.Name,
				PrintHelpErr.Error(),
			)
		}
		option.Ref.Pointer.SetBool(true)
		modified[option.Name] = true
	}

	return nil
}

func parseOptionValue(
	option *internal.ParsedOption,
	parts []string,
	args []string,
	index *int,
) error {
	switch option.Name {
	case internal.HelpOption.Name:
		return PrintHelpErr
	case internal.VersionOption.Name:
		return PrintVersionErr
	}

	if option.Kind == reflect.Bool {
		option.Ref.Pointer.SetBool(true)
		return nil
	}

	if len(parts) == 2 {
		return option.Set(parts[1])
	}

	if *index+1 >= len(args) {
		return fmt.Errorf(
			"option %q requires a value <%s>. %s",
			args[*index],
			option.Kind.String(),
			PrintHelpErr.Error(),
		)
	}

	*index++
	return option.Set(args[*index])
}

func applyDefaultOptions(program *Program, modified map[string]bool) error {
	for _, option := range program.parsedOptions {
		if modified[option.Name] || option.Default == nil {
			continue
		}
		if err := option.Set(*option.Default); err != nil {
			return err
		}
	}
	return nil
}

func assignPositionalArgs(program *Program, input []string) error {
	for i, arg := range program.parsedArgs {
		if i >= len(input) {
			if arg.Tag == internal.ARG_DEFAULT {
				return fmt.Errorf(
					"expected an argument %q but got none. %s",
					arg.String(),
					PrintHelpErr.Error(),
				)
			}
			break
		}

		if arg.Tag == internal.ARG_VARIADIC {
			arg.Ref.Pointer.Set(reflect.ValueOf(input[i:]))
			return nil
		}

		arg.Set(input[i])
	}
	return nil
}
