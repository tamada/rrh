package remove

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

type removeOptions struct {
	inquiry   bool
	recursive bool
	verbose   bool
	args      []string
}

/*
RemoveCommand represents a command.
*/
type RemoveCommand struct {
	Options *removeOptions
}

/*
RemoveCommandFactory returns an instance of the RemoveCommand.
*/
func RemoveCommandFactory() (cli.Command, error) {
	return &RemoveCommand{&removeOptions{}}, nil
}

func (options *removeOptions) printIfVerbose(message string) {
	if options.verbose {
		fmt.Println(message)
	}
}

func (rm *RemoveCommand) perform(db *common.Database) int {
	var result = 0
	for _, target := range rm.Options.args {
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
func (rm *RemoveCommand) Run(args []string) int {
	var options, err = rm.parse(args)
	if err != nil {
		return 1
	}
	rm.Options = options
	var config = common.OpenConfig()
	var db, err1 = common.Open(config)
	if err1 != nil {
		fmt.Println(err1.Error())
		return 2
	}
	return rm.perform(db)
}

func (rm *RemoveCommand) buildFlagSet() (*flag.FlagSet, *removeOptions) {
	var options = removeOptions{false, false, false, []string{}}
	flags := flag.NewFlagSet("rm", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(rm.Help()) }
	flags.BoolVar(&options.inquiry, "i", false, "inquiry flag")
	flags.BoolVar(&options.verbose, "v", false, "verbose flag")
	flags.BoolVar(&options.recursive, "r", false, "recursive flag")
	flags.BoolVar(&options.inquiry, "inquiry", false, "inquiry flag")
	flags.BoolVar(&options.verbose, "verbose", false, "verbose flag")
	flags.BoolVar(&options.recursive, "recursive", false, "recursive flag")
	return flags, &options
}

func (rm *RemoveCommand) parse(args []string) (*removeOptions, error) {
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
func (rm *RemoveCommand) Help() string {
	return `rrh rm [OPTIONS] <REPO_ID|GROUP_ID|REPO_ID/GROUP_ID...>
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
func (rm *RemoveCommand) Synopsis() string {
	return "remove given repository from database."
}
