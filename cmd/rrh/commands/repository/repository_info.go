package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dustin/go-humanize/english"
	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/utils"
	"github.com/tamada/rrh/common"
)

func newInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "show repository information",
		Args:  validateArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return utils.PerformRrhCommand(c, args, performInfo)
		},
	}
	flags := cmd.Flags()
	flags.StringSliceVarP(&listOpts.entries, "entry", "e", []string{"id", "group", "path"}, "print the repository information")
	return cmd
}

func validateArgs(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("requires one and more arguments")
	}
	_, err := NewEntries(listOpts.entries)
	if err != nil {
		return err
	}
	return nil
}

func performInfo(c *cobra.Command, args []string, db *rrh.Database) error {
	el := common.NewErrorList()
	e, _ := NewEntries(listOpts.entries)
	for _, arg := range args {
		repo := db.FindRepository(arg)
		if repo == nil {
			el = el.Append(fmt.Errorf("%s: repository not found", arg))
		} else {
			err := printRepository(c, repo, e, db)
			el = el.Append(err)
		}
	}
	return el.NilOrThis()
}

func findColoredGroup(db *rrh.Database, repoID string) []string {
	deco := db.Config.Decorator
	groups := db.FindRelationsOfRepository(repoID)
	results := []string{}
	for _, group := range groups {
		results = append(results, deco.GroupName(group))
	}
	return results
}

func printRepository(c *cobra.Command, repo *rrh.Repository, e Entries, db *rrh.Database) error {
	deco := db.Config.Decorator
	if e.IsId() {
		c.Printf("Repository Id: %s\n", deco.RepositoryID(repo.ID))
	}
	if e.IsGroup() {
		groups := findColoredGroup(db, repo.ID)
		if len(groups) > 0 {
			c.Printf("%s: %s\n", english.PluralWord(len(groups), "Group", "Groups"), strings.Join(groups, ", "))
		}
	}
	if e.IsDesc() {
		c.Printf("Description:   %s\n", repo.Description)
	}
	if e.IsPath() {
		c.Printf("Path: %s\n", repo.Path)
	}
	if e.IsRemotes() {
		if len(repo.Remotes) > 0 {
			c.Printf("%s:\n", english.PluralWord(len(repo.Remotes), "Remote", "Remotes"))
			for _, remote := range repo.Remotes {
				c.Printf("    %s: %s\n", remote.Name, remote.URL)
			}
		}
	}
	return nil
}
