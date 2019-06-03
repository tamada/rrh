package internal

import (
	"fmt"
	"time"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/lib"
)

/*
StatusCommand represents a command.
*/
type StatusCommand struct {
	options *statusOptions
}

const timeformat = "2006-01-02 03:04:05-07"

const (
	relative     = "relative"
	absolute     = "absolute"
	notSpecified = "not_specified"
)

type statusOptions struct {
	csv    bool
	option *lib.StatusOption
	format string
}

func (options *statusOptions) strftime(time *time.Time, config *lib.Config) string {
	if time == nil {
		return ""
	}
	switch options.format {
	case relative:
		return lib.HumanizeTime(*time)
	case notSpecified:
		return lib.Strftime(*time, config)
	}
	return time.Format(timeformat)
}

/*
StatusCommandFactory returns an instance of the StatusCommand.
*/
func StatusCommandFactory() (cli.Command, error) {
	return &StatusCommand{&statusOptions{false, lib.NewStatusOption(), notSpecified}}, nil
}

/*
Help returns the help message for the user.
*/
func (status *StatusCommand) Help() string {
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

func (status *StatusCommand) parseFmtString(results []lib.Status) string {
	var max = 0
	for _, result := range results {
		var len = len(result.BranchName)
		if len > max {
			max = len
		}
	}
	return fmt.Sprintf("        %%-%ds    %%-22s    %%s\n", max)
}

func (status *StatusCommand) printResultInCsv(results []lib.Status, config *lib.Config) {
	for _, result := range results {
		var timeString = status.options.strftime(result.LastModified, config)
		fmt.Printf("%s,%s,%s,%s,%s\n", result.Relation.GroupName, result.Relation.RepositoryID, result.BranchName, timeString, result.Description)
	}
}

func (status *StatusCommand) printResult(results []lib.Status, config *lib.Config) {
	var groupName = results[0].Relation.GroupName
	var repositoryName = results[0].Relation.RepositoryID
	fmt.Printf("%s\n    %s\n", config.Color.ColorizedGroupName(groupName), config.Color.ColorizedRepositoryID(repositoryName))
	var fmtString = status.parseFmtString(results)
	for _, result := range results {
		if groupName != result.Relation.GroupName {
			fmt.Println(config.Color.ColorizedGroupName(result.Relation.GroupName))
			groupName = result.Relation.GroupName
		}
		if repositoryName != result.Relation.RepositoryID {
			fmt.Printf("    %s\n", config.Color.ColorizedRepositoryID(result.Relation.RepositoryID))
			repositoryName = result.Relation.RepositoryID
		}
		var time = ""
		if result.LastModified != nil {
			time = status.options.strftime(result.LastModified, config)
		}
		fmt.Printf(fmtString, result.BranchName, time, result.Description)
	}
}

func (status *StatusCommand) runStatus(db *lib.Database, arg string) int {
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
func (status *StatusCommand) Run(args []string) int {
	var config = lib.OpenConfig()
	arguments, err := status.parse(args, config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	db, err := lib.Open(config)
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

func (status *StatusCommand) buildFlagSet() (*flag.FlagSet, *statusOptions) {
	var options = statusOptions{false, lib.NewStatusOption(), notSpecified}
	flags := flag.NewFlagSet("status", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(status.Help()) }
	flags.BoolVarP(&options.csv, "csv", "c", false, "csv format")
	flags.BoolVarP(&options.option.RemoteStatus, "remote", "r", false, "remote branch status")
	flags.BoolVarP(&options.option.BranchStatus, "branches", "b", false, "local branch status")
	flags.StringVarP(&options.format, "time-format", "f", notSpecified, "specifies time format")
	return flags, &options
}

func (status *StatusCommand) parse(args []string, config *lib.Config) ([]string, error) {
	var flags, options = status.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	status.options = options
	if len(flags.Args()) == 0 {
		return []string{config.GetValue(lib.RrhDefaultGroupName)}, nil
	}
	return flags.Args(), nil
}

/*
Synopsis returns the help message of the command.
*/
func (status *StatusCommand) Synopsis() string {
	return "show git status of repositories."
}

func (status *StatusCommand) executeStatus(db *lib.Database, name string) ([]lib.Status, []error) {
	if db.HasGroup(name) {
		return status.executeStatusOnGroup(db, name)
	}
	if db.HasRepository(name) {
		var results, err = status.options.option.StatusOfRepository(db, &lib.Relation{GroupName: "unknown-group", RepositoryID: name})
		if err != nil {
			return results, []error{err}
		}
		return results, []error{}
	}
	return nil, []error{fmt.Errorf("%s: group and repository not found", name)}
}

func (status *StatusCommand) executeStatusOnGroup(db *lib.Database, groupName string) ([]lib.Status, []error) {
	var group = db.FindGroup(groupName)
	if group == nil {
		return nil, []error{fmt.Errorf("%s: group not found", groupName)}
	}
	var errors = []error{}
	var results = []lib.Status{}
	for _, repoID := range db.FindRelationsOfGroup(groupName) {
		var sr, err = status.options.option.StatusOfRepository(db, &lib.Relation{GroupName: groupName, RepositoryID: repoID})
		if err != nil {
			errors = append(errors, err)
		} else {
			results = append(results, sr...)
		}
	}
	return results, errors
}
