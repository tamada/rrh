package prune

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/common"
)

type pruneOptions struct {
	dryRunFlag bool
}

var pruneOpts = &pruneOptions{}

func New() *cobra.Command {
	pruneCommand := &cobra.Command{
		Use:   "prune",
		Short: "prune unnecessary entries in the rrh database",
		RunE: func(c *cobra.Command, args []string) error {
			return common.PerformRrhCommand(c, args, func(c *cobra.Command, args []string, db *rrh.Database) error {
				dryRunFlag := pruneOpts.dryRunFlag
				if perform(c, db) || !dryRunFlag {
					db.StoreAndClose()
				}
				return nil
			})
		},
	}

	flags := pruneCommand.Flags()
	flags.BoolVarP(&pruneOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")

	return pruneCommand
}

func dryRunMode(c *cobra.Command) string {
	if pruneOpts.dryRunFlag {
		return " (dry-run mode)"
	}
	return ""
}

func perform(c *cobra.Command, db *rrh.Database) bool {
	var count = removeNotExistRepository(c, db, pruneOpts.dryRunFlag)
	var repos, groups = db.Prune()
	c.Printf("Pruned %d groups, %d repositories%s\n", len(groups), len(repos)+count, dryRunMode(c))
	if pruneOpts.dryRunFlag {
		printResults(c, repos, groups)
	}

	return true
}

func deleteNotExistRepository(c *cobra.Command, db *rrh.Database, repo string) int {
	rrh.PrintIfVerbose(c, fmt.Sprintf("%s: repository removed, not exists in path", repo))
	var err = db.DeleteRepository(repo)
	if err != nil {
		return 0
	}
	return 1
}

func removeNotExistRepository(c *cobra.Command, db *rrh.Database, dryRunMode bool) int {
	var removeRepos = []string{}
	for _, repo := range db.Repositories {
		var _, err = os.Stat(repo.Path)
		if os.IsNotExist(err) {
			removeRepos = append(removeRepos, repo.ID)
		}
	}

	var count = 0
	for _, repo := range removeRepos {
		count += deleteNotExistRepository(c, db, repo)
	}
	return count
}

func printResults(c *cobra.Command, repos []*rrh.Repository, groups []*rrh.Group) {
	for _, repo := range repos {
		fmt.Printf("%s: repository pruned (no relations)\n", repo.ID)
	}
	for _, group := range groups {
		fmt.Printf("%s: group pruned (no relations)\n", group.Name)
	}
}
