package status

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type StatusCommand struct {
}

type statusOptions struct {
	csv    bool
	branch bool
	remote bool
	args   []string
}

func (options *statusOptions) isRemoteTarget(name plumbing.ReferenceName) bool {
	return options.remote && name.IsRemote()
}

func (options *statusOptions) isBranchTarget(name plumbing.ReferenceName) bool {
	return options.branch && name.IsBranch()
}

func StatusCommandFactory() (cli.Command, error) {
	return &StatusCommand{}, nil
}

func (status *StatusCommand) Help() string {
	return `rrh status [OPTIONS] [REPOSITORIES|GROUPS...]
OPTIONS
    -b, --branches  show the status of the local branches.
	-r, --remote    show the status of the remote branches.
    -c, --csv       print result in csv format.
ARGUMENTS
    REPOSITORIES    target repositories.  If no repository was specified
                    the command shows the result of the default group.
    GROUPS          target groups.  If no group was specified,
                    the command shows the result of the default group.`
}

func (status *StatusCommand) parseFmtString(results []StatusResult) string {
	var max = 0
	for _, result := range results {
		var len = len(result.BranchName)
		if len > max {
			max = len
		}
	}
	return fmt.Sprintf("        %%-%ds    %%-12s    %%s\n", max)
}

func (status *StatusCommand) printResultInCsv(results []StatusResult, config *common.Config) {
	for _, result := range results {
		fmt.Printf("%s,%s,%s,%s,%s\n", result.GroupName, result.RepositoryName, result.BranchName, common.Strftime(*result.LastModified, config), result.Description)
	}
}

func (status *StatusCommand) printResult(results []StatusResult, config *common.Config) {
	var groupName = results[0].GroupName
	var repositoryName = results[0].RepositoryName
	fmt.Printf("%s\n    %s\n", groupName, repositoryName)
	var fmtString = status.parseFmtString(results)
	for _, result := range results {
		if groupName != result.GroupName {
			fmt.Println(result.GroupName)
			groupName = result.GroupName
		}
		if repositoryName != result.RepositoryName {
			fmt.Printf("    %s\n", result.RepositoryName)
			repositoryName = result.RepositoryName
		}
		var time = ""
		if result.LastModified != nil {
			time = common.Strftime(*result.LastModified, config)
		}
		fmt.Printf(fmtString, result.BranchName, time, result.Description)
	}
}

func (status *StatusCommand) runStatus(db *common.Database, arg string, options *statusOptions) int {
	var errorFlag = 0
	var result, err = status.executeStatus(db, arg, options)
	if len(err) != 0 {
		for _, item := range err {
			fmt.Println(item.Error())
			errorFlag = 1
		}
	} else {
		if options.csv {
			status.printResultInCsv(result, db.Config)
		} else {
			status.printResult(result, db.Config)
		}
	}
	return errorFlag
}

func (status *StatusCommand) Run(args []string) int {
	var config = common.OpenConfig()
	options, err := status.parse(args, config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	db, err := common.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var errorFlag = 0
	for _, arg := range options.args {
		errorFlag += status.runStatus(db, arg, options)
	}

	return errorFlag
}

func (status *StatusCommand) parse(args []string, config *common.Config) (*statusOptions, error) {
	var options = statusOptions{false, false, false, []string{}}
	flags := flag.NewFlagSet("status", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(status.Help()) }
	flags.BoolVar(&options.csv, "c", false, "csv format")
	flags.BoolVar(&options.csv, "csv", false, "csv format")
	flags.BoolVar(&options.remote, "r", false, "remote branch status")
	flags.BoolVar(&options.remote, "remote", false, "remote branch status")
	flags.BoolVar(&options.branch, "b", false, "local branch status")
	flags.BoolVar(&options.branch, "branch", false, "local branch status")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	options.args = flags.Args()
	if len(options.args) == 0 {
		options.args = []string{config.GetValue(common.RrhDefaultGroupName)}
	}
	return &options, nil
}

/*
Synopsis returns the help message of the command.
*/
func (status *StatusCommand) Synopsis() string {
	return "show git status of repositories."
}
