package group

import (
	"fmt"

	"github.com/dustin/go-humanize/english"
	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

func createGroupInfoCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "info <GROUP>",
		Short: "print the information of the specified group",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, printGroupInfos)
		},
	}
	return command
}

func printGroupInfos(c *cobra.Command, args []string, db *rrh.Database) error {
	el := common.NewErrorList()
	for _, groupName := range args {
		err := printGroupInfo(c, groupName, db)
		el = el.Append(err)
	}
	return el.NilOrThis()
}

func printGroupInfo(c *cobra.Command, groupName string, db *rrh.Database) error {
	group := db.FindGroup(groupName)
	if group == nil {
		return fmt.Errorf("%s: group not found", groupName)
	}
	count := db.ContainsCount(group.Name)
	c.Printf("%s: %s (%s, abbrev: %v)\n", group.Name, group.Description, english.Plural(count, "repository", ""), group.OmitList)
	return nil
}
