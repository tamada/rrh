package list

import (
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamadalab/rrh/common"
)

type ListAllCommand struct{}

func ListAllCommandFactory() (cli.Command, error) {
	return &ListAllCommand{}, nil
}

func (la *ListAllCommand) Run(args []string) int {
	var list = ListCommand{}
	options, err := list.parse(args)
	if err != nil {
		fmt.Printf(la.Help())
		return 1
	}
	var config = common.OpenConfig()
	db, err := common.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var names = []string{}
	for _, group := range db.Groups {
		names = append(names, group.Name)
	}
	options.args = names
	results, err := list.FindResults(db, options)
	list.printResults(results, options)

	return 0
}

/*
Synopsis returns the help message of the command.
*/
func (list *ListAllCommand) Synopsis() string {
	return "print managed repositories and their groups."
}

/*
Help function shows the help message.
*/
func (list *ListAllCommand) Help() string {
	return `rrh list-all [OPTIONS]
OPTIONS
    -a, --all       print all (default).
    -d, --desc      print description of group.
    -p, --path      print local paths.
    -r, --remote    print remote urls.
                    if any options of above are specified, '-a' are specified.

    -c, --csv       print result as csv format.`
}
