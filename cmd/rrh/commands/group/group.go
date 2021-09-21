package group

import (
	"github.com/spf13/cobra"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

func New() *cobra.Command {
	groupCommand := &cobra.Command{
		Use:   "group [subcommand]",
		Short: "manage the groups for the rrh database",
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, listGroups)
		},
	}
	registerGroupCommands(groupCommand)
	return groupCommand
}

func registerGroupCommands(c *cobra.Command) {
	c.AddCommand(createGroupAddCommand())
	c.AddCommand(createGroupInfoCommand())
	c.AddCommand(createGroupOfCommand())
	c.AddCommand(createGroupListCommand())
	c.AddCommand(createGroupUpdateCommand())
	c.AddCommand(createGroupRemoveCommand())
}
