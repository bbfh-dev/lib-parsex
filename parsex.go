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

			switch option.Name {
			case internal.HelpOption.Name:
				return program, PrintHelpErr
			case internal.VersionOption.Name:
				return program, PrintVersionErr
			}

			were_modified = append(were_modified, option.Name)

			if option.Type == reflect.Bool {
				option.Ref.SetBool(true)
				continue
			}

			if len(parts) == 2 {
				if err := option.Set(parts[1]); err != nil {
					return program, err
				}
				continue
			}

			if i+1 >= len(args) {
				return program, fmt.Errorf(
					"option %q requires a value <%s>. %s",
					args[i],
					option.Type.String(),
					PrintHelpErr.Error(),
				)
			}
			i++

			if err := option.Set(args[i]); err != nil {
				return program, err
			}
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

	// TODO: Validate and save input
	fmt.Println(input)

	return program, nil
}

func parseCommand(command *Program, args []string) (*Program, error) {
	if err := command.Parse(); err != nil {
		return command, err
	}
	return parseInput(command, args)
}
