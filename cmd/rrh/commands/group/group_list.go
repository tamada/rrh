package group

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

type groupEntry int

const (
	name   groupEntry = 1
	count             = 2
	desc              = 4
	repo              = 8
	abbrev            = 16
	all               = name | desc | repo | abbrev
)

type listOptions struct {
	format  string
	entries []string
	header  bool
}

var listOpts = &listOptions{}

func createGroupListCommand() *cobra.Command {
	command := &cobra.Command{
		Use:  "list",
		Args: cobra.NoArgs,
		PreRunE: func(c *cobra.Command, args []string) error {
			return validateListOpts(listOpts)
		},
		Short: "print the group list",
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, listGroups)
		},
	}
	flags := command.Flags()
	flags.StringVarP(&listOpts.format, "format", "f", "table", "specifies the output format. Available values: csv, json, and table.")
	flags.StringSliceVarP(&listOpts.entries, "entry", "e", []string{"name", "count"}, "specifies the printing entries separated with comma. Available vaues: all, name, desc, repo, abbrev, and count")
	flags.BoolVarP(&listOpts.header, "without-header", "H", false, "print without headers")

	return command
}

func validateListOpts(opts *listOptions) error {
	err := common.ValidateGroupEntries(opts.entries)
	if err != nil {
		return err
	}
	return common.ValidateFormatter(opts.format)
}

func listGroups(c *cobra.Command, args []string, db *rrh.Database) error {
	if len(listOpts.entries) == 0 {
		listOpts.entries = []string{"name", "count"}
	}
	entry, err := common.NewGroupEntries(listOpts.entries)
	if err != nil {
		return err
	}
	headers, results, err := groupListResult(entry, db)
	if err != nil {
		return err
	}
	formatter, err := common.NewFormatter(listOpts.format, !listOpts.header)
	if err != nil {
		return err
	}
	return formatter.Print(c, headers, results)
}

func groupListResult(ge common.GroupEntry, db *rrh.Database) (headers []string, values [][]string, err error) {
	results := [][]string{}
	for _, group := range db.Groups {
		resultItems := []string{}
		if ge.IsName() {
			resultItems = append(resultItems, group.Name)
		}
		if ge.IsDesc() {
			resultItems = append(resultItems, group.Description)
		}
		if ge.IsAbbrev() {
			resultItems = append(resultItems, strconv.FormatBool(group.OmitList))
		}
		if ge.IsRepo() {
			list := db.FindRelationsOfGroup(group.Name)
			resultItems = append(resultItems, fmt.Sprintf("%v", list))
		}
		if ge.IsCount() {
			resultItems = append(resultItems, repositoryCount(db, group.Name))
		}
		results = append(results, resultItems)
	}
	return ge.StringArray(), results, nil
}
