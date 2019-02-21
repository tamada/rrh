package fetch

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamadalab/rrh/common"
)

type FetchAllCommand struct{}

func FetchAllCommandFactory() (cli.Command, error) {
	return &FetchAllCommand{}, nil
}

func (fetch *FetchAllCommand) Help() string {
	return `rrh fetch-all [OPTIONS]
OPTIONS
    -r, --remote <REMOTE>   specify the remote name. Default is "origin."`
}

func (fetch *FetchAllCommand) Run(args []string) int {
	var config = common.OpenConfig()

	var fetchOptions, err = fetch.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	if len(fetchOptions.args) != 0 {
		fmt.Println(fetch.Help())
		return 1
	}
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 1
	}
	return fetch.execFetch(db, fetchOptions)
}

func (fetch *FetchAllCommand) printError(errs []error) {
	for _, err := range errs {
		fmt.Println(err.Error())
	}
}

func (fetch *FetchAllCommand) execFetch(db *common.Database, fetchOptions *FetchOptions) int {
	var onError = db.Config.GetValue(common.RrhOnError)

	var fetch2 = FetchCommand{}
	var errorlist = []error{}
	for _, group := range db.Groups {
		var errs = fetch2.FetchGroup(db, group.Name, fetchOptions)
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

func (fetch *FetchAllCommand) parse(args []string) (*FetchOptions, error) {
	var options = FetchOptions{"origin", []string{}}
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

func (fetch *FetchAllCommand) Synopsis() string {
	return "run \"git fetch\" in the all repositories"
}
