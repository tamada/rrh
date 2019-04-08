package list

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

type options struct {
	all           bool
	description   bool
	localPath     bool
	remoteURL     bool
	csv           bool
	noOmit        bool
	repoNameOnly  bool
	groupRepoName bool
	args          []string
}

/*
Command represents a command.
*/
type Command struct {
	options *options
}

/*
CommandFactory returns an instance of the ListCommand.
*/
func CommandFactory() (cli.Command, error) {
	return &Command{&options{}}, nil
}

func (options *options) isChecked(target bool) bool {
	return target || options.all
}

func (options *options) printResultAsCsv(result Result, repo Repo, remote *common.Remote) {
	fmt.Printf("%s", result.GroupName)
	if options.isChecked(options.description) {
		fmt.Printf(",%s", result.Description)
	}
	fmt.Printf(",%s", repo.Name)
	if options.isChecked(options.localPath) {
		fmt.Printf(",%s", repo.Path)
	}
	if remote != nil && options.isChecked(options.remoteURL) {
		fmt.Printf(",%s,%s", remote.Name, remote.URL)
	}
	fmt.Println()
}

func (options *options) printRepoAsCsv(repo Repo, result Result) {
	if len(repo.Remotes) > 0 && (options.remoteURL || options.all) {
		for _, remote := range repo.Remotes {
			options.printResultAsCsv(result, repo, &remote)
		}
	} else {
		options.printResultAsCsv(result, repo, nil)
	}
}

func (options *options) printResultsAsCsv(results []Result) int {
	for _, result := range results {
		for _, repo := range result.Repos {
			options.printRepoAsCsv(repo, result)
		}
	}
	return 0
}

/*
generateFormatString returns the formatter for `Printf` to arrange the length of repository names.
*/
func (options *options) generateFormatString(repos []Repo) string {
	var max = len("Description")
	for _, repo := range repos {
		var len = len(repo.Name)
		if len > max {
			max = len
		}
	}
	return fmt.Sprintf("    %%-%ds", max)
}

func (options *options) printRepo(repo Repo, result Result, formatString string, config *common.Config) {
	fmt.Printf(formatString, common.RepositoryID(repo.Name, config))
	if options.localPath || options.all {
		fmt.Printf("  %s", repo.Path)
	}
	if options.remoteURL || options.all {
		for _, remote := range repo.Remotes {
			fmt.Println()
			fmt.Printf("        %s  %s", remote.Name, remote.URL)
		}
	}
	fmt.Println()
}

func (options *options) isPrintSimple(result Result) bool {
	return !options.noOmit && result.OmitList && len(options.args) == 0
}

func (options *options) printGroupName(result Result, config *common.Config) int {
	if len(result.Repos) == 1 {
		fmt.Printf("%s (1 repository)\n", common.GroupName(result.GroupName, config))
	} else {
		fmt.Printf("%s (%d repositories)\n", common.GroupName(result.GroupName, config), len(result.Repos))
	}
	return len(result.Repos)
}

func (options *options) printResult(result Result, config *common.Config) int {
	var repoCount = options.printGroupName(result, config)
	if !options.isPrintSimple(result) {
		if options.description || options.all {
			fmt.Printf("    Description  %s", result.Description)
			fmt.Println()
		}
		var formatString = options.generateFormatString(result.Repos)
		for _, repo := range result.Repos {
			options.printRepo(repo, result, formatString, config)
		}
	}
	return repoCount
}

func (options *options) printSimpleResult(repo Repo, result Result) {
	if options.repoNameOnly {
		fmt.Println(repo.Name)
	} else if options.groupRepoName {
		fmt.Printf("%s/%s\n", result.GroupName, repo.Name)
	}
}

func (options *options) printSimpleResults(results []Result) int {
	for _, result := range results {
		for _, repo := range result.Repos {
			options.printSimpleResult(repo, result)
		}
	}
	return 0
}

func printGroupAndRepoCount(groupCount int, repoCount int) {
	var groupLabel = "groups"
	var repoLabel = "repositories"
	if groupCount == 1 {
		groupLabel = "group"
	}
	if repoCount == 1 {
		repoLabel = "repository"
	}
	fmt.Printf("%d %s, %d %s\n", groupCount, groupLabel, repoCount, repoLabel)
}

func (options *options) printResults(results []Result, config *common.Config) int {
	if options.csv {
		return options.printResultsAsCsv(results)
	} else if options.repoNameOnly || options.groupRepoName {
		return options.printSimpleResults(results)
	}
	var repoCount int
	for _, result := range results {
		repoCount += options.printResult(result, config)
	}
	printGroupAndRepoCount(len(results), repoCount)
	return 0
}

/*
Run performs the command.
*/
func (list *Command) Run(args []string) int {
	var _, err = list.parse(args)
	if err != nil {
		fmt.Printf(list.Help())
		return 1
	}
	var config = common.OpenConfig()
	db, err := common.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 2
	}
	results, err := list.FindResults(db)
	if err != nil {
		fmt.Println(err.Error())
		return 3
	}
	return list.options.printResults(results, config)
}

/*
Synopsis returns the help message of the command.
*/
func (list *Command) Synopsis() string {
	return "print managed repositories and their groups."
}

/*
Help function shows the help message.
*/
func (list *Command) Help() string {
	return `rrh list [OPTIONS] [GROUPS...]
OPTIONS
    -d, --desc          print description of group.
    -p, --path          print local paths (default).
    -r, --remote        print remote urls.
    -A, --all-entries   print all entries of each repository.

    -a, --all           print all repositories, no omit repositories.
    -c, --csv           print result as csv format.
ARGUMENTS
    GROUPS    print managed repositories categorized in the groups.
              if no groups are specified, all groups are printed.`
}

func (list *Command) buildFlagSet() (*flag.FlagSet, *options) {
	var options = options{args: []string{}}
	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(list.Help()) }
	flags.BoolVar(&options.all, "A", false, "show all entries")
	flags.BoolVar(&options.all, "all-entries", false, "show all entries")
	flags.BoolVar(&options.description, "d", false, "description flag")
	flags.BoolVar(&options.description, "desc", false, "description flag")
	flags.BoolVar(&options.localPath, "p", false, "local path flag")
	flags.BoolVar(&options.localPath, "path", false, "local path flag")
	flags.BoolVar(&options.remoteURL, "r", false, "remote url flag")
	flags.BoolVar(&options.remoteURL, "remote", false, "remote url flag")
	flags.BoolVar(&options.noOmit, "a", false, "no omit repositories")
	flags.BoolVar(&options.noOmit, "all", false, "no omit repositories")
	flags.BoolVar(&options.csv, "c", false, "print as csv format")
	flags.BoolVar(&options.csv, "csv", false, "print as csv format")
	flags.BoolVar(&options.repoNameOnly, "only-repositoryname", false, "show only repository names")
	flags.BoolVar(&options.groupRepoName, "group-repository-form", false, "show group and repository pair form")
	return flags, &options
}

func (list *Command) parse(args []string) (*options, error) {
	var flags, options = list.buildFlagSet()

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if !(options.all || options.description || options.localPath || options.remoteURL) {
		options.localPath = true
	}
	options.args = flags.Args()
	list.options = options
	return options, nil
}
