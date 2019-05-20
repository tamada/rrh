package fetch

import (
	"fmt"

	"github.com/mitchellh/cli"
	flag "github.com/ogier/pflag"
	"github.com/tamada/rrh/common"
)

/*
Command represents a command.
*/
type Command struct {
	options *options
}

/*
CommandFactory returns an instance of command.
*/
func CommandFactory() (cli.Command, error) {
	return &Command{&options{}}, nil
}

/*
Help returns the help message of the command.
*/
func (fetch *Command) Help() string {
	return `rrh fetch [OPTIONS] [GROUPS...]
OPTIONS
    -r, --remot=<REMOTE>   specify the remote name. Default is "origin."
ARGUMENTS
    GROUPS                  run "git fetch" command on each repository on the group.
                            if no value is specified, run on the default group.`
}

/*
Synopsis returns the help message of the command.
*/
func (fetch *Command) Synopsis() string {
	return "run \"git fetch\" on the given groups."
}

/*
Run performs the command.
*/
func (fetch *Command) Run(args []string) int {
	var options, err = fetch.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var config = common.OpenConfig()
	if len(options.args) == 0 {
		options.args = []string{config.GetValue(common.RrhDefaultGroupName)}
	}
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 1
	}
	return fetch.perform(db)
}

func (fetch *Command) buildProgress(db *common.Database, groupNames []string) *Progress {
	var progress = Progress{total: 0, current: 0}
	for _, name := range groupNames {
		var count = db.ContainsCount(name)
		progress.total += count
	}
	return &progress
}

func (fetch *Command) perform(db *common.Database) int {
	var errorFlag = 0
	var onError = db.Config.GetValue(common.RrhOnError)
	var progress = fetch.buildProgress(db, fetch.options.args)
	fmt.Printf("before progress: %s\n", progress)
	for _, groupName := range fetch.options.args {
		var list = fetch.FetchGroup(db, groupName, progress)
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

type options struct {
	remote string
	// key      string
	// userName string
	// password string
	args []string
}

func (fetch *Command) parse(args []string) (*options, error) {
	var options = options{"origin", []string{}}
	flags := flag.NewFlagSet("fetch", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(fetch.Help()) }
	flags.StringVarP(&options.remote, "remote", "r", "origin", "remote name")
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
