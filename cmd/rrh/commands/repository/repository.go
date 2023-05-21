package repository

import (
	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/add"
	"github.com/tamada/rrh/cmd/rrh/commands/utils"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repository",
		Short: "manages repositories",
		RunE: func(c *cobra.Command, args []string) error {
			return utils.PerformRrhCommand(c, args, perform)
		},
	}
	registerSubCommand(cmd)

	return cmd
}

func registerSubCommand(c *cobra.Command) {
	c.AddCommand(add.New())
	c.AddCommand(newListCommand())
	c.AddCommand(newInfoCommand())
	c.AddCommand(newOfCommand())
	c.AddCommand(newUpdateCommand())
}

func perform(c *cobra.Command, args []string, db *rrh.Database) error {
	return performList(c, args, db)
}
