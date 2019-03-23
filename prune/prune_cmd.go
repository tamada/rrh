package prune

import (
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
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

func (prune *PruneCommand) perform(db *common.Database) bool {
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
	var config = common.OpenConfig()
	var db, err = common.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	if prune.perform(db) {
		db.StoreAndClose()
	}
	return 0
}

/*
Synopsis returns the help message of the command.
*/
func (prune *PruneCommand) Synopsis() string {
	return "prune unnecessary repositories and groups."
}
