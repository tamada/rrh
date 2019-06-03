package internal

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/lib"
)

/*
PruneCommand represents a command.
*/
type PruneCommand struct {
}

/*
PruneCommandFactory returns an instance of the PruneCommand.
*/
func PruneCommandFactory() (cli.Command, error) {
	return &PruneCommand{}, nil
}

func (prune *PruneCommand) perform(db *lib.Database) bool {
	var count = prune.removeNotExistRepository(db)
	var gCount, rCount = db.Prune()
	fmt.Printf("Pruned %d groups, %d repositories\n", gCount, rCount+count)
	return true
}

/*
Help function shows the help message.
*/
func (prune *PruneCommand) Help() string {
	return `rrh prune`
}

/*
Run performs the command.
*/
func (prune *PruneCommand) Run(args []string) int {
	var config = lib.OpenConfig()
	var db, err = lib.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	if prune.perform(db) {
		db.StoreAndClose()
	}
	return 0
}

func (prune *PruneCommand) removeNotExistRepository(db *lib.Database) int {
	var removeRepos = []string{}
	for _, repo := range db.Repositories {
		var _, err = os.Stat(repo.Path)
		if os.IsNotExist(err) {
			removeRepos = append(removeRepos, repo.ID)
		}
	}

	var count = 0
	for _, repo := range removeRepos {
		var err = db.DeleteRepository(repo)
		if err == nil {
			count++
		}
	}
	return count
}

/*
Synopsis returns the help message of the command.
*/
func (prune *PruneCommand) Synopsis() string {
	return "prune unnecessary repositories and groups."
}
