package clone

import (
	"fmt"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/common"
)

/*
Command represents a command.
*/
type Command struct {
	options *options
}

/*
CommandFactory returns an instance of the CloneCommand.
*/
func CommandFactory() (cli.Command, error) {
	return &Command{&options{}}, nil
}

type options struct {
	group   string
	dest    string
	verbose bool
}

/*
Help function shows the help message.
*/
func (clone *Command) Help() string {
	return `rrh clone [OPTIONS] <REMOTE_REPOS...>
OPTIONS
    -g, --group <GROUP>   print managed repositories categorized in the group.
    -d, --dest <DEST>     specify the destination.
    -v, --verbose         verbose mode.
ARGUMENTS
    REMOTE_REPOS          repository urls`
}

/*
Synopsis returns the help message of the command.
*/
func (clone *Command) Synopsis() string {
	return "run \"git clone\" and register it to a group."
}

func (clone *Command) printIfVerbose(message string) {
	if clone.options.verbose {
		fmt.Println(message)
	}
}

func (options *options) showError(list []error) {
	for _, err := range list {
		fmt.Println(err.Error())
	}
}

/*
Run performs the command.
*/
func (clone *Command) Run(args []string) int {
	var config = common.OpenConfig()
	arguments, err := clone.parse(args, config)
	if err != nil || len(arguments) == 0 {
		fmt.Printf(clone.Help())
		return 1
	}
	db, err := common.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 2
	}
	return clone.perform(db, arguments)
}

func (clone *Command) perform(db *common.Database, arguments []string) int {
	var count, list = clone.DoClone(db, arguments)
	if len(list) != 0 {
		clone.options.showError(list)
		var onError = db.Config.GetValue(common.RrhOnError)
		if onError == common.Fail || onError == common.FailImmediately {
			return 1
		}
	}
	db.StoreAndClose()
	printResult(count, clone.options.dest, clone.options.group)
	return 0
}

func printResult(count int, dest string, group string) {
	switch count {
	case 0:
		fmt.Println("no repositories cloned")
	case 1:
		fmt.Printf("a repository cloned into %s and registered to group %s\n", dest, group)
	default:
		fmt.Printf("%d repositories cloned into %s and registered to group %s\n", count, dest, group)
	}
}

func (clone *Command) buildFlagSets(config *common.Config) (*flag.FlagSet, *options) {
	var defaultGroup = config.GetValue(common.RrhDefaultGroupName)
	var destination = config.GetValue(common.RrhCloneDestination)
	var options = options{defaultGroup, ".", false}
	flags := flag.NewFlagSet("clone", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(clone.Help()) }
	flags.StringVarP(&options.group, "group", "g", defaultGroup, "belonging group")
	flags.StringVarP(&options.dest, "dest", "d", destination, "destination")
	flags.BoolVarP(&options.verbose, "verbose", "v", false, "verbose mode")
	return flags, &options
}

func (clone *Command) parse(args []string, config *common.Config) ([]string, error) {
	var flags, options = clone.buildFlagSets(config)
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	clone.options = options

	return flags.Args(), nil
}
