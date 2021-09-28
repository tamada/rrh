package group

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

type updateOptions struct {
	name       string
	desc       string
	abbrev     string
	dryRunFlag bool
}

var updateOpts = &updateOptions{}

func createGroupUpdateCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "update <GROUP>",
		Short: "update the information of the specified group",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, updateGroup)
		},
	}
	flags := command.Flags()
	flags.StringVarP(&updateOpts.name, "name", "n", "", "change the group name to the given string")
	flags.StringVarP(&updateOpts.desc, "desc", "d", "", "change the description to the given string")
	flags.StringVarP(&updateOpts.abbrev, "abbrev", "a", "false", "change abbrev flag of the group, the given string must be true or false")
	flags.BoolVarP(&updateOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")
	return command
}

func updateGroup(c *cobra.Command, args []string, db *rrh.Database) error {
	abbrevFlag, err := strconv.ParseBool(updateOpts.abbrev)
	if err != nil {
		return err
	}
	group := &rrh.Group{Name: updateOpts.name, Description: updateOpts.desc, OmitList: abbrevFlag}
	if !db.UpdateGroup(args[0], group) {
		return fmt.Errorf("%s: update failed", args[0])
	}
	c.Printf("update %s -> %v", args[0], group)
	if updateOpts.dryRunFlag {
		c.Println("(dry run mode)")
	} else {
		c.Println()
	}
	return nil
}
