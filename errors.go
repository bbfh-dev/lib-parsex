package parsex

import (
	"errors"
)

var (
	ErrMustbePointer     = errors.New("Program.Data must be a pointer")
	ErrMustPointToStruct = errors.New("Program.Data must point to a struct{}")
	ErrExecIsNil         = errors.New("runtime.Exec function is nil")
	ErrNotEnoughArgs     = errors.New("not enough arguments are provided")
	ErrUnknownOption     = errors.New("unknown option. Refer to --help for usage information")
	ErrOptionNeedsValue  = errors.New("option needs a value. Refer to --help for usage information")
	ErrSettingOption     = errors.New("setting option")
	ErrUnknownCluster    = errors.New("unknown option or cluster. Refer to --help for usage information")
	ErrMistypedCluster   = errors.New("cluster can only contain flags but options are found. Refer to --help for usage information")
)
