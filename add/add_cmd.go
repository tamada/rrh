package add

import (
	"fmt"

	"github.com/mitchellh/cli"
	flag "github.com/ogier/pflag"
	"github.com/tamada/rrh/common"
)

/*
Command shows the subcommand of rrh.
*/
type Command struct {
	options *options
}

/*
CommandFactory generates the object of AddCommand.
*/
func CommandFactory() (cli.Command, error) {
	return &Command{}, nil
}

/*
Help function shows the help message.
*/
func (add *Command) Help() string {
	return `rrh add [OPTIONS] <REPOSITORY_PATHS...>
OPTIONS
    -g, --group=<GROUP>        add repository to RRH database.
    -r, --repository-id=<ID>   specified repository id of the given repository path.
                               Specifying this option fails with multiple arguments.
ARGUMENTS
    REPOSITORY_PATHS           the local path list of the git repositories.`
}

func (add *Command) showError(errorlist []error, onError string) {
	if len(errorlist) == 0 || onError == common.Ignore {
		return
	}
	for _, item := range errorlist {
		fmt.Println(item.Error())
	}
}

func (add *Command) perform(db *common.Database, opt *options) int {
	var onError = db.Config.GetValue(common.RrhOnError)

	var errorlist = add.AddRepositoriesToGroup(db, opt)

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
func (add *Command) Run(args []string) int {
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
		return 2
	}
	return add.perform(db, opt)
}

type options struct {
	group  string
	repoID string
	args   []string
}

func (add *Command) buildFlagSet(config *common.Config) (*flag.FlagSet, *options) {
	var opt = options{}
	var defaultGroup = config.GetValue(common.RrhDefaultGroupName)
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(add.Help()) }
	flags.StringVarP(&opt.group, "group", "g", defaultGroup, "target group")
	flags.StringVarP(&opt.repoID, "repository-id", "r", "", "specifying repository id")
	return flags, &opt
}

func (add *Command) parse(args []string, config *common.Config) (*options, error) {
	var flags, opt = add.buildFlagSet(config)
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	opt.args = flags.Args()
	add.options = opt

	return opt, nil
}

/*
Synopsis returns the simple help message of the command.
*/
func (add *Command) Synopsis() string {
	return "add repositories on the local path to RRH."
}
