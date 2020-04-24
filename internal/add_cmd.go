package internal

import (
	"fmt"
	"path/filepath"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh"
)

/*
AddCommand shows the subcommand of rrh.
*/
type AddCommand struct {
	options *addOptions
}

/*
AddCommandFactory generates the object of AddCommand.
*/
func AddCommandFactory() (cli.Command, error) {
	return &AddCommand{options: &addOptions{}}, nil
}

/*
Help function shows the help message.
*/
func (add *AddCommand) Help() string {
	return `rrh add [OPTIONS] <REPOSITORY_PATHS...>
OPTIONS
    -g, --group <GROUP>        add repository to rrh database.
    -r, --repository-id <ID>   specified repository id of the given repository path.
                               Specifying this option fails with multiple arguments.
ARGUMENTS
    REPOSITORY_PATHS           the local path list of the git repositories.`
}

/*
Synopsis returns the simple help message of the command.
*/
func (add *AddCommand) Synopsis() string {
	return "add repositories on the local path to rrh."
}

func (add *AddCommand) showError(errorlist []error, onError string) {
	if len(errorlist) == 0 || onError == rrh.Ignore {
		return
	}
	for _, item := range errorlist {
		fmt.Println(item.Error())
	}
}

func (add *AddCommand) perform(db *rrh.Database, opt *addOptions) int {
	var onError = db.Config.GetValue(rrh.OnError)

	var errorlist = add.AddRepositoriesToGroup(db, opt)

	add.showError(errorlist, onError)

	if onError == rrh.Fail || onError == rrh.FailImmediately {
		return 1
	}
	var err2 = db.StoreAndClose()
	if err2 != nil {
		fmt.Println(err2.Error())
	}

	return 0
}

/*
Run function performs the command.
*/
func (add *AddCommand) Run(args []string) int {
	var config = rrh.OpenConfig()
	var opt, err = add.parse(args, config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var db, err2 = rrh.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	return add.perform(db, opt)
}

type addOptions struct {
	group  string
	repoID string
	args   []string
}

func (add *AddCommand) buildFlagSet(config *rrh.Config) (*flag.FlagSet, *addOptions) {
	var opt = addOptions{}
	var defaultGroup = config.GetValue(rrh.DefaultGroupName)
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(add.Help()) }
	flags.StringVarP(&opt.group, "group", "g", defaultGroup, "target group")
	flags.StringVarP(&opt.repoID, "repository-id", "r", "", "specifying repository id")
	return flags, &opt
}

func (add *AddCommand) parse(args []string, config *rrh.Config) (*addOptions, error) {
	var flags, opt = add.buildFlagSet(config)
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	opt.args = flags.Args()
	add.options = opt

	return opt, nil
}

func isDuplicateRepository(db *rrh.Database, repoID string, path string) error {
	var repo = db.FindRepository(repoID)
	if repo != nil && repo.Path != path {
		return fmt.Errorf("%s: duplicate repository id", repoID)
	}
	return nil
}

func findIDFromPath(repoID string, absPath string) string {
	if repoID == "" {
		return filepath.Base(absPath)
	}
	return repoID
}

func (add *AddCommand) addRepositoryToGroup(db *rrh.Database, rel rrh.Relation, path string) []error {
	var absPath, _ = filepath.Abs(path)
	var id = findIDFromPath(rel.RepositoryID, absPath)
	if err1 := rrh.IsExistAndGitRepository(absPath, path); err1 != nil {
		return []error{err1}
	}
	if err1 := isDuplicateRepository(db, id, absPath); err1 != nil {
		return []error{err1}
	}
	var remotes, err2 = rrh.FindRemotes(absPath)
	if err2 != nil {
		return []error{err2}
	}
	db.CreateRepository(id, absPath, "", remotes)

	var err = db.Relate(rel.GroupName, id)
	if err != nil {
		return []error{fmt.Errorf("%s: cannot create relation to group %s", id, rel.GroupName)}
	}
	return []error{}
}

func validateArguments(args []string, repoID string) error {
	if repoID != "" && len(args) > 1 {
		return fmt.Errorf("specifying repository id do not accept multiple arguments: %v", args)
	}
	return nil
}

/*
AddRepositoriesToGroup registers the given repositories to the specified group.
*/
func (add *AddCommand) AddRepositoriesToGroup(db *rrh.Database, opt *addOptions) []error {
	var _, err = db.AutoCreateGroup(opt.group, "", false)
	if err != nil {
		return []error{err}
	}
	if err := validateArguments(opt.args, opt.repoID); err != nil {
		return []error{err}
	}
	var errorlist = []error{}
	for _, item := range opt.args {
		var list = add.addRepositoryToGroup(db, rrh.Relation{RepositoryID: opt.repoID, GroupName: opt.group}, item)
		errorlist = append(errorlist, list...)
	}
	return errorlist
}
