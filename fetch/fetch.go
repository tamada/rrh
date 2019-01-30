package fetch

import "github.com/mitchellh/cli"

type FetchCommand struct{}

func FetchCommandFactory() (cli.Command, error) {
	return &FetchCommand{}, nil
}

func (fetch *FetchCommand) Help() string {
	return `grim fetch [GROUPS...]
ARGUMENTS
    GROUPS    run "git fetch" command on each repository on the group.`
}

func (fetch *FetchCommand) Run(args []string) int {
	return 1
}

func (fetch *FetchCommand) Synopsis() string {
	return "run \"git fetch\" (not yet)"
}
