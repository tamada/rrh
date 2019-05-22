package fetch

import (
	"fmt"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
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
    -r, --remote <REMOTE>   specify the remote name. Default is "origin."`
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
	return handleError(fetchAll.execFetch(db, options), config.GetValue(common.RrhOnError))
}

func convertToGroupNames(groups []common.Group) []string {
	var result = []string{}
	for _, group := range groups {
		result = append(result, group.Name)
	}
	return result
}

func (fetchAll *AllCommand) execFetch(db *common.Database, options *options) []error {
	var onError = db.Config.GetValue(common.RrhOnError)
	var errorlist = []error{}
	var fetch = Command{options}
	var relations = fetch.FindTargets(db, convertToGroupNames(db.Groups))
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
