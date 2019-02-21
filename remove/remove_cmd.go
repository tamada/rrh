package remove

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

type RemoveCommand struct{}

type removeOptions struct {
	inquiry   bool
	recursive bool
	verbose   bool
	args      []string
}

func RemoveCommandFactory() (cli.Command, error) {
	return &RemoveCommand{}, nil
}

func (rm *RemoveCommand) printIfVerbose(message string, options *removeOptions) {
	if options.verbose {
		fmt.Println(message)
	}
}

func (rm *RemoveCommand) Run(args []string) int {
	var rmOptions, err = rm.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var config = common.OpenConfig()
	var db, err1 = common.Open(config)
	if err1 != nil {
		fmt.Println(err1.Error())
		return 1
	}

	var result = 0
	for _, target := range rmOptions.args {
		var err = rm.executeRemove(db, target, rmOptions)
		if err != nil {
			fmt.Println(err.Error())
			result = 1
		}
	}
	if result == 0 {
		if config.GetValue(common.RrhAutoDeleteGroup) == "yes" {
			db.Prune()
		}
		db.StoreAndClose()
	}

	return result
}

func (rm *RemoveCommand) parse(args []string) (*removeOptions, error) {
	var options = removeOptions{false, false, false, []string{}}
	flags := flag.NewFlagSet("rm", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(rm.Help()) }
	flags.BoolVar(&options.inquiry, "i", false, "inquiry flag")
	flags.BoolVar(&options.verbose, "v", false, "verbose flag")
	flags.BoolVar(&options.recursive, "r", false, "recursive flag")
	flags.BoolVar(&options.inquiry, "inquiry", false, "inquiry flag")
	flags.BoolVar(&options.verbose, "verbose", false, "verbose flag")
	flags.BoolVar(&options.recursive, "recursive", false, "recursive flag")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	options.args = flags.Args()
	return &options, nil
}

/*
Help returns the help message.
*/
func (rm *RemoveCommand) Help() string {
	return `rrh rm [OPTION] <REPO_ID|GROUP_ID|REPO_ID/GROUP_ID...>
OPTION
    -i, --inquiry       inquiry mode.
    -r, --recursive     recursive mode.
    -v, --verbose       verbose mode.

ARGUMENTS
    REPOY_ID            repository name for removing.
    GROUP_ID            group name. if the group contains repositories,
                        remove will fail without '-r' option.
    GROUP_ID/REPO_ID    remove given REPO_ID from GROUP_ID.`
}

/*
Synopsis returns the help message of the command.
*/
func (rm *RemoveCommand) Synopsis() string {
	return "remove given repository from database."
}
