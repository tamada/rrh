package list

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

type listOptions struct {
	all         bool
	description bool
	localPath   bool
	remoteURL   bool
	csv         bool
	args        []string
}

/*
ListCommand represents a command.
*/
type ListCommand struct {
	Options *listOptions
}

/*
ListCommandFactory returns an instance of the ListCommand.
*/
func ListCommandFactory() (cli.Command, error) {
	return &ListCommand{&listOptions{}}, nil
}

func (list *ListCommand) printResultAsCsv(result ListResult, repo Repo, remote *common.Remote) {
	fmt.Printf("%s", result.GroupName)
	if list.Options.description || list.Options.all {
		fmt.Printf(",%s", result.Description)
	}
	fmt.Printf(",%s", repo.Name)
	if list.Options.localPath || list.Options.all {
		fmt.Printf(",%s", repo.Path)
	}
	if remote != nil && (list.Options.remoteURL || list.Options.all) {
		fmt.Printf(",%s,%s", remote.Name, remote.URL)
	}
	fmt.Println()
}

func (list *ListCommand) printResultsAsCsv(results []ListResult) int {
	for _, result := range results {
		for _, repo := range result.Repos {
			if len(repo.Remotes) > 0 && (list.Options.remoteURL || list.Options.all) {
				for _, remote := range repo.Remotes {
					list.printResultAsCsv(result, repo, &remote)
				}
			} else {
				list.printResultAsCsv(result, repo, nil)
			}
		}
	}
	return 0
}

/*
generateFormatString returns the formatter for `Printf` to arrange the length of repository names.
*/
func (list *ListCommand) generateFormatString(repos []Repo) string {
	var max = len("Description")
	for _, repo := range repos {
		var len = len(repo.Name)
		if len > max {
			max = len
		}
	}
	return fmt.Sprintf("    %%-%ds", max)
}

func (list *ListCommand) printResults(results []ListResult) int {
	if list.Options.csv {
		return list.printResultsAsCsv(results)
	}
	for _, result := range results {
		fmt.Println(result.GroupName)
		if list.Options.description || list.Options.all {
			fmt.Printf("    Description  %s", result.Description)
			fmt.Println()
		}
		var formatString = list.generateFormatString(result.Repos)
		for _, repo := range result.Repos {
			fmt.Printf(formatString, repo.Name)
			if list.Options.localPath || list.Options.all {
				fmt.Printf("  %s", repo.Path)
			}
			if list.Options.remoteURL || list.Options.all {
				for _, remote := range repo.Remotes {
					fmt.Println()
					fmt.Printf("        %s  %s", remote.Name, remote.URL)
				}
			}
			fmt.Println()
		}
	}

	return 1
}

/*
Run performs the command.
*/
func (list *ListCommand) Run(args []string) int {
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
	list.printResults(results)
	return 0
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
    -a, --all       print all entries of each repository.
    -d, --desc      print description of group.
    -p, --path      print local paths (default).
    -r, --remote    print remote urls.
                    if any options of above are specified, '-a' are specified.

    -c, --csv       print result as csv format.
ARGUMENTS
    GROUPS    print managed repositories categorized in the groups.
              if no groups are specified, all groups are printed.`
}

func (list *ListCommand) parse(args []string) (*listOptions, error) {
	var options = listOptions{false, false, false, false, false, []string{}}
	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(list.Help()) }
	flags.BoolVar(&options.all, "a", false, "all flag")
	flags.BoolVar(&options.all, "all", false, "all flag")
	flags.BoolVar(&options.description, "d", false, "description flag")
	flags.BoolVar(&options.description, "desc", false, "description flag")
	flags.BoolVar(&options.localPath, "p", false, "local path flag")
	flags.BoolVar(&options.localPath, "path", false, "local path flag")
	flags.BoolVar(&options.remoteURL, "r", false, "remote url flag")
	flags.BoolVar(&options.remoteURL, "remote", false, "remote url flag")
	flags.BoolVar(&options.csv, "c", false, "print as csv format")
	flags.BoolVar(&options.csv, "csv", false, "print as csv format")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if !(options.all || options.description || options.localPath || options.remoteURL) {
		options.localPath = true
	}
	options.args = flags.Args()
	list.Options = &options
	return &options, nil
}
