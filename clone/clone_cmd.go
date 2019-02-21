package clone

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamadalab/rrh/common"
)

type CloneCommand struct{}

func CloneCommandFactory() (cli.Command, error) {
	return &CloneCommand{}, nil
}

type cloneOptions struct {
	group string
	dest  string
	args  []string
}

/*
Help function shows the help message.
*/
func (clone *CloneCommand) Help() string {
	return `rrh clone [OPTION] <REMOTE_REPOS...>
OPTION
    -g, --group <GROUP>   print managed repositories categoried in the group.
    -d, --dest <DEST>     specify the destination.
ARGUMENTS
    REMOTE_REPOS          repository urls`
}

/*
Synopsis returns the help message of the command.
*/
func (clone *CloneCommand) Synopsis() string {
	return "run \"git clone\""
}

func (clone *CloneCommand) showError(list []error) {
	for _, err := range list {
		fmt.Println(err.Error())
	}
}

func (clone *CloneCommand) Run(args []string) int {
	var config = common.OpenConfig()
	options, err := clone.parse(args, config)
	if err != nil || len(options.args) == 0 {
		fmt.Printf(clone.Help())
		return 1
	}
	db, err := common.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	return clone.perform(db, options)
}

func (clone *CloneCommand) perform(db *common.Database, options *cloneOptions) int {
	var count, list = clone.DoClone(db, options)
	if len(list) != 0 {
		clone.showError(list)
		var onError = db.Config.GetValue(common.RrhOnError)
		if onError == common.Fail || onError == common.FailImmediately {
			return 1
		}
	}
	db.StoreAndClose()
	if count == 0 {
		fmt.Printf("a repository cloned into %s and registered to group %s\n", options.dest, options.group)
	} else {
		fmt.Printf("%d repositories cloned into %s and registered to group %s\n", count, options.dest, options.group)
	}

	return 0
}

func (clone *CloneCommand) parse(args []string, config *common.Config) (*cloneOptions, error) {
	var defaultGroup = config.GetDefaultValue(common.RrhDefaultGroupName)
	var options = cloneOptions{defaultGroup, ".", []string{}}
	flags := flag.NewFlagSet("clone", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(clone.Help()) }
	flags.StringVar(&options.group, "g", defaultGroup, "belonging group")
	flags.StringVar(&options.group, "group", defaultGroup, "belonging group")
	flags.StringVar(&options.dest, "d", ".", "destination")
	flags.StringVar(&options.dest, "dest", ".", "destination")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	options.args = flags.Args()

	return &options, nil
}
