package group

import (
	"errors"
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

type removeResult struct {
	groupName  string
	removeDone bool
	err        error
	message    string
}

func (rr *removeResult) Err() error {
	return rr.err
}

func newRemoveResult(name string, doneFlag bool, err error, message string) *removeResult {
	return &removeResult{
		groupName:  name,
		removeDone: doneFlag,
		err:        err,
		message:    message,
	}
}

func executeGroupRemove(c *cobra.Command, args []string, db *rrh.Database) error {
	// verbose, _ := c.PersistentFlags().GetBool("verbose")
	messages := []*removeResult{}
	for _, groupName := range args {
		group := db.FindGroup(groupName)
		if group == nil && !removeOpts.force {
			messages = append(messages, newRemoveResult(groupName, false, errors.New("group not found"), ""))
		}
		if group != nil {
			if !inquiryRemovingGroup(removeOpts.inquiry, groupName) {
				messages = append(messages, newRemoveResult(groupName, false, nil, "not remove by the user request"))
			}
			if err := removeGroupsImpl(c, db, groupName); err != nil {
				messages = append(messages, newRemoveResult(groupName, false, err, ""))
			} else {
				messages = append(messages, newRemoveResult(groupName, true, nil, "remove done"))
			}
		}
	}
	return storeDatabase(c, db, messages)
}

func storeDatabase(c *cobra.Command, db *rrh.Database, messages []*removeResult) error {
	if removeOpts.dryRunFlag {
		for _, result := range messages {
			if result.err != nil {
				c.Printf("%s: %s (remove error)\n", result.groupName, result.err.Error())
			} else if result.removeDone {
				c.Printf("%s: %s (dry-run mode)\n", result.groupName, result.message)
			} else {
				c.Printf("%s: %s\n", result.groupName, result.message)
			}
		}
	} else {
		db.StoreAndClose()
	}
	return common.MergeErrors(convert(messages))
}

func convert(results []*removeResult) []common.Resulter {
	resulters := []common.Resulter{}
	for _, result := range results {
		resulters = append(resulters, result)
	}
	return resulters
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
