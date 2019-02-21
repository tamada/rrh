package prune

import (
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamadalab/rrh/common"
)

type PruneCommand struct {
}

func PruneCommandFactory() (cli.Command, error) {
	return &PruneCommand{}, nil
}

func (prune *PruneCommand) perform(db *common.Database) bool {
	var count = prune.removeNotExistRepository(db)
	var gCount, rCount = db.Prune()
	fmt.Printf("Pruned %d groups, %d repositories\n", gCount, rCount+count)
	return true
}

func (prune *PruneCommand) Help() string {
	return `rrh prune`
}

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
