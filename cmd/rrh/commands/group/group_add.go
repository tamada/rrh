package group

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/utils"
)

type createOptions struct {
	desc       string
	abbrevFlag string
	dryRunFlag bool
}

var createOpts = &createOptions{}

func createGroupAddCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "add <GROUP_NAME>",
		Short: "add groups to the rrh database",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return utils.PerformRrhCommand(c, args, func(c *cobra.Command, args []string, db *rrh.Database) error {
				if err := createAndAddGroup(c, args[0], db); err != nil {
					return err
				}
				if !createOpts.dryRunFlag {
					return db.StoreAndClose()
				}
				c.Printf("create group Name: %s, Description: %s, Abbrev: %v\n", args[0], createOpts.desc, createOpts.abbrevFlag)
				return nil
			})
		},
	}
	flags := command.Flags()
	flags.StringVarP(&createOpts.desc, "note", "n", "", "specifies the note")
	flags.StringVarP(&createOpts.abbrevFlag, "abbrev", "a", "false", "set the group as the abbrev")
	flags.BoolVarP(&createOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")
	return command
}

func createAndAddGroup(c *cobra.Command, groupName string, db *rrh.Database) error {
	abbrevFlag, err := strconv.ParseBool(createOpts.abbrevFlag)
	if err != nil {
		return err
	}
	_, err = db.CreateGroup(groupName, createOpts.desc, abbrevFlag)
	return err
}
