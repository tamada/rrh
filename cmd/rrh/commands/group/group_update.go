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
	group := db.FindGroup(origName)
	if group == nil {
		return nil, fmt.Errorf("%s: group not found", origName)
	}
	newGroup := &rrh.Group{}
	*newGroup = *group
	if updateOpts.name != "" {
		newGroup.Name = updateOpts.name
	}
	if updateOpts.desc != "" {
		newGroup.Description = updateOpts.desc
	}
	if updateOpts.abbrev != "" {
		abbrevFlag, err := strconv.ParseBool(updateOpts.abbrev)
		if err != nil {
			return nil, err
		}
		newGroup.OmitList = abbrevFlag
	}
	return newGroup, nil
}

func updateGroup(c *cobra.Command, args []string, db *rrh.Database) error {
	group, err := findGroup(args[0], db)
	if err != nil {
		return err
	}
	if err := db.UpdateGroup(args[0], group); err != nil {
		return err
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
	c.Printf("update(%s) = %s (Note: %s, Abbrev: %v)", origName, group.Name, group.Description, group.OmitList)
}
