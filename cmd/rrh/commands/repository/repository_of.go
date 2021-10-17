package repository

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/list"
	"github.com/tamada/rrh/cmd/rrh/commands/utils"
	"github.com/tamada/rrh/common"
)

func newOfCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "of",
		Short: "show repository list of the given groups. This is an alias of \"rrh list\" command",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return utils.PerformRrhCommand(c, args, performOf)
		},
	}
	return cmd
}

func performOf(c *cobra.Command, args []string, db *rrh.Database) error {
	el := common.NewErrorList()
	for _, arg := range args {
		if !db.HasGroup(arg) {
			el.Append(fmt.Errorf("%s: group not found", arg))
		}
	}
	if el.IsErr() {
		return el
	}
	return list.Perform(c, args, db)
}
