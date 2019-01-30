package clone

import "github.com/mitchellh/cli"

type CloneCommand struct{}

func CloneCommandFactory() (cli.Command, error) {
	return &CloneCommand{}, nil
}

func (clone *CloneCommand) Help() string {
	return `grim clone [OPTION] <REOMOTE_REPO>
OPTION
	-g <GROUP>    print managed repositories categoried in the group.
	-d <DEST>     specify the destination.
ARGUMENTS
	REMOTE_REPO`
}

func (clone *CloneCommand) Run(args []string) int {
	return 1
}

func (clone *CloneCommand) Synopsis() string {
	return "run \"git clone\" (not yet)"
}
