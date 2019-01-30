package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/tamadalab/grim/add"
	"github.com/tamadalab/grim/clone"
	"github.com/tamadalab/grim/common"
	"github.com/tamadalab/grim/export"
	"github.com/tamadalab/grim/fetch"
	"github.com/tamadalab/grim/group"
	"github.com/tamadalab/grim/list"
	"github.com/tamadalab/grim/remove"
	"github.com/tamadalab/grim/status"
)

func main() {
	c := cli.NewCLI("grim", "1.0.0")
	c.Args = os.Args[1:]
	c.Autocomplete = true
	c.Commands = map[string]cli.CommandFactory{
		"add":         add.AddCommandFactory,
		"batch_clone": clone.BatchCloneCommandFactory,
		"clone":       clone.CloneCommandFactory,
		"config":      common.ConfigCommandFactory,
		"export":      export.ExportCommandFactory,
		"group":       group.GroupCommandFactory,
		"list":        list.ListCommandFactory,
		"rm":          remove.RemoveCommandFactory,
		"status":      status.StatusCommandFactory,
		"fetch":       fetch.FetchCommandFactory,
		"fetch_all":   fetch.FetchAllCommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
