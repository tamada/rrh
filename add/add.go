package add

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/cli"
	"github.com/tamadalab/grim/common"
)

type AddCommand struct {
}

func AddCommandFactory() (cli.Command, error) {
	return &AddCommand{}, nil
}

func (add *AddCommand) Help() string {
	return `grim add [OPTION] <REPOSITORY_PATHS...>
OPTION
    -g <GROUP>        add repository to GRIM database.
ARGUMENTS
    REPOSITORY_PATHS  local path list of git repository.`
}

func findGroup(db *common.Database, groupName string, config *common.Config) (*common.Group, error) {
	var group = db.FindGroup(groupName)
	if group == nil {
		if config.GetValue(common.GrimAutoCreateGroup) == "true" {
			return db.AddGroup(groupName, ""), nil
		}
		return nil, fmt.Errorf("%s: group not found", groupName)
	}
	return group, nil
}

func (add *AddCommand) showError(errorlist []error, onError string) {
	if len(errorlist) == 0 || onError == common.Ignore {
		return
	}
	for _, item := range errorlist {
		fmt.Fprintf(os.Stderr, "%s\n", item.Error())
	}
}

func (add *AddCommand) addRepositoryToGroup(db *common.Database, group *common.Group, path string, config *common.Config, list []error) []error {
	var absPath, _ = filepath.Abs(path)
	var id = filepath.Base(absPath)
	var repoPath = common.NormalizePath(absPath)
	var remote, err = common.FindRemoteUrlFromRepository(absPath)
	if err != nil {
		list = append(list, err)
		if config.GetValue(common.GrimOnError) == common.FailImmediately {
			return list
		}
	} else {
		var repo = common.Repository{ID: id, Path: repoPath, URL: remote}
		db.AddRepository(&repo, group)
	}

	return list
}

func (add *AddCommand) perform(db *common.Database, config *common.Config, args []string, groupName string) int {
	var group, err = findGroup(db, groupName, config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	var onError = config.GetValue(common.GrimOnError)
	var errorlist = []error{}
	for _, item := range args {
		errorlist = add.addRepositoryToGroup(db, group, item, config, errorlist)
	}
	add.showError(errorlist, onError)

	if onError == common.Fail {
		return 1
	}
	var err2 = db.StoreAndClose(config)
	if err2 != nil {
		fmt.Println(err2.Error())
	}

	return 0
}

func (add *AddCommand) Run(args []string) int {
	var config = common.OpenConfig()
	var opt, err = add.parse(args, config)
	if err != nil {
		fmt.Println(add.Help())
		return 1
	}
	return add.perform(common.Open(config), config, opt.args, opt.group)
}

type addoptions struct {
	group string
	args  []string
}

func (add *AddCommand) parse(args []string, config *common.Config) (*addoptions, error) {
	var opt = addoptions{}
	flags := flag.NewFlagSet("add", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(add.Help()) }
	flags.StringVar(&opt.group, "g", config.GetValue(common.GrimDefaultGroupName), "target group")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	opt.args = flags.Args()
	if opt.group == "" {
		opt.group = config.GetValue(common.GrimDefaultGroupName)
	}

	return &opt, nil
}

func (add *AddCommand) Synopsis() string {
	return "add repository on the local path to GRIM"
}
