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

type ListCommand struct{}

func ListCommandFactory() (cli.Command, error) {
	return &ListCommand{}, nil
}

func (list *ListCommand) printResultAsCsv(result ListResult, repo Repo, remote *common.Remote, options *listOptions) {
	fmt.Printf("%s", result.GroupName)
	if options.description || options.all {
		fmt.Printf(",%s", result.Description)
	}
	fmt.Printf(",%s", repo.Name)
	if options.localPath || options.all {
		fmt.Printf(",%s", repo.Path)
	}
	if (options.remoteURL || options.all) && remote != nil {
		fmt.Printf(",%s,%s", remote.Name, remote.URL)
	}
	fmt.Println()
}

func (list *ListCommand) printResultsAsCsv(results []ListResult, options *listOptions) int {

	for _, result := range results {
		for _, repo := range result.Repos {
			if options.remoteURL || options.all {
				for _, remote := range repo.Remotes {
					list.printResultAsCsv(result, repo, &remote, options)
				}
			} else {
				list.printResultAsCsv(result, repo, nil, options)
			}
		}
	}
	return 0
}

func (list *ListCommand) printResults(results []ListResult, options *listOptions) int {
	if options.csv {
		return list.printResultsAsCsv(results, options)
	}
	for _, result := range results {
		fmt.Println(result.GroupName)
		if options.description || options.all {
			fmt.Printf("    Description: %s\n", result.Description)
		}
		fmt.Println("    Repositories:")
		for _, repo := range result.Repos {
			fmt.Printf("        %s", repo.Name)
			if options.localPath || options.all {
				fmt.Printf(",%s", repo.Path)
			}
			if options.remoteURL || options.all {
				for _, remote := range repo.Remotes {
					fmt.Printf("\n            %s,%s", remote.Name, remote.URL)
				}
			}
			fmt.Println()
		}
	}

	return 1
}

func (list *ListCommand) Run(args []string) int {
	options, err := list.parse(args)
	if err != nil {
		fmt.Printf(list.Help())
		return 1
	}
	var config = common.OpenConfig()
	db, err := common.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	results, err := list.FindResults(db, options)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	list.printResults(results, options)
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
    -a, --all       print all (default).
    -d, --desc      print description of group.
    -p, --path      print local paths.
    -r, --remote    print remote urls.
                    if any options of above are specified, '-a' are specified.

    -c, --csv       print result as csv format.

ARGUMENTS
    GROUPS    print managed repositories categoried in the groups.
              if no groups are specified, default groups are printed.`
}

func (list *ListCommand) parse(args []string) (*listOptions, error) {
	var options = listOptions{false, false, false, false, false, []string{}}
	flags := flag.NewFlagSet("list", flag.ExitOnError)
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
		options.all = true
	}
	options.args = flags.Args()
	return &options, nil
}