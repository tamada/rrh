package path

import (
	"flag"
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

/*
PathCommand represents a command.
*/
type PathCommand struct {
	options *pathOptions
}

type pathOptions struct {
	partialMatch bool
	showRepoID   bool
	args         []string
}

type pathResult struct {
	id   string
	path string
}

/*
PathCommandFactory returns an instance of the PruneCommand.
*/
func PathCommandFactory() (cli.Command, error) {
	return &PathCommand{}, nil
}

func (options *pathOptions) buildFormatter(results []pathResult) string {
	var maxLength = 0
	for _, r := range results {
		var len = len(r.id)
		if len > maxLength {
			maxLength = len
		}
	}
	return fmt.Sprintf("%%-%ds", maxLength)
}

func (options *pathOptions) showErrorIfNeeded(results []pathResult) int {
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

func (path *PathCommand) perform(db *common.Database) int {
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

func (options *pathOptions) matchEach(id string, arg string) bool {
	if options.partialMatch {
		return strings.Contains(id, arg)
	}
	return id == arg
}

func (options *pathOptions) match(id string) bool {
	for _, arg := range options.args {
		var bool = options.matchEach(id, arg)
		if bool {
			return true
		}
	}
	return len(options.args) == 0
}

func (path *PathCommand) findResult(db *common.Database) []pathResult {
	var results = []pathResult{}
	for _, repo := range db.Repositories {
		if path.options.match(repo.ID) {
			results = append(results, pathResult{id: repo.ID, path: repo.Path})
		}
	}
	return results
}

func (path *PathCommand) buildFlagSet() (*flag.FlagSet, *pathOptions) {
	var options = pathOptions{partialMatch: false, args: []string{}}
	flags := flag.NewFlagSet("path", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(path.Help()) }
	flags.BoolVar(&options.partialMatch, "m", false, "partial match mode")
	flags.BoolVar(&options.partialMatch, "partial-match", false, "partial match mode")
	flags.BoolVar(&options.showRepoID, "r", false, "show path only")
	flags.BoolVar(&options.showRepoID, "show-repository-id", false, "show path only")
	return flags, &options
}

func (path *PathCommand) parse(args []string) error {
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
func (path *PathCommand) Run(args []string) int {
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
func (path *PathCommand) Help() string {
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
func (path *PathCommand) Synopsis() string {
	return "print paths of specified repositories."
}
