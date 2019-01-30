package remove

import "github.com/mitchellh/cli"

type RemoveCommand struct{}

func RemoveCommandFactory() (cli.Command, error) {
	return &RemoveCommand{}, nil
}

func (rm *RemoveCommand) Run(args []string) int {
	return 1
}

func (rm *RemoveCommand) Help() string {
	return `grim rm [OPTION] <REPOSITORY_ID/GROUP_NAME...>
OPTION
	-f               force remove.
	-i               inquiry.
	-r               remove group and repositories in the group.

ARGUMENTS
	REPOSITORY_ID    repository name.
	GROUP_NAME       group name.`
}

func (rm *RemoveCommand) Synopsis() string {
	return "remove given repository from database. (not yet)"
}
