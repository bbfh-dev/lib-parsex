package libparsex_test

import (
	"strings"
	"testing"

	libparsex "github.com/bbfh-dev/lib-parsex/v3"
	"gotest.tools/assert"
)

type ArgTestCase struct {
	Input         []string
	ExpectOptions ExpectedOptions
	ExpectArgs    ExpectedArgs
	ExpectErr     string
}

var ArgTestCases = []ArgTestCase{
	{
		Input: []string{"--verbose", "0"},
		ExpectOptions: ExpectedOptions{
			Verbose:    true,
			OtherValue: 69,
		},
		ExpectArgs: ExpectedArgs{
			Input: []string{},
		},
	},
	{
		Input: []string{"--stdin-file-path", "Hello World!", "0"},
		ExpectOptions: ExpectedOptions{
			StdinFilePath: "Hello World!",
			OtherValue:    69,
		},
		ExpectArgs: ExpectedArgs{
			Input: []string{},
		},
	},
	{
		Input:     []string{"--stdin-file-path", "Hello World!", "--verbose", "--other-value"},
		ExpectErr: "--other-value",
	},
	{
		Input: []string{"--stdin-file-path", "Hello World!", "--verbose", "--other-value=123", "0"},
		ExpectOptions: ExpectedOptions{
			StdinFilePath: "Hello World!",
			Verbose:       true,
			OtherValue:    123,
		},
		ExpectArgs: ExpectedArgs{
			Input: []string{},
		},
	},
	{
		Input: []string{"-v", "-o=123", "0"},
		ExpectOptions: ExpectedOptions{
			Verbose:    true,
			OtherValue: 123,
		},
		ExpectArgs: ExpectedArgs{
			Input: []string{},
		},
	},
	{
		Input: []string{"-v", "-o", "123", "0"},
		ExpectOptions: ExpectedOptions{
			Verbose:    true,
			OtherValue: 123,
		},
		ExpectArgs: ExpectedArgs{
			Input: []string{},
		},
	},
	{
		Input:     []string{"-vo"},
		ExpectErr: "contains a non-flag",
	},
	{
		Input: []string{"-v", "123"},
		ExpectOptions: ExpectedOptions{
			Verbose:    true,
			OtherValue: 69,
		},
		ExpectArgs: ExpectedArgs{
			Count: 123,
			Input: []string{},
		},
	},
	{
		Input: []string{"-v", "123", "file1"},
		ExpectOptions: ExpectedOptions{
			Verbose:    true,
			OtherValue: 69,
		},
		ExpectArgs: ExpectedArgs{
			Count: 123,
			Input: []string{"file1"},
		},
	},
	{
		Input: []string{"123", "file1", "file2", "file3"},
		ExpectOptions: ExpectedOptions{
			OtherValue: 69,
		},
		ExpectArgs: ExpectedArgs{
			Count: 123,
			Input: []string{"file1", "file2", "file3"},
		},
	},
	{
		Input:     []string{"-v"},
		ExpectErr: "expected an argument",
	},
}

func TestArgs(test *testing.T) {
	for _, test_case := range ArgTestCases {
		test.Run(strings.Join(test_case.Input, "__"), func(test *testing.T) {
			Options.Verbose = false
			Options.StdinFilePath = ""
			Options.OtherValue = 0
			Args.Count = 0
			Args.Input = []string{}

			err := libparsex.Run(&Program, test_case.Input)

			if test_case.ExpectErr != "" {
				assert.ErrorContains(test, err, test_case.ExpectErr)
			} else {
				assert.NilError(test, err)
				assert.DeepEqual(test, test_case.ExpectOptions.Verbose, Options.Verbose)
				assert.DeepEqual(test, test_case.ExpectOptions.StdinFilePath, Options.StdinFilePath)
				assert.DeepEqual(test, test_case.ExpectOptions.OtherValue, Options.OtherValue)
				assert.DeepEqual(test, test_case.ExpectArgs.Count, Args.Count)
				assert.DeepEqual(test, test_case.ExpectArgs.Input, Args.Input)
			}
		})
	}
}
