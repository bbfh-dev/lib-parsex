package libparsex_test

import (
	"errors"
	"testing"

	libparsex "github.com/bbfh-dev/lib-parsex/v3"
	"gotest.tools/assert"
)

var DidRun bool

type ExpectedOptions struct {
	Verbose       bool
	StdinFilePath string
	OtherValue    int
}

var Options struct {
	Verbose       bool   `alt:"v" desc:"Print verbose debug information"`
	StdinFilePath string `        desc:"Path to the file to pretend that stdin comes from"`
	OtherValue    int    `alt:"o" default:"69"`
}

type ExpectedArgs struct {
	Count int
	Input []string
}

var Args struct {
	Count int
	Input []string
}

var Program = libparsex.Program{
	Name:        "example",
	Version:     "0.1.2-beta.1",
	Description: "This is an example program",
	Options:     &Options,
	Args:        &Args,
	Commands: []*libparsex.Program{
		{Name: "nested", Description: "Example nested command"},
	},
	EntryPoint: func(rawArgs []string) error {
		DidRun = true
		if len(rawArgs) == 0 {
			return errors.New("this is wrong")
		}
		return nil
	},
}

func TestHelp(test *testing.T) {
	err := libparsex.Run(&Program, []string{"--help"})
	assert.NilError(test, err)
}

func TestCall(test *testing.T) {
	DidRun = false
	err := libparsex.Run(&Program, []string{})
	assert.NilError(test, err)
	assert.Equal(test, DidRun, true)
}
