package fetch

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

/*
AllCommand represents a command.
*/
type AllCommand struct{}

/*
AllCommandFactory returns an instance of the FetchAllCommand.
*/
func AllCommandFactory() (cli.Command, error) {
	return &AllCommand{}, nil
}

/*
Help returns the help message.
*/
func (fetch *AllCommand) Help() string {
	return `rrh fetch-all [OPTIONS]
OPTIONS
    -r, --remote <REMOTE>   specify the remote name. Default is "origin."`
}

/*
Run performs the command.
*/
func (fetch *AllCommand) Run(args []string) int {
	var config = common.OpenConfig()

	var options, err = fetch.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	if len(options.args) != 0 {
		fmt.Println(fetch.Help())
		return 1
	}
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 1
	}
	return fetch.execFetch(db, options)
}

func (fetch *AllCommand) printError(errs []error) {
	for _, err := range errs {
		fmt.Println(err.Error())
	}
}

func (fetch *AllCommand) execFetch(db *common.Database, options *options) int {
	var onError = db.Config.GetValue(common.RrhOnError)

	var fetch2 = Command{}
	fetch2.options = options
	var errorlist = []error{}
	for _, group := range db.Groups {
		var errs = fetch2.FetchGroup(db, group.Name)
		errorlist = append(errorlist, errs...)
		if onError == common.FailImmediately {
			fetch.printError(errorlist)
			return 1
		}
	}
	if onError == common.Fail || onError == common.Warn {
		for _, err := range errorlist {
			fmt.Println(err.Error())
		}
		if onError == common.Fail {
			return 1
		}
	}

	return 0
}

func (fetch *AllCommand) parse(args []string) (*options, error) {
	var options = options{"origin", []string{}}
	flags := flag.NewFlagSet("fetch-all", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(fetch.Help()) }
	flags.StringVar(&options.remote, "r", "origin", "remote name")
	// flags.StringVar(&options.key, "k", "", "private key")
	// flags.StringVar(&options.userName, "u", "", "user name")
	// flags.StringVar(&options.password, "p", "", "password")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	options.args = flags.Args()
	return &options, nil
}

/*
Synopsis returns the help message of the command.
*/
func (fetch *AllCommand) Synopsis() string {
	return "run \"git fetch\" in the all repositories."
}
