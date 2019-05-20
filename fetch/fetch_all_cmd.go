package fetch

import (
	"fmt"

	"github.com/mitchellh/cli"
	flag "github.com/ogier/pflag"
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
func (fetchAll *AllCommand) Help() string {
	return `rrh fetch-all [OPTIONS]
OPTIONS
    -r, --remote=<REMOTE>   specify the remote name. Default is "origin."`
}

func (fetchAll *AllCommand) validateArguments(args []string) (*options, error) {
	var options, err = fetchAll.parse(args)
	if err == nil {
		if len(options.args) != 0 {
			return nil, fmt.Errorf("fetch-all must be no arguments")
		}
	}
	return options, err
}

/*
Run performs the command.
*/
func (fetchAll *AllCommand) Run(args []string) int {
	var config = common.OpenConfig()

	var options, err = fetchAll.validateArguments(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 1
	}
	return fetchAll.execFetch(db, options)
}

func (fetchAll *AllCommand) printError(errs []error) {
	for _, err := range errs {
		fmt.Println(err.Error())
	}
}

func (fetchAll *AllCommand) printErrors(onError string, errs []error) int {
	if onError == common.Fail || onError == common.Warn {
		fetchAll.printError(errs)
		if onError == common.Fail {
			return 1
		}
	}
	return 0
}

func convertToGroupName(groups []common.Group) []string {
	var result = []string{}
	for _, group := range groups {
		result = append(result, group.Name)
	}
	return result
}

func (fetchAll *AllCommand) execFetch(db *common.Database, options *options) int {
	var onError = db.Config.GetValue(common.RrhOnError)

	var fetch = Command{options}
	var errorlist = []error{}
	var progress = fetch.buildProgress(db, convertToGroupName(db.Groups))
	for _, group := range db.Groups {
		var errs = fetch.FetchGroup(db, group.Name, progress)
		errorlist = append(errorlist, errs...)
		if onError == common.FailImmediately {
			fetchAll.printError(errs)
			return 1
		}
	}
	return fetchAll.printErrors(onError, errorlist)
}

func (fetchAll *AllCommand) parse(args []string) (*options, error) {
	var options = options{"origin", []string{}}
	flags := flag.NewFlagSet("fetch-all", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(fetchAll.Help()) }
	flags.StringVarP(&options.remote, "remote", "r", "origin", "remote name")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	options.args = flags.Args()
	return &options, nil
}

/*
Synopsis returns the help message of the command.
*/
func (fetchAll *AllCommand) Synopsis() string {
	return "run \"git fetch\" in the all repositories."
}
