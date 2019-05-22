package status

import (
	"fmt"
	"time"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/common"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

/*
Command represents a command.
*/
type Command struct {
	options *options
}

const timeformat = "2006-01-02 03:04:05-07"

const (
	relative     = "relative"
	absolute     = "absolute"
	notSpecified = "not_specified"
)

type options struct {
	csv    bool
	branch bool
	remote bool
	format string
}

func (options *options) strftime(time *time.Time, config *common.Config) string {
	if time == nil {
		return ""
	}
	switch options.format {
	case relative:
		return common.HumanizeTime(*time)
	case notSpecified:
		return common.Strftime(*time, config)
	}
	return time.Format(timeformat)
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
	return &Command{&options{false, false, false, notSpecified}}, nil
}

/*
Help returns the help message for the user.
*/
func (status *Command) Help() string {
	return `rrh status [OPTIONS] [REPOSITORIES|GROUPS...]
OPTIONS
    -b, --branches               show the status of the local branches.
    -r, --remote                 show the status of the remote branches.
    -c, --csv                    print result in csv format.
    -f, --time-format <FORMAT>   specifies time format. Available value is
                                 'relative' ad 'absolute'
ARGUMENTS
    REPOSITORIES                 target repositories.  If no repository was specified
                                 the command shows the result of the default group.
    GROUPS                       target groups.  If no group was specified,
                                 the command shows the result of the default group.`
}

func (status *Command) parseFmtString(results []result) string {
	var max = 0
	for _, result := range results {
		var len = len(result.branchName)
		if len > max {
			max = len
		}
	}
	return fmt.Sprintf("        %%-%ds    %%-22s    %%s\n", max)
}

func (status *Command) printResultInCsv(results []result, config *common.Config) {
	for _, result := range results {
		var timeString = status.options.strftime(result.lastModified, config)
		fmt.Printf("%s,%s,%s,%s,%s\n", result.relation.gname, result.relation.rname, result.branchName, timeString, result.description)
	}
}

func (status *Command) printResult(results []result, config *common.Config) {
	var groupName = results[0].relation.gname
	var repositoryName = results[0].relation.rname
	fmt.Printf("%s\n    %s\n", config.Color.ColorizedGroupName(groupName), config.Color.ColorizedRepositoryID(repositoryName))
	var fmtString = status.parseFmtString(results)
	for _, result := range results {
		if groupName != result.relation.gname {
			fmt.Println(config.Color.ColorizedGroupName(result.relation.gname))
			groupName = result.relation.gname
		}
		if repositoryName != result.relation.rname {
			fmt.Printf("    %s\n", config.Color.ColorizedRepositoryID(result.relation.rname))
			repositoryName = result.relation.rname
		}
		var time = ""
		if result.lastModified != nil {
			time = status.options.strftime(result.lastModified, config)
		}
		fmt.Printf(fmtString, result.branchName, time, result.description)
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
	arguments, err := status.parse(args, config)
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
	for _, arg := range arguments {
		errorFlag += status.runStatus(db, arg)
	}

	return errorFlag
}

func (status *Command) buildFlagSet() (*flag.FlagSet, *options) {
	var options = options{false, false, false, notSpecified}
	flags := flag.NewFlagSet("status", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(status.Help()) }
	flags.BoolVarP(&options.csv, "csv", "c", false, "csv format")
	flags.BoolVarP(&options.remote, "remote", "r", false, "remote branch status")
	flags.BoolVarP(&options.branch, "branches", "b", false, "local branch status")
	flags.StringVarP(&options.format, "time-format", "f", notSpecified, "specifies time format")
	return flags, &options
}

func (status *Command) parse(args []string, config *common.Config) ([]string, error) {
	var flags, options = status.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	status.options = options
	if len(flags.Args()) == 0 {
		return []string{config.GetValue(common.RrhDefaultGroupName)}, nil
	}
	return flags.Args(), nil
}

/*
Synopsis returns the help message of the command.
*/
func (status *Command) Synopsis() string {
	return "show git status of repositories."
}
