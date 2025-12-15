# Parsex

Parsex `/pɑːrsɛks/` — a GNU-/POSIX-compiant CLI argument parsing and validation library.

## Usage

```bash
# Get parsex
go get -u github.com/bbfh-dev/lib-parsex/v3
```

```go
var Options struct {
    // --verbose, -v
    Verbose       bool   `alt:"v" desc:"Print verbose debug information"`
    // --stdin-file-path
    StdinFilePath string `        desc:"Path to the file to pretend that stdin comes from"`
    // --other-value, -o
    OtherValue    int    `alt:"o" default:"69"`
}

var Args struct {
    // <count>
    Count int
    // <input...>
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
        // Your program logic

        if len(rawArgs) == 0 {
            // Print --help when no arguments are provided
            return libparsex.PrintHelpErr
        }

        // It is recommended to use your [Args] because it respects the datatypes.
        // rawArgs is mostly for edge-cases where a slice is more convinient
        // (e.g. checking the length, as in this example)
        fmt.Prinf("count argument is: %d", Args.Count)
        return nil
    },
}

func TestHelp(test *testing.T) {
    err := libparsex.Run(&Program, os.Args[1:])
    if err != nil {
        os.Stderr.WriteString(err.Error())
        os.Exit(1)
    }
}
```
