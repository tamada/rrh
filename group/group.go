package group

import "github.com/mitchellh/cli"

type GroupCommand struct{}

func GroupCommandFactory() (cli.Command, error) {
	return &GroupCommand{}, nil
}

func (group *GroupCommand) Help() string {
	return `group`
}

func (group *GroupCommand) Run(args []string) int {
	return 1
}

func (group *GroupCommand) Synopsis() string {
	return "print groups. (not yet)"
}
