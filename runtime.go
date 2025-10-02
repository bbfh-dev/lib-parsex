package parsex

import (
	"os"
	"strings"

	"github.com/bbfh-dev/berr"
	"github.com/bbfh-dev/parsex/v2/internal"
)

// [runtimeType] is created from [Program] and it's what actually handles everything
type runtimeType struct {
	// These are set from [Program]
	data any
	exec *internal.ContextExecutable
	name string
	desc string
	// (Optional) Only needed for the primary executable program.
	// SemVer is adviced. Use [Program.SetVersion(...)] to edit
	version string
	// (Optional) Positional arguments for the program
	posArgs []string
	// (Optional) Other subcommands
	branches *internal.OrderedMap[*runtimeType]

	// --- Internal
	genOptions    *internal.OrderedMap[internal.Option]
	genOptionAlts map[string]string
	genReqPosArgs []string
}

func newRuntime(program *Program) *runtimeType {
	return &runtimeType{
		data:          program.Data,
		exec:          internal.NewContextExecutable(program.Exec),
		name:          program.Name,
		desc:          program.Desc,
		version:       "",
		posArgs:       []string{},
		branches:      internal.NewOrderedMap[*runtimeType](),
		genOptions:    internal.NewOrderedMap[internal.Option](),
		genOptionAlts: map[string]string{},
	}
}

// Sets the version of the program (SemVer is adviced). Enables built-in `--version` flag.
//
// Designed to be used in main programs, not branches/subcommands.
func (runtime *runtimeType) SetVersion(version string) *runtimeType {
	runtime.version = version
	return runtime
}

// Specifies program positional arguments.
//
// - Include `?` in the argument if it's optional;
//
// - Suffix with `...` in the argument if its variadic;
//
// Prints the arguments in --help menu as provided.
func (runtime *runtimeType) SetPosArgs(args ...string) *runtimeType {
	runtime.posArgs = args
	for _, arg := range args {
		if !strings.Contains(arg, "?") {
			runtime.genReqPosArgs = append(runtime.genReqPosArgs, arg)
		}
	}
	return runtime
}

// Registers a branch/subcommand, which is just another [Program] object.
//
// Prints the command in --help menu.
func (runtime *runtimeType) RegisterCommand(command *runtimeType) *runtimeType {
	runtime.branches.Add(command.name, command)
	return runtime
}

// Run processes options, validates them, and then executes the command.
func (runtime *runtimeType) Run(inputArgs []string) error {
	runtime.exec.Clear()
	if err := runtime.preprocess(); err != nil {
		return err
	}

iterate:
	for i := 0; i < len(inputArgs); i++ {
		arg := inputArgs[i]
		if !strings.HasPrefix(arg, "-") {
			// Branch or positional argument.
			if branch, exists := runtime.branches.Get(arg); exists {
				return branch.Run(inputArgs[i+1:])
			}
			runtime.exec.AddArg(arg)
			continue
		}

		// Handle help and version shortcuts.
		switch arg {
		case "--help":
			runtime.printHelp(os.Stdout)
			return nil
		case "--version":
			runtime.PrintVersion(os.Stdout)
			return nil
		case "--":
			runtime.exec.Args = append(runtime.exec.Args, inputArgs[i+1:]...)
			break iterate
		}

		var err error
		if strings.HasPrefix(arg, "--") {
			err = runtime.processLongOption(arg, &i, inputArgs)
		} else {
			err = runtime.processShortOption(arg, &i, inputArgs)
		}
		if err != nil {
			return err
		}
	}
	lenProv := len(runtime.exec.Args)
	lenReq := len(runtime.genReqPosArgs)
	if lenProv < lenReq {
		var args strings.Builder
		runtime.printArgs(&args)
		return berr.WithContext(
			runtime.name,
			ErrNotEnoughArgs,
			"required length", lenReq,
			"provided length", lenProv,
			"exec args", runtime.exec.Args,
			"arg printer", args.String(),
		)
	}

	if runtime.exec.Function == nil {
		return berr.New(runtime.name, ErrExecIsNil)
	}
	if err := runtime.exec.Run(); err != nil {
		return berr.New(runtime.name, err)
	}
	return nil
}
