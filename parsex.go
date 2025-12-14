package libparsex

import (
	"errors"
	"os"
)

// Return this error in [EntryPoint] to print help message and exit.
var PrintHelpErr = errors.New("--help")

func Run(program *Program, args []string) error {
	if !program.didParse {
		if err := program.Parse(); err != nil {
			return err
		}
	}

	if err := parseInput(program, args); err != nil {
		return err
	}

	err := program.EntryPoint()
	if err == PrintHelpErr {
		program.PrintHelp(os.Stdout)
		return nil
	}

	return err
}

func parseInput(program *Program, args []string) error {
	if len(args) == 0 {
		return nil
	}

	return nil
}
