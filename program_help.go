package libparsex

import (
	"fmt"
	"io"
	"strings"

	libescapes "github.com/bbfh-dev/lib-ansi-escapes"
)

var printIndent = strings.Repeat(" ", 4)

func (program *Program) String() string {
	return fmt.Sprintf(
		"%[1]s%[2]s\n%[1]s%[1]s%[4]s# %[3]s%[5]s\n",
		printIndent,
		program.Name,
		program.Description,
		libescapes.Optional(libescapes.TextColorWhite),
		libescapes.Optional(libescapes.ColorReset),
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

	printHeader(writer, "\n[?] Usage:\n", libescapes.TextColorBrightYellow)
	fmt.Fprint(writer, printIndent)
	fmt.Fprint(writer, program.Name, " [options...]")

	for _, argument := range program.parsedArgs {
		fmt.Fprint(writer, " ", argument.String())
	}
	fmt.Fprint(writer, "\n")

	if len(program.Commands) > 0 {
		printHeader(writer, "\n[>] Commands:\n", libescapes.TextColorBrightBlue)
		for _, command := range program.Commands {
			fmt.Fprint(writer, command.String())
		}
	}

	printHeader(writer, "\n[#] Options:\n", libescapes.TextColorBrightMagenta)
	for _, option := range program.parsedOptions {
		fmt.Fprint(writer, printIndent)
		fmt.Fprint(writer, option.String())
	}
}

func printHeader(writer io.Writer, text, color string) {
	writer.Write([]byte(libescapes.Optional(color) +
		text +
		libescapes.Optional(libescapes.ColorReset)))
}
