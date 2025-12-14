package libparsex

import (
	"bytes"
	"io"
	"strings"
)

var printIndent = bytes.Repeat([]byte{' '}, 4)
var printNewLine = []byte{'\n'}

func (program *Program) String() string {
	var builder strings.Builder
	program.PrintHelp(&builder)
	return builder.String()
}

func (program *Program) PrintVersion(writer io.Writer) {
	writer.Write([]byte(program.Name))
	if program.HasVersion() {
		writer.Write([]byte(" " + program.Version))
	}
	writer.Write(printNewLine)
}

func (program *Program) PrintHelp(writer io.Writer) {
	program.PrintVersion(writer)
	writer.Write(printNewLine)
	if program.Description != "" {
		writer.Write([]byte(program.Description))
		writer.Write(printNewLine)
	}

	writer.Write(printNewLine)
	writer.Write([]byte("[?] Usage:"))
	writer.Write(printNewLine)
	writer.Write(printIndent)
	writer.Write([]byte(program.Name + " [options...]"))
	for _, arg := range program.parsedArgs {
		writer.Write([]byte(" " + arg.String()))
	}
	writer.Write(printNewLine)

	writer.Write(printNewLine)
	writer.Write([]byte("[#] Options:"))
	writer.Write(printNewLine)
	for _, option := range program.parsedOptions {
		writer.Write(printIndent)
		writer.Write([]byte(option.String()))
	}
}
