package group

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

func createGroupOfCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "of <REPOSITORY_ID>",
		Short: "print the group name of the specified repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, func(c *cobra.Command, args []string, db *rrh.Database) error {
				if !db.HasRepository(args[0]) {
					return fmt.Errorf("%s: repository not found", args[0])
				}
				var groups = db.FindRelationsOfRepository(args[0])
				fmt.Printf("%s,%v\n", args[0], groups)
				return nil
			})
		},
	}
	return command
}
