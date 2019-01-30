package list

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamadalab/grim/common"
)

type listOptions struct {
	all         bool
	description bool
	group       bool
	localPath   bool
	remoteURL   bool
	grep        string
	args        []string
}

type ListCommand struct{}

func (list *ListCommand) printRepository(item string, db *common.Database, options *listOptions) {
	var repo = db.FindRepository(item)
	if repo != nil {
		fmt.Printf("\t%s", item)
		if options.localPath || options.all {
			fmt.Printf(",%s", repo.Path)
		}
		if options.remoteURL || options.all {
			fmt.Printf(",%s", repo.URL)
		}
		fmt.Println()
	}
}

func (list *ListCommand) printGroup(group common.Group, db *common.Database, options *listOptions) {
	if options.group || options.all {
		fmt.Printf("%s\n", group.Name)
	}
	if options.description || options.all {
		fmt.Printf("Desc: %s\n", group.Description)
	}
	for _, repository := range group.Items {
		list.printRepository(repository, db, options)
	}
}

func (list *ListCommand) printResult(db *common.Database, options *listOptions) int {
	var resultGroup = db.FindGroups(options.args)
	for _, group := range resultGroup {
		list.printGroup(group, db, options)
	}

	return 0
}

func ListCommandFactory() (cli.Command, error) {
	return &ListCommand{}, nil
}

func (list *ListCommand) Run(args []string) int {
	var options, err = list.parse(args)
	if err != nil {
		fmt.Printf(list.Help())
		return 1
	}
	var config = common.OpenConfig()
	var db = common.Open(config)
	return list.printResult(db, options)
}

func (list *ListCommand) Synopsis() string {
	return "print managed repositories and their groups."
}

func (list *ListCommand) Help() string {
	return `grim list [OPTIONS] [GROUPS...]
OPTIONS
	-a           print all.
	-d           print description of group.
	-g           print group name.
	-p           print local paths.
	-r           print remote urls.
          	     if any options are specified, -a are specified.
    -G <KEYWORD> grep by given keywords (regex).

ARGUMENTS
    GROUPS    print managed repositories categoried in the groups.
	          if no groups are specified, all groups are printed.`
}

func (list *ListCommand) parse(args []string) (*listOptions, error) {
	var options = listOptions{false, false, false, false, false, "", []string{}}
	flags := flag.NewFlagSet("list", flag.ExitOnError)
	flags.Usage = func() { fmt.Printf("%s\n", list.Help()) }
	flags.BoolVar(&options.all, "a", false, "all flag")
	flags.BoolVar(&options.description, "d", false, "description flag")
	flags.BoolVar(&options.group, "g", false, "group flag")
	flags.BoolVar(&options.localPath, "p", false, "local path flag")
	flags.BoolVar(&options.remoteURL, "r", false, "remote url flag")
	flags.StringVar(&options.grep, "G", "*", "regex")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if !(options.all || options.description || options.group || options.localPath || options.remoteURL) {
		options.all = true
	}
	options.args = flags.Args()
	return &options, nil
}
