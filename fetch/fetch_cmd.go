package fetch

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

/*
FetchCommand represents a command.
*/
type FetchCommand struct {
	options *fetchOptions
}

/*
FetchCommandFactory returns an instance of command.
*/
func FetchCommandFactory() (cli.Command, error) {
	return &FetchCommand{&fetchOptions{}}, nil
}

/*
Help returns the help message of the command.
*/
func (fetch *FetchCommand) Help() string {
	return `rrh fetch [OPTIONS] [GROUPS...]
OPTIONS
    -r, --remote <REMOTE>   specify the remote name. Default is "origin."
ARGUMENTS
    GROUPS                  run "git fetch" command on each repository on the group.
                            if no value is specified, run on the default group.`
}

/*
Synopsis returns the help message of the command.
*/
func (fetch *FetchCommand) Synopsis() string {
	return "run \"git fetch\" on the given groups."
}

/*
Run performs the command.
*/
func (fetch *FetchCommand) Run(args []string) int {
	var fetchOptions, err = fetch.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var config = common.OpenConfig()
	if len(fetchOptions.args) == 0 {
		fetchOptions.args = []string{config.GetValue(common.RrhDefaultGroupName)}
	}
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 1
	}
	return fetch.perform(db)
}

func (fetch *FetchCommand) perform(db *common.Database) int {
	var errorFlag = 0
	var onError = db.Config.GetValue(common.RrhOnError)
	for _, groupName := range fetch.options.args {
		var list = fetch.FetchGroup(db, groupName)
		for _, err := range list {
			if onError != common.Ignore {
				fmt.Println(err.Error())
				errorFlag = 1
			}
		}
	}
	if onError == common.Warn {
		return 0
	}
	return errorFlag
}

type fetchOptions struct {
	remote string
	// key      string
	// userName string
	// password string
	args []string
}

func (fetch *FetchCommand) parse(args []string) (*fetchOptions, error) {
	var options = fetchOptions{"origin", []string{}}
	flags := flag.NewFlagSet("fetch", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(fetch.Help()) }
	flags.StringVar(&options.remote, "r", "origin", "remote name")
	flags.StringVar(&options.remote, "remote", "origin", "remote name")
	// flags.StringVar(&options.key, "k", "", "private key path")
	// flags.StringVar(&options.userName, "u", "", "user name")
	// flags.StringVar(&options.password, "p", "", "password")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	options.args = flags.Args()
	fetch.options = &options
	return &options, nil
}
