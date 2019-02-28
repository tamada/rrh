package clone

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

type CloneCommand struct {
	Options *cloneOptions
}

func CloneCommandFactory() (cli.Command, error) {
	return &CloneCommand{&cloneOptions{}}, nil
}

type cloneOptions struct {
	group   string
	dest    string
	verbose bool
}

/*
Help function shows the help message.
*/
func (clone *CloneCommand) Help() string {
	return `rrh clone [OPTIONS] <REMOTE_REPOS...>
OPTIONS
    -g, --group <GROUP>   print managed repositories categoried in the group.
    -d, --dest <DEST>     specify the destination.
    -v, --verbose         verbose mode.
ARGUMENTS
    REMOTE_REPOS          repository urls`
}

/*
Synopsis returns the help message of the command.
*/
func (clone *CloneCommand) Synopsis() string {
	return "run \"git clone\""
}

func (clone *CloneCommand) printIfVerbose(message string) {
	if clone.Options.verbose {
		fmt.Println(message)
	}
}

func (clone *CloneCommand) showError(list []error) {
	for _, err := range list {
		fmt.Println(err.Error())
	}
}

func (clone *CloneCommand) Run(args []string) int {
	var config = common.OpenConfig()
	arguments, err := clone.parse(args, config)
	if err != nil || len(arguments) == 0 {
		fmt.Printf(clone.Help())
		return 1
	}
	db, err := common.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	return clone.perform(db, arguments)
}

func (clone *CloneCommand) perform(db *common.Database, arguments []string) int {
	var count, list = clone.DoClone(db, arguments)
	if len(list) != 0 {
		clone.showError(list)
		var onError = db.Config.GetValue(common.RrhOnError)
		if onError == common.Fail || onError == common.FailImmediately {
			return 1
		}
	}
	db.StoreAndClose()
	if count == 0 {
		fmt.Printf("a repository cloned into %s and registered to group %s\n", clone.Options.dest, clone.Options.group)
	} else {
		fmt.Printf("%d repositories cloned into %s and registered to group %s\n", count, clone.Options.dest, clone.Options.group)
	}

	return 0
}

func (clone *CloneCommand) parse(args []string, config *common.Config) ([]string, error) {
	var defaultGroup = config.GetDefaultValue(common.RrhDefaultGroupName)
	var options = cloneOptions{defaultGroup, ".", false}
	flags := flag.NewFlagSet("clone", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(clone.Help()) }
	flags.StringVar(&options.group, "g", defaultGroup, "belonging group")
	flags.StringVar(&options.group, "group", defaultGroup, "belonging group")
	flags.StringVar(&options.dest, "d", ".", "destination")
	flags.StringVar(&options.dest, "dest", ".", "destination")
	flags.BoolVar(&options.verbose, "v", false, "verbose mode")
	flags.BoolVar(&options.verbose, "verbose", false, "verbose mode")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	clone.Options = &options

	return flags.Args(), nil
}
