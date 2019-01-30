package fetch

import "github.com/mitchellh/cli"

type FetchAllCommand struct{}

func FetchAllCommandFactory() (cli.Command, error) {
	return &FetchAllCommand{}, nil
}

func (fetch *FetchAllCommand) Help() string {
	return `grim fetch [GROUPS...]
ARGUMENTS
    GROUPS    run "git fetch" command on each repository on the group.`
}

func (fetch *FetchAllCommand) Run(args []string) int {
	return 1
}

func (fetch *FetchAllCommand) Synopsis() string {
	return "run \"git fetch\" in the all repositories (not yet)"
}
