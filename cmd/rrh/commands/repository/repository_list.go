package repository

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

type listOptions struct {
	entries []string
}

var listOpts = &listOptions{}

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [GROUPs...]",
		Short: "list repositories",
		Args: func(c *cobra.Command, args []string) error {
			return ValidateEntries(listOpts.entries)
		},
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, performList)
		},
	}
	flags := cmd.Flags()
	flags.StringSliceVarP(&listOpts.entries, "entry", "e", []string{"id", "group", "path"}, "print the repository information")
	return cmd
}

func performList(c *cobra.Command, args []string, db *rrh.Database) error {
	entries, err := NewEntries(listOpts.entries)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return executeList(c, db, db.Groups, entries)
	}
	el := common.NewErrorList()
	groups, errs := findGroups(args, db)
	el = el.Append(errs)
	err2 := executeList(c, db, groups, entries)
	el = el.Append(err2)
	return el.NilOrThis()
}

func findGroups(args []string, db *rrh.Database) ([]*rrh.Group, error) {
	results := []*rrh.Group{}
	errs := common.ErrorList{}
	for _, arg := range args {
		group := db.FindGroup(arg)
		if group != nil {
			results = append(results, group)
		} else {
			errs.Append(fmt.Errorf("%s: group not found", arg))
		}
	}
	return results, errs.NilOrThis()
}

func executeList(c *cobra.Command, db *rrh.Database, groups []*rrh.Group, li Entries) error {
	err := common.NewErrorList()
	results := [][]string{}
	for _, group := range groups {
		repos, errs := findRepositories(group, db)
		for _, repo := range repos {
			r := formatRepository(c, group, repo, li)
			results = append(results, r)
		}
		err.Append(errs)
	}
	printAll(c, results, li)
	return err.NilOrThis()
}

func calculateMaxWidth(data [][]string, extractor func(datum []string) string) int {
	max := -1
	for _, datum := range data {
		str := extractor(datum)
		max = rrh.MaxInt(max, len(str))
	}
	return max
}

func printAll(c *cobra.Command, r [][]string, li Entries) {
	maxWidth := calculateMaxWidth(r, func(rr []string) string { return rr[0] })
	formatter := fmt.Sprintf("%%-%ds", maxWidth)
	for _, result := range r {
		first := 0
		if li.IsGroup() || li.IsId() {
			c.Printf(formatter, result[0])
			first = first + 1
		}
		for i := first; i < len(result); i++ {
			c.Printf("    %s", result[i])
		}
		c.Println()
	}
	/*
		buf := bytes.NewBuffer([]byte{})
		table := tablewriter.NewWriter(buf)
		table.SetBorder(false)
		table.SetHeaderLine(false)
		table.AppendBulk(r)
		table.Render()
		c.Println(buf.String())
	*/
}

func findRepositories(group *rrh.Group, db *rrh.Database) ([]*rrh.Repository, common.ErrorList) {
	errs := common.NewErrorList()
	results := []*rrh.Repository{}
	relations := db.FindRelationsOfGroup(group.Name)
	for _, rel := range relations {
		repo := db.FindRepository(rel)
		if repo != nil {
			results = append(results, repo)
		} else {
			errs.Append(fmt.Errorf("%s: repository not found", rel))
		}
	}
	return results, errs
}

func formatRepository(c *cobra.Command, group *rrh.Group, repo *rrh.Repository, li Entries) []string {
	results := []string{formatRepositoryName(group, repo, li)}
	if li.IsPath() {
		results = append(results, repo.Path)
	}
	if li.IsDesc() {
		results = append(results, repo.Description)
	}
	return results
}

func formatRepositoryName(group *rrh.Group, repo *rrh.Repository, li Entries) string {
	if li.IsGroup() && li.IsId() {
		return fmt.Sprintf("%s/%s", group.Name, repo.ID)
	} else if !li.IsGroup() && li.IsId() {
		return repo.ID
	} else if li.IsCount() && !li.IsId() {
		return group.Name
	}
	return ""
}
