package internal

import (
	"fmt"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh"
)

/*
FetchAllCommand represents a command.
*/
type FetchAllCommand struct{}

/*
FetchAllCommandFactory returns an instance of the FetchAllCommand.
*/
func FetchAllCommandFactory() (cli.Command, error) {
	return &FetchAllCommand{}, nil
}

/*
Help returns the help message.
*/
func (fetchAll *FetchAllCommand) Help() string {
	return `rrh fetch-all [OPTIONS]
OPTIONS
    -r, --remote <REMOTE>   specify the remote name. Default is "origin."`
}

func (fetchAll *FetchAllCommand) validateArguments(args []string) (*fetchOptions, error) {
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
func (fetchAll *FetchAllCommand) Run(args []string) int {
	var config = rrh.OpenConfig()

	var options, err = fetchAll.validateArguments(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var db, err2 = rrh.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 1
	}
	return printErrors(config, fetchAll.execFetch(db, options))
}

func convertToGroupNames(groups []rrh.Group) []string {
	var result = []string{}
	for _, group := range groups {
		result = append(result, group.Name)
	}
	return result
}

func (fetchAll *FetchAllCommand) execFetch(db *rrh.Database, options *fetchOptions) []error {
	var onError = db.Config.GetValue(rrh.OnError)
	var errorlist = []error{}
	var fetch = FetchCommand{options}
	var relations = rrh.FindTargets(db, convertToGroupNames(db.Groups))
	var progress = NewProgress(len(relations))
	for _, relation := range relations {
		var err = fetch.FetchRepository(db, &relation, progress)
		if err != nil {
			if onError == rrh.FailImmediately {
				return []error{err}
			}
			errorlist = append(errorlist, err)
		}
	}
	return errorlist
}

func (fetchAll *FetchAllCommand) parse(args []string) (*fetchOptions, error) {
	var options = fetchOptions{remote: "origin", args: []string{}}
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
func (fetchAll *FetchAllCommand) Synopsis() string {
	return "run \"git fetch\" in the all repositories."
}
