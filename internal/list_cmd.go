package internal

import (
	"fmt"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/lib"
)

type listOptions struct {
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
ListCommand represents a command.
*/
type ListCommand struct {
	options *listOptions
}

/*
ListCommandFactory returns an instance of the ListCommand.
*/
func ListCommandFactory() (cli.Command, error) {
	return &ListCommand{&listOptions{}}, nil
}

func (options *listOptions) isChecked(target bool) bool {
	return target || options.all
}

func (options *listOptions) printResultAsCsv(result Result, repo Repo, remote *lib.Remote) {
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

func (options *listOptions) printRepoAsCsv(repo Repo, result Result) {
	if len(repo.Remotes) > 0 && (options.remoteURL || options.all) {
		for _, remote := range repo.Remotes {
			options.printResultAsCsv(result, repo, &remote)
		}
	} else {
		options.printResultAsCsv(result, repo, nil)
	}
}

func (options *listOptions) printResultsAsCsv(results []Result) int {
	for _, result := range results {
		for _, repo := range result.Repos {
			options.printRepoAsCsv(repo, result)
		}
	}
	return 0
}

func findMaxLengthOfRepositoryName(repos []Repo) int {
	var max = len("Description")
	for _, repo := range repos {
		var len = len(repo.Name)
		if len > max {
			max = len
		}
	}
	return max
}

/*
printColoriezdRepositoryID prints the repository name in color.
Coloring escape sequence breaks the printf position arrangement.
Therefore, we arranges the positions by spacing behind the colored repository name.
*/
func printColoriezdRepositoryID(repoName string, length int, config *lib.Config) {
	var formatter = fmt.Sprintf("    %%s%%%ds", length-len(repoName))
	fmt.Printf(formatter, config.Color.ColorizedRepositoryID(repoName), "")
}

func (options *listOptions) printRepo(repo Repo, result Result, maxLength int, config *lib.Config) {
	printColoriezdRepositoryID(repo.Name, maxLength, config)
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

func (options *listOptions) isPrintSimple(result Result) bool {
	return !options.noOmit && result.OmitList && len(options.args) == 0
}

func printGroupName(result Result, config *lib.Config) int {
	if len(result.Repos) == 1 {
		fmt.Printf("%s (1 repository)\n", config.Color.ColorizedGroupName(result.GroupName))
	} else {
		fmt.Printf("%s (%d repositories)\n", config.Color.ColorizedGroupName(result.GroupName), len(result.Repos))
	}
	return len(result.Repos)
}

func (options *listOptions) printResult(result Result, config *lib.Config) int {
	var repoCount = printGroupName(result, config)
	if !options.isPrintSimple(result) {
		if options.description || options.all {
			fmt.Printf("    Description  %s", result.Description)
			fmt.Println()
		}
		var maxLength = findMaxLengthOfRepositoryName(result.Repos)
		for _, repo := range result.Repos {
			options.printRepo(repo, result, maxLength, config)
		}
	}
	return repoCount
}

func (options *listOptions) printSimpleResult(repo Repo, result Result) {
	if options.repoNameOnly {
		fmt.Println(repo.Name)
	} else if options.groupRepoName {
		fmt.Printf("%s/%s\n", result.GroupName, repo.Name)
	}
}

func (options *listOptions) printSimpleResults(results []Result) int {
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

func (options *listOptions) printResults(results []Result, config *lib.Config) int {
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

func (list *ListCommand) findAndPrintResult(db *lib.Database) int {
	results, err := list.FindResults(db)
	if err != nil {
		fmt.Println(err.Error())
		return 3
	}
	return list.options.printResults(results, db.Config)
}

func (list *ListCommand) printError(err error, printHelpFlag bool, statusCode int) int {
	fmt.Println(err.Error())
	if printHelpFlag {
		fmt.Println(list.Help())
	}
	return statusCode
}

/*
Run performs the command.
*/
func (list *ListCommand) Run(args []string) int {
	var _, err = list.parse(args)
	if err != nil {
		return list.printError(err, true, 1)
	}
	var config = lib.OpenConfig()
	db, err := lib.Open(config)
	if err != nil {
		return list.printError(err, false, 2)
	}
	return list.findAndPrintResult(db)
}

/*
Synopsis returns the help message of the command.
*/
func (list *ListCommand) Synopsis() string {
	return "print managed repositories and their groups."
}

/*
Help function shows the help message.
*/
func (list *ListCommand) Help() string {
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

func (list *ListCommand) buildFlagSet() (*flag.FlagSet, *listOptions) {
	var options = listOptions{args: []string{}}
	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(list.Help()) }
	flags.BoolVarP(&options.all, "all-entries", "A", false, "show all entries")
	flags.BoolVarP(&options.description, "desc", "d", false, "description flag")
	flags.BoolVarP(&options.localPath, "path", "p", false, "local path flag")
	flags.BoolVarP(&options.remoteURL, "remote", "r", false, "remote url flag")
	flags.BoolVarP(&options.noOmit, "all", "a", false, "no omit repositories")
	flags.BoolVarP(&options.csv, "csv", "c", false, "print as csv format")
	flags.BoolVar(&options.repoNameOnly, "only-repositoryname", false, "show only repository names")
	flags.BoolVar(&options.groupRepoName, "group-repository-form", false, "show group and repository pair form")
	return flags, &options
}

func (list *ListCommand) parse(args []string) (*listOptions, error) {
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

/*
Repo represents the result for showing of repositories.
*/
type Repo struct {
	Name    string
	Path    string
	Remotes []lib.Remote
}

/*
Result represents the result for showing.
*/
type Result struct {
	GroupName   string
	Description string
	OmitList    bool
	Repos       []Repo
}

func (list *ListCommand) findList(db *lib.Database, groupName string) (*Result, error) {
	var repos = []Repo{}
	var group = db.FindGroup(groupName)
	if group == nil {
		return nil, fmt.Errorf("%s: group not found", groupName)
	}
	for _, relation := range db.Relations {
		if relation.GroupName == groupName {
			var repo = db.FindRepository(relation.RepositoryID)
			if repo == nil {
				return nil, fmt.Errorf("%s: repository not found", relation.RepositoryID)
			}
			repos = append(repos, Repo{repo.ID, repo.Path, repo.Remotes})
		}
	}

	return &Result{group.Name, group.Description, group.OmitList, repos}, nil
}

func (list *ListCommand) findAllGroupNames(db *lib.Database) []string {
	var names = []string{}
	for _, group := range db.Groups {
		names = append(names, group.Name)
	}
	return names
}

/*
FindResults returns the result list of list command.
*/
func (list *ListCommand) FindResults(db *lib.Database) ([]Result, error) {
	var groups = list.options.args
	if len(groups) == 0 {
		groups = list.findAllGroupNames(db)
	}
	var results = []Result{}
	for _, group := range groups {
		var list, err = list.findList(db, group)
		if err != nil {
			return nil, err
		}
		results = append(results, *list)
	}
	return results, nil
}
