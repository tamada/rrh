package list

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list groups and repositories",
		Args:  cobra.ArbitraryArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, Perform)
		},
	}
	flags := cmd.Flags()
	flags.StringSliceVarP(&listOpts.entries, "entry", "e", []string{"group", "count", "id", "path", "summary"}, "specifies the printing entries.\navailables: all, group, note, count, id, desc, path, remote, and summary")
	flags.StringVarP(&listOpts.format, "format", "f", "default", "specifies the output format. availables: csv, default, and json")
	flags.BoolVarP(&listOpts.noAbbrev, "no-abbrev", "a", false, "no abbrev mode")
	flags.BoolVarP(&listOpts.header, "no-header", "H", false, "print without header")
	return cmd
}

func validateOpts(c *cobra.Command) error {
	if err := ValidateEntries(listOpts.entries); err != nil {
		return err
	}
	if err := validateFormat(listOpts.format); err != nil {
		return err
	}
	return nil
}

func Perform(c *cobra.Command, args []string, db *rrh.Database) error {
	if err := validateOpts(c); err != nil {
		return err
	}
	c.SilenceUsage = true
	return performImpl(c, args, db)
}

func performImpl(c *cobra.Command, args []string, db *rrh.Database) error {
	results, err := FindResults(db, args)
	if err != nil {
		return err
	}
	le, err := newListEntry(listOpts.entries)
	if err != nil {
		return err
	}
	formatter, err := newFormatter(listOpts.format, listOpts.header)
	if err != nil {
		return err
	}
	err = formatter.Format(c.OutOrStdout(), results, le, listOpts.noAbbrev)
	if err != nil {
		return err
	}
	return nil
}

var listOpts = &listOptions{}

type listOptions struct {
	format   string
	entries  []string
	noAbbrev bool
	header   bool
}

/*
Repo represents the result for showing of repositories.
*/
type Repo struct {
	Name    string        `json:"id"`
	Path    string        `json:"path"`
	Desc    string        `json:"desc"`
	Remotes []*rrh.Remote `json:"remote"`
}

/*
Result represents the result for showing.
*/
type Result struct {
	GroupName string  `json:"group"`
	Note      string  `json:"note"`
	Abbrev    bool    `json:"-"`
	Repos     []*Repo `json:"repositories"`
}

func findList(db *rrh.Database, groupName string) (*Result, error) {
	var repos = []*Repo{}
	var group = db.FindGroup(groupName)
	if group == nil {
		return nil, fmt.Errorf("%s: group not found", groupName)
	}
	for _, relation := range db.Relations {
		if relation.GroupName == groupName {
			var repo = db.FindRepository(relation.RepositoryID)
			if repo == nil {
				return nil, fmt.Errorf("%s: repository not found", relation.RepositoryID)
			}
			repos = append(repos, &Repo{Name: repo.ID, Path: repo.Path, Desc: repo.Description, Remotes: repo.Remotes})
		}
	}

	return &Result{GroupName: group.Name, Note: group.Description, Abbrev: group.OmitList, Repos: repos}, nil
}

func findAllGroupNames(db *rrh.Database) []string {
	var names = []string{}
	for _, group := range db.Groups {
		names = append(names, group.Name)
	}
	return names
}

/*
FindResults returns the result list of list command.
*/
func FindResults(db *rrh.Database, args []string) ([]*Result, error) {
	groups := findGroupNames(db, args)
	results := []*Result{}
	for _, group := range groups {
		var list, err = findList(db, group)
		if err != nil {
			return nil, err
		}
		results = append(results, list)
	}
	return results, nil
}

func findGroupNames(db *rrh.Database, args []string) []string {
	if len(args) == 0 {
		return findAllGroupNames(db)
	}
	return args
}
