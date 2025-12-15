package libparsex

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
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

	err = program.EntryPoint()
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

	were_modified := make([]string, 0, len(args))
	input := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		// Regular argument or command
		if !strings.HasPrefix(args[i], "-") {
			command, ok := program.parsedCommands[args[i]]
			if !ok {
				input = append(input, args[i])
				continue
			}

			return parseCommand(command, args[i+1:])
		}

		// Long option
		if prefix := "--"; strings.HasPrefix(args[i], prefix) {
			arg := args[i][len(prefix):]

			if arg == "" {
				// Everything after '--' is a positional argument
				input = append(input, args[i+1:]...)
				break
			}

			parts := strings.SplitN(arg, "=", 2)
			arg = parts[0]

			option, ok := program.parsedOptionMap[arg]
			if !ok {
				return program, fmt.Errorf("unknown option %q. %s", args[i], PrintHelpErr.Error())
			}

			were_modified = append(were_modified, option.Name)

			if err := parseLongOption(option, parts, args, &i); err != nil {
				return program, err
			}
			continue
		}

		// Short options
		arg := args[i][1:]
		if arg == "" {
			input = append(input, args[i])
			continue
		}

		parts := strings.SplitN(arg, "=", 2)
		arg = parts[0]

		option, ok := program.parsedOptionMap[arg]
		if ok {
			were_modified = append(were_modified, option.Name)
			if err := parseLongOption(option, parts, args, &i); err != nil {
				return program, err
			}
			continue
		}

		if len(arg) == 1 {
			option, ok := program.parsedAltMap[arg]
			if !ok {
				return program, fmt.Errorf("unknown option %q. %s", args[i], PrintHelpErr.Error())
			}
			were_modified = append(were_modified, option.Name)
			if option.Type == reflect.Bool {
				option.Ref.SetBool(true)
				continue
			}
			if err := parseLongOption(option, parts, args, &i); err != nil {
				return program, err
			}
			continue
		}

		// combined flag
		for char := range strings.SplitSeq(arg, "") {
			option, ok := program.parsedAltMap[char]
			if !ok {
				return program, fmt.Errorf(
					"provided option %q is neither a long or combined option. %s",
					args[i],
					PrintHelpErr.Error(),
				)
			}
			if option.Type != reflect.Bool {
				return program, fmt.Errorf(
					"provided combined flag %q contains a non-flag option %q. %s",
					args[i],
					option.Name,
					PrintHelpErr.Error(),
				)
			}
			option.Ref.SetBool(true)
		}
	}

	// Set default values
	for _, option := range program.parsedOptions {
		if slices.Contains(were_modified, option.Name) {
			continue
		}

		if option.Default != nil {
			if err := option.Set(*option.Default); err != nil {
				return program, err
			}
		}
	}

	for i, arg := range program.parsedArgs {
		if i >= len(input) {
			if arg.Tag != internal.ARG_DEFAULT {
				break
			}
			return program, fmt.Errorf(
				"expected an argument %q but got none. %s",
				arg.String(),
				PrintHelpErr.Error(),
			)
		}
		switch arg.Tag {

		case internal.ARG_VARIADIC:
			arg.Ref.Set(reflect.ValueOf(input[i:]))

		default:
			arg.Set(input[i])
		}
	}

	return program, nil
}

func parseCommand(command *Program, args []string) (*Program, error) {
	if err := command.Parse(); err != nil {
		return command, err
	}
	return parseInput(command, args)
}

func parseLongOption(option *internal.ParsedOption, parts []string, args []string, i *int) error {
	switch option.Name {
	case internal.HelpOption.Name:
		return PrintHelpErr
	case internal.VersionOption.Name:
		return PrintVersionErr
	}

	if option.Type == reflect.Bool {
		option.Ref.SetBool(true)
		return nil
	}

	if len(parts) == 2 {
		if err := option.Set(parts[1]); err != nil {
			return err
		}
		return nil
	}

	if *i+1 >= len(args) {
		return fmt.Errorf(
			"option %q requires a value <%s>. %s",
			args[*i],
			option.Type.String(),
			PrintHelpErr.Error(),
		)
	}

	*i++
	if err := option.Set(args[*i]); err != nil {
		return err
	}

	return nil
}
