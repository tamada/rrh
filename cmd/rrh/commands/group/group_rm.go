package group

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

type removeOptions struct {
	inquiry    bool
	force      bool
	dryRunFlag bool
}

var removeOpts = &removeOptions{}

func createGroupRemoveCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "rm <GROUP NAMEs...>",
		Args: cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, executeGroupRemove)
		},
	}
	flags := c.Flags()
	flags.BoolVarP(&removeOpts.inquiry, "inquiry", "i", false, "inquiry mode")
	flags.BoolVarP(&removeOpts.force, "force", "f", false, "force remove")
	flags.BoolVarP(&removeOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")

	return c
}

func executeGroupRemove(c *cobra.Command, args []string, db *rrh.Database) error {
	// verbose, _ := c.PersistentFlags().GetBool("verbose")
	for _, groupName := range args {
		group := db.FindGroup(groupName)
		if group == nil && !removeOpts.force {
			return fmt.Errorf("%s: group not found", groupName)
		}
		if group != nil || !inquiryRemovingGroup(removeOpts.inquiry, groupName) {
			return nil
		}
		if err := removeGroupsImpl(c, db, groupName); err != nil {
			return err
		}
	}
	return nil
}

func removeGroupsImpl(c *cobra.Command, db *rrh.Database, groupName string) error {
	if removeOpts.force {
		db.ForceDeleteGroup(groupName)
		c.Printf("%s: group removed", groupName)
	} else if db.ContainsCount(groupName) == 0 {
		db.DeleteGroup(groupName)
		c.Printf("%s: group removed", groupName)
	} else {
		return fmt.Errorf("%s: cannot remove group. the group has relations", groupName)
	}
	return nil
}

func inquiryRemovingGroup(inquiryFlag bool, groupName string) bool {
	if !inquiryFlag {
		return true
	}
	return rrh.IsInputYes(fmt.Sprintf("%s: remove group? [yN]", groupName))
}
