package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/add"
	"github.com/tamada/rrh/clone"
	"github.com/tamada/rrh/common"
	"github.com/tamada/rrh/export"
	"github.com/tamada/rrh/fetch"
	"github.com/tamada/rrh/group"
	"github.com/tamada/rrh/list"
	"github.com/tamada/rrh/move"
	"github.com/tamada/rrh/path"
	"github.com/tamada/rrh/prune"
	"github.com/tamada/rrh/remove"
	"github.com/tamada/rrh/status"
)

func buildCommandFactoryMap() map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"add":       add.CommandFactory,
		"clone":     clone.CommandFactory,
		"config":    common.CommandFactory,
		"export":    export.CommandFactory,
		"fetch":     fetch.CommandFactory,
		"fetch-all": fetch.AllCommandFactory,
		"group":     group.CommandFactory,
		"import":    export.ImportCommandFactory,
		"list":      list.CommandFactory,
		"mv":        move.MoveCommandFactory,
		"path":      path.PathCommandFactory,
		"prune":     prune.PruneCommandFactory,
		"rm":        remove.RemoveCommandFactory,
		"status":    status.StatusCommandFactory,
	}
}

func main() {
	c := cli.NewCLI("rrh", common.VERSION)
	c.Name = "rrh"
	c.Args = os.Args[1:]
	c.Autocomplete = true
	c.Commands = buildCommandFactoryMap()

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
