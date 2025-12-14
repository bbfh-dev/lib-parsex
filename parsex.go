package libparsex

func Run(program *Program, args []string) error {
	if !program.didParse {
		if err := program.Parse(); err != nil {
			return err
		}
	}

	// TODO: Parse args

	return program.EntryPoint()
}
