package internal

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh"
)

/*
PruneCommand represents a command.
*/
type PruneCommand struct {
	verbose bool
	dryrun  bool
	mc      *rrh.MessageCenter
}

/*
PruneCommandFactory returns an instance of the PruneCommand.
*/
func PruneCommandFactory() (cli.Command, error) {
	return &PruneCommand{verbose: false, dryrun: false, mc: rrh.NewMessageCenter()}, nil
}

func printResults(prune *PruneCommand, repos []rrh.Repository, groups []rrh.Group) {
	prune.mc.Print(os.Stdout, rrh.VERBOSE)
	for _, repo := range repos {
		fmt.Printf("%s: repository pruned (no relations)\n", repo.ID)
	}
	for _, group := range groups {
		fmt.Printf("%s: group pruned (no relations)\n", group.Name)
	}
}

func dryrunMode(mode bool) string {
	if !mode {
		return ""
	}
	return " (dry-run mode)"
}

func (prune *PruneCommand) perform(db *rrh.Database) bool {
	var count = prune.removeNotExistRepository(db)
	var repos, groups = db.PruneTargets()
	fmt.Printf("Pruned %d groups, %d repositories%s\n", len(groups), len(repos)+count, dryrunMode(prune.dryrun))
	db.Prune()
	if prune.verbose || prune.dryrun {
		printResults(prune, repos, groups)
	}

	return true
}

/*
Help function shows the help message.
*/
func (prune *PruneCommand) Help() string {
	return `rrh prune [OPTIONS]
OPTIONS
    -d, --dry-run    dry-run mode.
    -v, --verbose    verbose mode.`
}

func (prune *PruneCommand) parseOptions(args []string) error {
	var flags = flag.NewFlagSet("prune", flag.ContinueOnError)
	flags.Usage = func() { prune.Help() }
	flags.BoolVarP(&prune.dryrun, "dry-run", "d", false, "dry-run mode")
	flags.BoolVarP(&prune.verbose, "verbose", "v", false, "verbose mode")
	return flags.Parse(args)
}

/*
Run performs the command.
*/
func (prune *PruneCommand) Run(args []string) int {
	if err := prune.parseOptions(args); err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var config = rrh.OpenConfig()
	var db, err = rrh.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	if prune.perform(db) || !prune.dryrun {
		db.StoreAndClose()
	}
	return 0
}

func (prune *PruneCommand) deleteNotExistRepository(db *rrh.Database, repo string) int {
	pushMessage(prune, repo, "not exists")
	var err = db.DeleteRepository(repo)
	if err != nil {
		return 0
	}
	return 1
}

func (prune *PruneCommand) removeNotExistRepository(db *rrh.Database) int {
	var removeRepos = []string{}
	for _, repo := range db.Repositories {
		var _, err = os.Stat(repo.Path)
		if os.IsNotExist(err) {
			removeRepos = append(removeRepos, repo.ID)
		}
	}

	var count = 0
	for _, repo := range removeRepos {
		count += prune.deleteNotExistRepository(db, repo)
	}
	return count
}

func pushMessage(prune *PruneCommand, repo, reason string) {
	prune.mc.PushVerbose(fmt.Sprintf("%s: repository pruned (%s)", repo, reason))
}

/*
Synopsis returns the help message of the command.
*/
func (prune *PruneCommand) Synopsis() string {
	return "prune unnecessary repositories and groups."
}
