package add

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamadalab/rrh/common"
)

/*
AddCommand shows the subcommand of rrh.
*/
type AddCommand struct {
}

/*
AddCommandFactory generates the object of AddCommand.
*/
func AddCommandFactory() (cli.Command, error) {
	return &AddCommand{}, nil
}

/*
Help function shows the help message.
*/
func (add *AddCommand) Help() string {
	return `rrh add [OPTION] <REPOSITORY_PATHS...>
OPTION
    -g, --group <GROUP>    add repository to RRH database.

ARGUMENTS
    REPOSITORY_PATHS       the local path list of the git repositories`
}

func (add *AddCommand) showError(errorlist []error, onError string) {
	if len(errorlist) == 0 || onError == common.Ignore {
		return
	}
	for _, item := range errorlist {
		fmt.Println(item.Error())
	}
}

func (add *AddCommand) perform(db *common.Database, args []string, groupName string) int {
	var onError = db.Config.GetValue(common.RrhOnError)

	var errorlist = add.addRepositoriesToGroup(db, args, groupName)

	add.showError(errorlist, onError)

	if onError == common.Fail || onError == common.FailImmediately {
		return 1
	}
	var err2 = db.StoreAndClose()
	if err2 != nil {
		fmt.Println(err2.Error())
	}

	return 0
}

/*
Run function performs the command.
*/
func (add *AddCommand) Run(args []string) int {
	var config = common.OpenConfig()
	var opt, err = add.parse(args, config)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(add.Help())
		return 1
	}
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 1
	}
	return add.perform(db, opt.args, opt.group)
}

type addoptions struct {
	group string
	args  []string
}

func (add *AddCommand) parse(args []string, config *common.Config) (*addoptions, error) {
	var opt = addoptions{}
	flags := flag.NewFlagSet("add", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(add.Help()) }
	flags.StringVar(&opt.group, "g", config.GetValue(common.RrhDefaultGroupName), "target group")
	flags.StringVar(&opt.group, "group", config.GetValue(common.RrhDefaultGroupName), "target group")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	opt.args = flags.Args()
	if opt.group == "" {
		opt.group = config.GetValue(common.RrhDefaultGroupName)
	}

	return &opt, nil
}

/*
Synopsis returns the simple help message of the command.
*/
func (add *AddCommand) Synopsis() string {
	return "add repositories on the local path to RRH"
}
