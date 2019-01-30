package clone

import "github.com/mitchellh/cli"

type BatchCloneCommand struct{}

func BatchCloneCommandFactory() (cli.Command, error) {
	return &BatchCloneCommand{}, nil
}

func (clone *BatchCloneCommand) Help() string {
	return `grim clone <GRIM_DATABASE_JSON>
ARGUMENTS
	GRIM_DATABASE_JSON   `
}

func (clone *BatchCloneCommand) Run(args []string) int {
	return 1
}

func (clone *BatchCloneCommand) Synopsis() string {
	return "run \"git clone\" (not yet)"
}
