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
	return handleError(fetch.perform(db), config.GetValue(common.RrhOnError))
}

func handleError(errors []error, onError string) int {
	if len(errors) > 0 {
		if onError != common.Ignore {
			printErrors(errors)
		}
		if onError == common.Fail || onError == common.FailImmediately {
			return 5
		}
	}
	return 0
}

func printErrors(errorlist []error) {
	for _, err := range errorlist {
		fmt.Println(err.Error())
	}
}

func (fetch *Command) perform(db *common.Database) []error {
	var errorlist = []error{}
	var onError = db.Config.GetValue(common.RrhOnError)
	var relations = fetch.FindTargets(db, fetch.options.args)
	var progress = Progress{total: len(relations)}

	for _, relation := range relations {
		var err = fetch.FetchRepository(db, &relation, &progress)
		if err != nil {
			if onError == common.FailImmediately {
				return []error{err}
			}
			errorlist = append(errorlist, err)
		}
	}
	return errorlist
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
