package prune

import (
	"os"

	"github.com/dustin/go-humanize/english"
	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/utils"
)

type pruneOptions struct {
	dryRunFlag bool
}

var pruneOpts = &pruneOptions{}

func Execute(args []string) error {
	pruneCommand := New()
	pruneCommand.SetArgs(args)
	return pruneCommand.Execute()
}

func New() *cobra.Command {
	pruneCommand := &cobra.Command{
		Use:   "prune",
		Short: "prune unnecessary entries in the rrh database",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return utils.PerformRrhCommand(c, args, performPrune)
		},
	}
	flags := pruneCommand.Flags()
	flags.BoolVarP(&pruneOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")

	return pruneCommand
}

func performPrune(c *cobra.Command, args []string, db *rrh.Database) error {
	dryRunFlag := pruneOpts.dryRunFlag
	err := perform(c, db)
	if err != nil {
		return err
	}
	if !dryRunFlag {
		db.StoreAndClose()
	}
	return nil
}

func dryRunMode() string {
	if pruneOpts.dryRunFlag {
		return " (dry-run mode)"
	}
	return ""
}

func perform(c *cobra.Command, db *rrh.Database) error {
	var repos = removeNotExistRepository(c, db, pruneOpts.dryRunFlag)
	var repos2, groups = db.Prune()
	c.Printf("Pruned %s", english.Plural(len(groups), "group", ""))
	c.Printf(" and %s", english.Plural(len(repos)+len(repos2), "repository", ""))
	c.Printf("%s\n", dryRunMode())
	if pruneOpts.dryRunFlag || utils.IsVerbose(c) {
		printNotExistRepository(c, repos)
		printNoRelationsResults(c, repos2, groups)
	}
	return nil
}

func removeNotExistRepository(c *cobra.Command, db *rrh.Database, dryRunMode bool) []string {
	var removeRepos = []string{}
	for _, repo := range db.Repositories {
		var _, err = os.Stat(repo.Path)
		if os.IsNotExist(err) {
			removeRepos = append(removeRepos, repo.ID)
		}
	}

	result := []string{}
	for _, repo := range removeRepos {
		err := db.DeleteRepository(repo)
		if err == nil {
			result = append(result, repo)
		}
	}
	return result
}

func printNotExistRepository(c *cobra.Command, ids []string) {
	for _, id := range ids {
		c.Printf("%s: repository pruned (not exists)\n", id)
	}
}

func printNoRelationsResults(c *cobra.Command, repos []*rrh.Repository, groups []*rrh.Group) {
	for _, repo := range repos {
		c.Printf("%s: repository pruned (no relations)\n", repo.ID)
	}
	for _, group := range groups {
		c.Printf("%s: group pruned (no relations)\n", group.Name)
	}
}
