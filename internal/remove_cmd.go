package internal

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/lib"
)

type removeOptions struct {
	inquiry   bool
	recursive bool
	verbose   bool
	args      []string
}

/*
RemoveCommand represents a command.
*/
type RemoveCommand struct {
	options *removeOptions
}

/*
RemoveCommandFactory returns an instance of the RemoveCommand.
*/
func RemoveCommandFactory() (cli.Command, error) {
	return &RemoveCommand{&removeOptions{}}, nil
}

func (options *removeOptions) printIfVerbose(message string) {
	if options.verbose {
		fmt.Println(message)
	}
}

func (rm *RemoveCommand) executeRemoveGroup(db *lib.Database, groupName string) error {
	var group = db.FindGroup(groupName)
	if group == nil {
		return fmt.Errorf("%s: group not found", groupName)
	}
	if rm.options.inquiry && !lib.IsInputYes(fmt.Sprintf("%s: Remove group? [yN]> ", groupName)) {
		rm.options.printIfVerbose(fmt.Sprintf("%s: group do not removed", groupName))
		return nil
	}
	var count = db.ContainsCount(groupName)
	if !rm.options.recursive && count > 0 {
		return fmt.Errorf("%s: cannot remove, it contains %d repository(es)", group.Name, count)
	}
	db.UnrelateFromGroup(groupName)
	var err = db.DeleteGroup(groupName)
	if err == nil {
		rm.options.printIfVerbose(fmt.Sprintf("%s: group removed", group.Name))
	}
	return err
}

func (rm *RemoveCommand) executeRemoveRepository(db *lib.Database, repoID string) error {
	if !db.HasRepository(repoID) {
		return fmt.Errorf("%s: repository not found", repoID)
	}
	if rm.options.inquiry && !lib.IsInputYes(fmt.Sprintf("%s: Remove repository? [yN]> ", repoID)) {
		rm.options.printIfVerbose(fmt.Sprintf("%s: repository do not removed", repoID))
		return nil
	}
	if err := db.DeleteRepository(repoID); err != nil {
		return err
	}
	rm.options.printIfVerbose(fmt.Sprintf("%s: repository removed", repoID))
	return nil
}

func (rm *RemoveCommand) executeRemoveFromGroup(db *lib.Database, groupName string, repoID string) error {
	db.Unrelate(groupName, repoID)
	rm.options.printIfVerbose(fmt.Sprintf("%s: removed from group %s", repoID, groupName))
	return nil
}

func (rm *RemoveCommand) executeRemove(db *lib.Database, target string) error {
	var data = strings.Split(target, "/")
	if len(data) == 2 {
		return rm.executeRemoveFromGroup(db, data[0], data[1])
	}
	var repoFlag = db.HasRepository(target)
	var groupFlag = db.HasGroup(target)
	if repoFlag && groupFlag {
		return fmt.Errorf("%s: exists in repositories and groups", target)
	}
	if repoFlag {
		return rm.executeRemoveRepository(db, target)
	}
	if groupFlag {
		return rm.executeRemoveGroup(db, target)
	}
	return fmt.Errorf("%s: not found in repositories and groups", target)
}

func (rm *RemoveCommand) perform(db *lib.Database) int {
	var result = 0
	for _, target := range rm.options.args {
		var err = rm.executeRemove(db, target)
		if err != nil {
			fmt.Println(err.Error())
			result = 3
		}
	}
	if result == 0 {
		if db.Config.IsSet(lib.RrhAutoDeleteGroup) {
			db.Prune()
		}
		db.StoreAndClose()
	}
	return result
}

/*
Run performs the command.
*/
func (rm *RemoveCommand) Run(args []string) int {
	var options, err = rm.parse(args)
	if err != nil {
		return 1
	}
	rm.options = options
	var config = lib.OpenConfig()
	var db, err1 = lib.Open(config)
	if err1 != nil {
		fmt.Println(err1.Error())
		return 2
	}
	return rm.perform(db)
}

func (rm *RemoveCommand) buildFlagSet() (*flag.FlagSet, *removeOptions) {
	var options = removeOptions{false, false, false, []string{}}
	var flags = flag.NewFlagSet("rm", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(rm.Help()) }
	flags.BoolVarP(&options.inquiry, "inquiry", "i", false, "inquiry flag")
	flags.BoolVarP(&options.verbose, "verbose", "v", false, "verbose flag")
	flags.BoolVarP(&options.recursive, "recursive", "r", false, "recursive flag")
	return flags, &options
}

func (rm *RemoveCommand) parse(args []string) (*removeOptions, error) {
	var flags, options = rm.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	options.args = flags.Args()
	return options, nil
}

/*
Help returns the help message.
*/
func (rm *RemoveCommand) Help() string {
	return `rrh rm [OPTIONS] <REPO_ID|GROUP_ID|GROUP_ID/REPO_ID...>
OPTIONS
    -i, --inquiry       inquiry mode.
    -r, --recursive     recursive mode.
    -v, --verbose       verbose mode.

ARGUMENTS
    REPOY_ID            repository name for removing.
    GROUP_ID            group name. if the group contains repositories,
                        remove will fail without '-r' option.
    GROUP_ID/REPO_ID    remove the relation between the given REPO_ID and GROUP_ID.`
}

/*
Synopsis returns the help message of the command.
*/
func (rm *RemoveCommand) Synopsis() string {
	return "remove given repository from database."
}
