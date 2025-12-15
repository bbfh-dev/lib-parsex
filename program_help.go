package libparsex

import (
	"fmt"
	"io"
	"strings"
)

var printIndent = strings.Repeat(" ", 4)

func (program *Program) String() string {
	return fmt.Sprintf(
		"%[1]s%[2]s\n%[1]s%[1]s# %[3]s\n",
		printIndent,
		program.Name,
		program.Description,
	)
}

func (program *Program) PrintVersion(writer io.Writer) {
	writer.Write([]byte(program.Name))
	if program.HasVersion() {
		writer.Write([]byte(" " + program.Version))
	}
	writer.Write([]byte{'\n'})
}

func (program *Program) PrintHelp(writer io.Writer) {
	program.PrintVersion(writer)

	if program.Description != "" {
		fmt.Fprintln(writer, "\n"+program.Description)
	}

	fmt.Fprintln(writer, "\n[?] Usage:")
	fmt.Fprint(writer, printIndent)
	fmt.Fprint(writer, program.Name, " [options...]")

	for _, argument := range program.parsedArgs {
		fmt.Fprint(writer, " ", argument.String())
	}
	fmt.Fprint(writer, "\n")

	if len(program.Commands) > 0 {
		fmt.Fprintln(writer, "\n[>] Commands:")
		for _, command := range program.Commands {
			fmt.Fprint(writer, command.String())
		}
	}

	fmt.Fprintln(writer, "\n[#] Options:")
	for _, option := range program.parsedOptions {
		fmt.Fprint(writer, printIndent)
		fmt.Fprint(writer, option.String())
	}
}
