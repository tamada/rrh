package group

import (
	"fmt"
	"strconv"
	"strings"

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

	return command
}

func validateListOpts(opts *listOptions) error {
	err := validateEntries(opts)
	if err != nil {
		return err
	}
	return common.ValidateFormatter(opts.format)
}

func validateEntries(opts *listOptions) error {
	availables := []string{"all", "name", "count", "desc", "repo", "abbrev"}
	messages := []string{}
	for _, entry := range opts.entries {
		entry = strings.ToLower(entry)
		if !rrh.FindIn(entry, availables) {
			messages = append(messages, entry)
		}
	}
	if len(messages) == 0 {
		return nil
	} else if len(messages) == 1 {
		return fmt.Errorf("%s: unknown entry. available values: %s", messages[0], strings.Join(availables, ","))
	}
	return fmt.Errorf("%s: unknown entries. available values: %s", strings.Join(messages, ","), strings.Join(availables, ", "))
}

func convertToFlags(entries []string) (groupEntry, error) {
	var result groupEntry = 0
	for _, entry := range entries {
		switch strings.ToLower(entry) {
		case "all":
			result |= all
		case "name":
			result |= name
		case "repo":
			result |= repo
		case "desc":
			result |= desc
		case "abbrev":
			result |= abbrev
		case "count":
			result |= count
		default:
			return 0, fmt.Errorf("%s: invalid entry, availables are: all, name, desc, repo, abbrev, count", entry)
		}
	}
	return result, nil
}

func listGroups(c *cobra.Command, args []string, db *rrh.Database) error {
	if len(listOpts.entries) == 0 {
		listOpts.entries = []string{"name", "count"}
	}
	entry, err := convertToFlags(listOpts.entries)
	if err != nil {
		return err
	}
	headers, results, err := groupListResult(entry, db)
	if err != nil {
		return err
	}
	formatter, err := common.NewFormatter(listOpts.format)
	if err != nil {
		return err
	}
	return formatter.Print(c, headers, results)
}

func createHeaders(printTarget groupEntry) []string {
	headers := []string{}
	if printTarget&name == name {
		headers = append(headers, "name")
	}
	if printTarget&desc == desc {
		headers = append(headers, "description")
	}
	if printTarget&abbrev == abbrev {
		headers = append(headers, "abbrev")
	}
	if printTarget&repo == repo {
		headers = append(headers, "repositories")
	}
	if printTarget&count == count {
		headers = append(headers, "repository count")
	}
	return headers
}

func groupListResult(printTarget groupEntry, db *rrh.Database) (headers []string, values [][]string, err error) {
	results := [][]string{}
	for _, group := range db.Groups {
		resultItems := []string{}
		if printTarget&name == name {
			resultItems = append(resultItems, group.Name)
		}
		if printTarget&desc == desc {
			resultItems = append(resultItems, group.Description)
		}
		if printTarget&abbrev == abbrev {
			resultItems = append(resultItems, strconv.FormatBool(group.OmitList))
		}
		if printTarget&repo == repo {
			list := db.FindRelationsOfGroup(group.Name)
			resultItems = append(resultItems, fmt.Sprintf("%v", list))
		}
		if printTarget&count == count {
			resultItems = append(resultItems, repositoryCount(db, group.Name))
		}
		results = append(results, resultItems)
	}
	return createHeaders(printTarget), results, nil
}

func printGroupList(c *cobra.Command, printTarget groupEntry, db *rrh.Database) error {
	for _, group := range db.Groups {
		resultItems := []string{}
		if printTarget&name == name {
			resultItems = append(resultItems, group.Name)
		}
		if printTarget&desc == desc {
			resultItems = append(resultItems, group.Description)
		}
		if printTarget&abbrev == abbrev {
			resultItems = append(resultItems, strconv.FormatBool(group.OmitList))
		}
		if printTarget&repo == repo {
			list := db.FindRelationsOfGroup(group.Name)
			resultItems = append(resultItems, fmt.Sprintf("%v", list))
		}
		if printTarget&count == count {
			resultItems = append(resultItems, repositoryCount(db, group.Name))
		}
		c.Println(strings.Join(resultItems, ","))
	}
	return nil
}
