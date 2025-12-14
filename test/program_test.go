package libparsex_test

import (
	"testing"

	libparsex "github.com/bbfh-dev/lib-parsex/v3"
	"gotest.tools/assert"
)

func TestProgram(test *testing.T) {
	var options struct {
		Verbose bool `desc:"Print verbose debug information"`
	}

	var args struct {
		count int
		input []string
	}

	program := libparsex.Program{
		Name:        "example",
		Version:     "0.1.2-beta.1",
		Description: "This is an example program",
		Options:     &options,
		Args:        &args,
		EntryPoint: func() error {
			return nil
		},
	}

	err := libparsex.Run(&program, []string{})
	assert.NilError(test, err)
}
