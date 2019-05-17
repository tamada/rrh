package remove

import (
	"fmt"

	"github.com/mitchellh/cli"
	flag "github.com/ogier/pflag"
	"github.com/tamada/rrh/common"
)

type options struct {
	inquiry   bool
	recursive bool
	verbose   bool
	args      []string
}

/*
Command represents a command.
*/
type Command struct {
	options *options
}

/*
CommandFactory returns an instance of the RemoveCommand.
*/
func CommandFactory() (cli.Command, error) {
	return &Command{&options{}}, nil
}

func (options *options) printIfVerbose(message string) {
	if options.verbose {
		fmt.Println(message)
	}
}

func (rm *Command) perform(db *common.Database) int {
	var result = 0
	for _, target := range rm.options.args {
		var err = rm.executeRemove(db, target)
		if err != nil {
			fmt.Println(err.Error())
			result = 3
		}
	}
	if result == 0 {
		if db.Config.IsSet(common.RrhAutoDeleteGroup) {
			db.Prune()
		}
		db.StoreAndClose()
	}
	return result
}

/*
Run performs the command.
*/
func (rm *Command) Run(args []string) int {
	var options, err = rm.parse(args)
	if err != nil {
		return 1
	}
	rm.options = options
	var config = common.OpenConfig()
	var db, err1 = common.Open(config)
	if err1 != nil {
		fmt.Println(err1.Error())
		return 2
	}
	return rm.perform(db)
}

func (rm *Command) buildFlagSet() (*flag.FlagSet, *options) {
	var options = options{false, false, false, []string{}}
	flags := flag.NewFlagSet("rm", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(rm.Help()) }
	flags.BoolVarP(&options.inquiry, "inquiry", "i", false, "inquiry flag")
	flags.BoolVarP(&options.verbose, "verbose", "v", false, "verbose flag")
	flags.BoolVarP(&options.recursive, "recursive", "r", false, "recursive flag")
	return flags, &options
}

func (rm *Command) parse(args []string) (*options, error) {
	var flags, options = rm.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	options.args = flags.Args()
	return options, nil
}

/*
Help returns the help message.
*/
func (rm *Command) Help() string {
	return `rrh rm [OPTIONS] <REPO_ID|GROUP_ID|GROUP_ID/REPO_ID...>
OPTIONS
    -i, --inquiry       inquiry mode.
    -r, --recursive     recursive mode.
    -v, --verbose       verbose mode.

ARGUMENTS
    REPOY_ID            repository name for removing.
    GROUP_ID            group name. if the group contains repositories,
                        remove will fail without '-r' option.
    GROUP_ID/REPO_ID    remove the relation between the given REPO_ID and GROUP_ID.`
}

/*
Synopsis returns the help message of the command.
*/
func (rm *Command) Synopsis() string {
	return "remove given repository from database."
}
