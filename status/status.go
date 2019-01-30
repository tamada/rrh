package status

import (
	"github.com/mitchellh/cli"
)

type StatusCommand struct {
}

func StatusCommandFactory() (cli.Command, error) {
	return &StatusCommand{}, nil
}

func (status *StatusCommand) Help() string {
	return `grim status [GROUPS...]
ARGUMENTS
    GROUPS        target groups.`
}

func (status *StatusCommand) Run(args []string) int {
	return 1
}

func (status *StatusCommand) Synopsis() string {
	return "show git status of repositories. (not yet)"
}
