package group

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/utils"
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
			return utils.PerformRrhCommand(c, args, updateGroup)
		},
	}
	flags := command.Flags()
	flags.StringVarP(&updateOpts.name, "name", "", "", "specify the new group name")
	flags.StringVarP(&updateOpts.desc, "note", "", "", "specify the new note of the group")
	flags.StringVarP(&updateOpts.abbrev, "abbrev", "", "false", "specify the new abbrev flag of the group, the given string must be true or false")
	flags.BoolVarP(&updateOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")
	return command
}

func findGroup(origName string, db *rrh.Database) (*rrh.Group, error) {
	group := *db.FindGroup(origName)
	if updateOpts.name != "" {
		group.Name = updateOpts.name
	}
	if updateOpts.desc != "" {
		group.Description = updateOpts.desc
	}
	if updateOpts.abbrev != "" {
		abbrevFlag, err := strconv.ParseBool(updateOpts.abbrev)
		if err != nil {
			return nil, err
		}
		group.OmitList = abbrevFlag
	}
	return &group, nil
}

func updateGroup(c *cobra.Command, args []string, db *rrh.Database) error {
	group, err := findGroup(args[0], db)
	if err != nil {
		return err
	}
	if !db.UpdateGroup(args[0], group) {
		return fmt.Errorf("%s: update failed", args[0])
	}
	printResult(c, args[0], group)
	if updateOpts.dryRunFlag {
		c.Println("(dry-run mode)")
		return nil
	}
	c.Println()
	return db.StoreAndClose()
}

func printResult(c *cobra.Command, origName string, group *rrh.Group) {
	c.Printf("update(%s) = %s (Note: %s, Abbrev: %v)\n", origName, group.Name, group.Description, group.OmitList)
}
