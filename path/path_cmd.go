package path

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/common"
)

/*
Command represents a command.
*/
type Command struct {
	options *options
}

type options struct {
	partialMatch bool
	showRepoID   bool
	args         []string
}

type result struct {
	id   string
	path string
}

/*
CommandFactory returns an instance of the PruneCommand.
*/
func CommandFactory() (cli.Command, error) {
	return &Command{}, nil
}

func (options *options) buildFormatter(results []result) string {
	var maxLength = 0
	for _, r := range results {
		var len = len(r.id)
		if len > maxLength {
			maxLength = len
		}
	}
	return fmt.Sprintf("%%-%ds", maxLength)
}

func (options *options) showErrorIfNeeded(results []result) int {
	if len(results) != 0 {
		return 0
	}
	var message = "found"
	if options.partialMatch {
		message = "match"
	}
	fmt.Printf("%s: repository not %s", message, options.args[0])
	return 5
}

func (path *Command) perform(db *common.Database) int {
	var results = path.findResult(db)
	var formatter = path.options.buildFormatter(results)
	for _, r := range results {
		if path.options.showRepoID {
			fmt.Printf(formatter+" %s\n", r.id, r.path)
		} else {
			fmt.Println(r.path)
		}
	}
	return path.options.showErrorIfNeeded(results)
}

func (options *options) matchEach(id string, arg string) bool {
	if options.partialMatch {
		return strings.Contains(id, arg)
	}
	return id == arg
}

func (options *options) match(id string) bool {
	for _, arg := range options.args {
		var bool = options.matchEach(id, arg)
		if bool {
			return true
		}
	}
	return len(options.args) == 0
}

func (path *Command) findResult(db *common.Database) []result {
	var results = []result{}
	for _, repo := range db.Repositories {
		if path.options.match(repo.ID) {
			results = append(results, result{id: repo.ID, path: repo.Path})
		}
	}
	return results
}

func (path *Command) buildFlagSet() (*flag.FlagSet, *options) {
	var options = options{partialMatch: false, args: []string{}}
	flags := flag.NewFlagSet("path", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(path.Help()) }
	flags.BoolVarP(&options.partialMatch, "partial-match", "m", false, "partial match mode")
	flags.BoolVarP(&options.showRepoID, "show-repository-id", "r", false, "show path only")
	return flags, &options
}

func (path *Command) parse(args []string) error {
	var flags, options = path.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return err
	}
	options.args = flags.Args()
	path.options = options
	return nil
}

/*
Run performs the command.
*/
func (path *Command) Run(args []string) int {
	fmt.Printf("path subcommand is deprecated.")
	var err = path.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var config = common.OpenConfig()
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	return path.perform(db)
}

/*
Help function shows the help message.
*/
func (path *Command) Help() string {
	return `rrh path [OPTIONS] <REPOSITORIES...>
OPTIONS
    -m, --partial-match        treats the arguments as the patterns.
    -r, --show-repository-id   show repository name.
ARGUMENTS
    REPOSITORIES               repository ids.`
}

/*
Synopsis returns the help message of the command.
*/
func (path *Command) Synopsis() string {
	return "print paths of specified repositories."
}
