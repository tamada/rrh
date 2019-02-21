package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/tamadalab/rrh/add"
	"github.com/tamadalab/rrh/clone"
	"github.com/tamadalab/rrh/common"
	"github.com/tamadalab/rrh/export"
	"github.com/tamadalab/rrh/fetch"
	"github.com/tamadalab/rrh/group"
	"github.com/tamadalab/rrh/list"
	"github.com/tamadalab/rrh/prune"
	"github.com/tamadalab/rrh/remove"
	"github.com/tamadalab/rrh/status"
)

func main() {
	c := cli.NewCLI("rrh", "1.0.0")
	c.Args = os.Args[1:]
	c.Autocomplete = true
	c.Commands = map[string]cli.CommandFactory{
		"add":       add.AddCommandFactory,
		"clone":     clone.CloneCommandFactory,
		"config":    common.ConfigCommandFactory,
		"export":    export.ExportCommandFactory,
		"fetch":     fetch.FetchCommandFactory,
		"fetch-all": fetch.FetchAllCommandFactory,
		"group":     group.GroupCommandFactory,
		"list":      list.ListCommandFactory,
		"list-all":  list.ListAllCommandFactory,
		"prune":     prune.PruneCommandFactory,
		"rm":        remove.RemoveCommandFactory,
		"status":    status.StatusCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
