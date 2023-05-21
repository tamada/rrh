package group

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/utils"
	"github.com/tamada/rrh/decorator"
)

func createGroupOfCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "of <REPOSITORY_ID>",
		Short: "print the group name of the specified repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return utils.PerformRrhCommand(c, args, func(c *cobra.Command, args []string, db *rrh.Database) error {
				if !db.HasRepository(args[0]) {
					return fmt.Errorf("%s: repository not found", args[0])
				}
				var groups = db.FindRelationsOfRepository(args[0])
				return printResults(args[0], groups, db.Config.Decorator)
			})
		},
	}
	return command
}

func printResults(repo string, groups []string, deco decorator.Decorator) error {
	fmt.Printf("%s: ", deco.RepositoryID(repo))
	for index, g := range groups {
		if index != 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%s", deco.GroupName(g))
	}
	fmt.Println()
	return nil
}
