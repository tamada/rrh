package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
)

func PruneCommand() *cobra.Command {
	pruneCommand := &cobra.Command{
		Use:   "prune",
		Short: "prune unnecessary entries in the rrh database",
		RunE: func(c *cobra.Command, args []string) error {
			var config = rrh.OpenConfig()
			var db, err = rrh.Open(config)
			if err != nil {
				return err
			}
			dryRunFlag, err := c.Flags().GetBool("dry-run")
			if perform(c, db) || err == nil && !dryRunFlag {
				db.StoreAndClose()
			}
			return nil
		},
	}

	flags := pruneCommand.Flags()
	flags.BoolP("dry-run", "d", false, "dry-run mode")

	return pruneCommand
}

func dryRunMode(c *cobra.Command) string {
	dryRunFlag, err := c.Flags().GetBool("dry-run")
	if err != nil && dryRunFlag {
		return " (dry-run mode)"
	}
	return ""
}

func perform(c *cobra.Command, db *rrh.Database) bool {
	dryRunFlag, err := c.Flags().GetBool("dry-run")
	var count = removeNotExistRepository(c, db, dryRunFlag && err != nil)
	var repos, groups = db.PruneTargets()
	c.Printf("Pruned %d groups, %d repositories%s\n", len(groups), len(repos)+count, dryRunMode(c))
	db.Prune()
	if err != nil && dryRunFlag {
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

func printResults(c *cobra.Command, repos []rrh.Repository, groups []rrh.Group) {
	for _, repo := range repos {
		fmt.Printf("%s: repository pruned (no relations)\n", repo.ID)
	}
	for _, group := range groups {
		fmt.Printf("%s: group pruned (no relations)\n", group.Name)
	}
}
