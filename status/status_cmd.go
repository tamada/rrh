package status

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

/*
Command represents a command.
*/
type Command struct {
	options *options
}

type options struct {
	csv    bool
	branch bool
	remote bool
	args   []string
}

func (options *options) isRemoteTarget(name plumbing.ReferenceName) bool {
	return options.remote && name.IsRemote()
}

func (options *options) isBranchTarget(name plumbing.ReferenceName) bool {
	return options.branch && name.IsBranch()
}

/*
CommandFactory returns an instance of the StatusCommand.
*/
func CommandFactory() (cli.Command, error) {
	return &Command{&options{false, false, false, []string{}}}, nil
}

/*
Help returns the help message for the user.
*/
func (status *Command) Help() string {
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

func (status *Command) parseFmtString(results []result) string {
	var max = 0
	for _, result := range results {
		var len = len(result.BranchName)
		if len > max {
			max = len
		}
	}
	return fmt.Sprintf("        %%-%ds    %%-12s    %%s\n", max)
}

func (status *Command) printResultInCsv(results []result, config *common.Config) {
	for _, result := range results {
		fmt.Printf("%s,%s,%s,%s,%s\n", result.GroupName, result.RepositoryName, result.BranchName, common.Strftime(*result.LastModified, config), result.Description)
	}
}

func (status *Command) printResult(results []result, config *common.Config) {
	var groupName = results[0].GroupName
	var repositoryName = results[0].RepositoryName
	fmt.Printf("%s\n    %s\n", common.ColorrizedGroupName(groupName), common.ColorrizedRepositoryID(repositoryName))
	var fmtString = status.parseFmtString(results)
	for _, result := range results {
		if groupName != result.GroupName {
			fmt.Println(common.ColorrizedGroupName(result.GroupName))
			groupName = result.GroupName
		}
		if repositoryName != result.RepositoryName {
			fmt.Printf("    %s\n", common.ColorrizedRepositoryID(result.RepositoryName))
			repositoryName = result.RepositoryName
		}
		var time = ""
		if result.LastModified != nil {
			time = common.Strftime(*result.LastModified, config)
		}
		fmt.Printf(fmtString, result.BranchName, time, result.Description)
	}
}

func (status *Command) runStatus(db *common.Database, arg string) int {
	var errorFlag = 0
	var result, err = status.executeStatus(db, arg)
	if len(err) != 0 {
		for _, item := range err {
			fmt.Println(item.Error())
			errorFlag = 1
		}
	} else {
		if status.options.csv {
			status.printResultInCsv(result, db.Config)
		} else {
			status.printResult(result, db.Config)
		}
	}
	return errorFlag
}

/*
Run performs the command.
*/
func (status *Command) Run(args []string) int {
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
		errorFlag += status.runStatus(db, arg)
	}

	return errorFlag
}

func (status *Command) buildFlagSet() (*flag.FlagSet, *options) {
	var options = options{false, false, false, []string{}}
	flags := flag.NewFlagSet("status", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(status.Help()) }
	flags.BoolVar(&options.csv, "c", false, "csv format")
	flags.BoolVar(&options.csv, "csv", false, "csv format")
	flags.BoolVar(&options.remote, "r", false, "remote branch status")
	flags.BoolVar(&options.remote, "remote", false, "remote branch status")
	flags.BoolVar(&options.branch, "b", false, "local branch status")
	flags.BoolVar(&options.branch, "branches", false, "local branch status")
	return flags, &options
}

func (status *Command) parse(args []string, config *common.Config) (*options, error) {
	var flags, options = status.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	options.args = flags.Args()
	if len(options.args) == 0 {
		options.args = []string{config.GetValue(common.RrhDefaultGroupName)}
	}
	status.options = options
	return options, nil
}

/*
Synopsis returns the help message of the command.
*/
func (status *Command) Synopsis() string {
	return "show git status of repositories."
}
