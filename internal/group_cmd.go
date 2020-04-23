package internal

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh"
)

/*
GroupCommand represents a command.
*/
type GroupCommand struct{}
type groupAddCommand struct {
	options *groupAddOptions
}
type groupListCommand struct{}
type groupOfCommand struct{}
type groupUpdateCommand struct{}
type groupRemoveCommand struct {
	options *groupRemoveOptions
}
type groupInfoCommand struct {
}

/*
GroupCommandFactory returns an instance of command.
*/
func GroupCommandFactory() (cli.Command, error) {
	return &GroupCommand{}, nil
}

func groupInfoCommandFactory() (cli.Command, error) {
	return &groupInfoCommand{}, nil
}

func groupAddCommandFactory() (cli.Command, error) {
	return &groupAddCommand{&groupAddOptions{}}, nil
}

func groupOfCommandFactory() (cli.Command, error) {
	return &groupOfCommand{}, nil
}

func groupListCommandFactory() (cli.Command, error) {
	return &groupListCommand{}, nil
}

func groupUpdateCommandFactory() (cli.Command, error) {
	return &groupUpdateCommand{}, nil
}

func groupRemoveCommandFactory() (cli.Command, error) {
	return &groupRemoveCommand{&groupRemoveOptions{}}, nil
}

func (gac *groupAddCommand) Help() string {
	return `rrh group add [OPTIONS] <GROUPS...>
OPTIONS
    -d, --desc <DESC>        gives the description of the group.
    -o, --omit-list <FLAG>   gives the omit list flag of the group.
ARGUMENTS
    GROUPS                   gives group names.`
}

func (glc *groupListCommand) Help() string {
	return `rrh group list [OPTIONS]
OPTIONS
    -d, --desc             show description.
    -r, --repository       show repositories in the group.
    -o, --only-groupname   show only group name. This option is prioritized.`
}

func (goc *groupOfCommand) Help() string {
	return `rrh group of <REPOSITORY_ID>
ARGUMENTS
    REPOSITORY_ID     show the groups of the repository.`
}

func (grc *groupRemoveCommand) Help() string {
	return `rrh group rm [OPTIONS] <GROUPS...>
OPTIONS
    -f, --force      force remove.
    -i, --inquiry    inquiry mode.
    -v, --verbose    verbose mode.
ARGUMENTS
    GROUPS           target group names.`
}

func (gic *groupInfoCommand) Help() string {
	return `rrh group info <GROUPS...>
ARGUMENTS
    GROUPS           group names to show the information.`
}

func (guc *groupUpdateCommand) Help() string {
	return `rrh group update [OPTIONS] <GROUP>
OPTIONS
    -n, --name <NAME>        change group name to NAME.
    -d, --desc <DESC>        change description to DESC.
    -o, --omit-list <FLAG>   change omit-list of the group. FLAG must be "true" or "false".
ARGUMENTS
    GROUP               update target group names.`
}

/*
Help returns the help message of the command.
*/
func (group *GroupCommand) Help() string {
	return `rrh group <SUBCOMMAND>
SUBCOMMAND
    add       add new group.
    info      show information of specified groups.
    list      list groups (default).
    of        shows groups of the specified repository.
    rm        remove group.
    update    update group.`
}

/*
Run peforms the command.
*/
func (group *GroupCommand) Run(args []string) int {
	c := cli.NewCLI("rrh group", rrh.VERSION)
	c.Args = args
	c.Autocomplete = true
	c.Commands = map[string]cli.CommandFactory{
		"add":    groupAddCommandFactory,
		"info":   groupInfoCommandFactory,
		"update": groupUpdateCommandFactory,
		"of":     groupOfCommandFactory,
		"rm":     groupRemoveCommandFactory,
		"list":   groupListCommandFactory,
	}
	if len(args) == 0 {
		new(groupListCommand).Run([]string{})
		return 0
	}
	var exitStatus, err = c.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
	return exitStatus
}

type groupAddOptions struct {
	desc string
	omit string
	args []string
}

func (gac *groupAddCommand) buildFlagSet() (*flag.FlagSet, *groupAddOptions) {
	var opt = groupAddOptions{}
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(gac.Help()) }
	flags.StringVarP(&opt.desc, "desc", "d", "", "description")
	flags.StringVarP(&opt.omit, "omit-list", "o", "", "omit list flag")
	return flags, &opt
}

func (gac *groupAddCommand) parse(args []string) error {
	var flags, opt = gac.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return err
	}
	opt.args = flags.Args()
	gac.options = opt
	return nil
}

/*
Run performs the command.
*/
func (gac *groupAddCommand) Run(args []string) int {
	var err = gac.parse(args)
	if err != nil {
		return 1
	}
	var config = rrh.OpenConfig()
	var db, err2 = rrh.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	return gac.perform(db)
}

func (gac *groupAddCommand) perform(db *rrh.Database) int {
	if len(gac.options.args) == 0 {
		fmt.Println(gac.Help())
		return 3
	}
	if err := gac.addGroups(db, gac.options); err != nil {
		fmt.Println(err.Error())
		return 4
	}
	db.StoreAndClose()

	return 0
}

type groupListOptions struct {
	desc         bool
	repositories bool
	nameOnly     bool
}

func (glc *groupListCommand) buildFlagSet() (*flag.FlagSet, *groupListOptions) {
	var opt = groupListOptions{}
	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(glc.Help()) }
	flags.BoolVarP(&opt.desc, "desc", "d", false, "show description")
	flags.BoolVarP(&opt.repositories, "repository", "r", false, "show repositories")
	flags.BoolVarP(&opt.nameOnly, "only-groupname", "o", false, "show only group names")
	return flags, &opt
}

func (glc *groupListCommand) parse(args []string) (*groupListOptions, error) {
	var flags, opt = glc.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	return opt, nil
}

func printGroupInfo(db *rrh.Database, group *rrh.Group) {
	count := db.ContainsCount(group.Name)
	unit := "repositories"
	if count == 1 {
		unit = "repository"
	}
	fmt.Printf("%s: %s (%d %s, omit: %v)\n", group.Name, group.Description, count, unit, group.OmitList)
}

func (gic *groupInfoCommand) perform(db *rrh.Database, args []string) int {
	errs := []error{}
	for _, arg := range args {
		group := db.FindGroup(arg)
		if group == nil {
			errs = append(errs, fmt.Errorf("%s: group not found", arg))
			continue
		}
		printGroupInfo(db, group)
	}
	return printErrors(db.Config, errs)
}

func (gic *groupInfoCommand) Run(args []string) int {
	if len(args) == 0 {
		fmt.Println(gic.Help())
		return 1
	}
	var config = rrh.OpenConfig()
	var db, err = rrh.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 2
	}
	return gic.perform(db, args)
}

func (goc *groupOfCommand) perform(db *rrh.Database, repositoryID string) int {
	if !db.HasRepository(repositoryID) {
		fmt.Printf("%s: repository not found\n", repositoryID)
		return 3
	}
	var groups = db.FindRelationsOfRepository(repositoryID)
	fmt.Printf("%s, %v\n", repositoryID, groups)
	return 0
}

func (goc *groupOfCommand) Run(args []string) int {
	if len(args) != 1 {
		fmt.Println(goc.Help())
		return 1
	}
	var config = rrh.OpenConfig()
	var db, err = rrh.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 2
	}
	return goc.perform(db, args[0])
}

func printRepositoryCount(count int) {
	if count == 1 {
		fmt.Print(",1 repository")
	} else {
		fmt.Printf(",%d repositories", count)
	}
}

func findGroupName(name string, nameOnlyFlag bool, config *rrh.Config) string {
	if nameOnlyFlag {
		return name
	}
	return config.Color.ColorizedGroupName(name)
}

func (glc *groupListCommand) printResult(result groupListResult, options *groupListOptions, config *rrh.Config) {
	fmt.Print(findGroupName(result.Name, options.nameOnly, config))
	if !options.nameOnly && options.desc {
		fmt.Printf(",%s", result.Description)
	}
	if !options.nameOnly && options.repositories {
		fmt.Printf(",%v", result.Repos)
	}
	if !options.nameOnly {
		printRepositoryCount(len(result.Repos))
	}
	fmt.Println()
}

func (glc *groupListCommand) printAll(results []groupListResult, options *groupListOptions, config *rrh.Config) {
	for _, result := range results {
		glc.printResult(result, options, config)
	}
}

/*
Run performs the command.
*/
func (glc *groupListCommand) Run(args []string) int {
	var listOption, err = glc.parse(args)
	if err != nil {
		return 1
	}
	var config = rrh.OpenConfig()
	var db, err2 = rrh.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	var results = glc.listGroups(db, listOption)
	glc.printAll(results, listOption, db.Config)

	return 0
}

type groupRemoveOptions struct {
	inquiry bool
	verbose bool
	force   bool
	args    []string
}

func (grc *groupRemoveCommand) printIfVerbose(message string) {
	if grc.options.verbose {
		fmt.Println(message)
	}
}

func (grc *groupRemoveCommand) Inquiry(groupName string) bool {
	// no inquiry option, do remove group.
	if !grc.options.inquiry {
		return true
	}
	return rrh.IsInputYes(fmt.Sprintf("%s: remove group? [yN]", groupName))
}

func (grc *groupRemoveCommand) buildFlagSet() (*flag.FlagSet, *groupRemoveOptions) {
	var opt = groupRemoveOptions{}
	flags := flag.NewFlagSet("rm", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(grc.Help()) }
	flags.BoolVarP(&opt.inquiry, "inquiry", "i", false, "inquiry mode")
	flags.BoolVarP(&opt.verbose, "verbose", "v", false, "verbose mode")
	flags.BoolVarP(&opt.force, "force", "f", false, "force remove")
	return flags, &opt
}

func (grc *groupRemoveCommand) parse(args []string) (*groupRemoveOptions, error) {
	var flags, opt = grc.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	opt.args = flags.Args()
	if len(opt.args) == 0 {
		return nil, fmt.Errorf("no arguments are specified")
	}
	grc.options = opt
	return opt, nil
}

/*
Run performs the command.
*/
func (grc *groupRemoveCommand) Run(args []string) int {
	var _, err = grc.parse(args)
	if err != nil {
		return 1
	}
	var config = rrh.OpenConfig()
	var db, err2 = rrh.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	if err := grc.removeGroups(db); err != nil {
		fmt.Println(err.Error())
		return 3
	}
	db.StoreAndClose()

	return 0
}

type groupUpdateOptions struct {
	newName  string
	desc     string
	omitList string
	target   string
}

/*
Run performs the command.
*/
func (guc *groupUpdateCommand) Run(args []string) int {
	var updateOption, err = guc.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var config = rrh.OpenConfig()
	var db, err2 = rrh.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	var err3 = guc.updateGroup(db, updateOption)
	if err3 != nil {
		fmt.Println(err3.Error())
		return 3
	}
	db.StoreAndClose()
	return 0
}

func (guc *groupUpdateCommand) buildFlagSet() (*flag.FlagSet, *groupUpdateOptions) {
	var opt = groupUpdateOptions{}
	flags := flag.NewFlagSet("update", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(guc.Help()) }
	flags.StringVarP(&opt.newName, "name", "n", "", "specify new group name")
	flags.StringVarP(&opt.desc, "desc", "d", "", "specify the description")
	flags.StringVarP(&opt.omitList, "omit-list", "o", "", "set the omit list flag.")
	return flags, &opt
}

func (guc *groupUpdateCommand) parse(args []string) (*groupUpdateOptions, error) {
	var flags, opt = guc.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	var arguments = flags.Args()
	if len(arguments) == 0 {
		return nil, fmt.Errorf("no arguments are specified")
	}
	if len(arguments) > 1 {
		return nil, fmt.Errorf("could not accept multiple arguments")
	}
	opt.target = arguments[0]
	return opt, nil
}

/*
Synopsis returns the help message of the command.
*/
func (group *GroupCommand) Synopsis() string {
	return "add/list/update/remove groups and show groups of the repository."
}

/*
Synopsis returns the help message of the command.
*/
func (gac *groupAddCommand) Synopsis() string {
	return "add group."
}

/*
Synopsis returns the help message of the command.
*/
func (glc *groupListCommand) Synopsis() string {
	return "list groups."
}

func (goc *groupOfCommand) Synopsis() string {
	return "show groups of the repository."
}

/*
Synopsis returns the help message of the command.
*/
func (grc *groupRemoveCommand) Synopsis() string {
	return "remove given group."
}

/*
Synopsis returns the help message of the command.
*/
func (gic *groupInfoCommand) Synopsis() string {
	return "show information of groups."
}

/*
Synopsis returns the help message of the command.
*/
func (guc *groupUpdateCommand) Synopsis() string {
	return "update group."
}

type groupListResult struct {
	Name        string
	Description string
	Repos       []string
}

func appendRelations(groupName string, relations []rrh.Relation) []string {
	var repos = []string{}
	for _, relation := range relations {
		if relation.GroupName == groupName {
			repos = append(repos, relation.RepositoryID)
		}
	}
	return repos
}

func (glc *groupListCommand) listGroups(db *rrh.Database, listOptions *groupListOptions) []groupListResult {
	var results = []groupListResult{}
	for _, group := range db.Groups {
		var result = groupListResult{group.Name, group.Description, []string{}}
		result.Repos = appendRelations(group.Name, db.Relations)
		results = append(results, result)
	}
	return results
}

func trueOrFalse(flag string) bool {
	var flagString = strings.ToLower(flag)
	if flagString == "true" {
		return true
	}
	return false
}

func (gac *groupAddCommand) addGroups(db *rrh.Database, options *groupAddOptions) error {
	for _, groupName := range options.args {
		var flag = trueOrFalse(options.omit)
		var _, err = db.CreateGroup(groupName, options.desc, flag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (grc *groupRemoveCommand) removeGroupsImpl(db *rrh.Database, groupName string) error {
	if grc.options.force {
		db.ForceDeleteGroup(groupName)
		grc.printIfVerbose(fmt.Sprintf("%s: group removed", groupName))
	} else if db.ContainsCount(groupName) == 0 {
		db.DeleteGroup(groupName)
		grc.printIfVerbose(fmt.Sprintf("%s: group removed", groupName))
	} else {
		return fmt.Errorf("%s: cannot remove group. the group has relations", groupName)
	}
	return nil
}

func (grc *groupRemoveCommand) removeGroups(db *rrh.Database) error {
	for _, groupName := range grc.options.args {
		if !db.HasGroup(groupName) || !grc.Inquiry(groupName) {
			return nil
		}
		if err := grc.removeGroupsImpl(db, groupName); err != nil {
			return err
		}
	}
	return nil
}

func createNewGroup(opt *groupUpdateOptions, prevGroup *rrh.Group) rrh.Group {
	var newGroup = rrh.Group{Name: opt.newName, Description: opt.desc, OmitList: strings.ToLower(opt.omitList) == "true"}
	if opt.desc == "" {
		newGroup.Description = prevGroup.Description
	}
	if opt.newName == "" {
		newGroup.Name = prevGroup.Name
	}
	if opt.omitList == "" {
		newGroup.OmitList = prevGroup.OmitList
	}
	return newGroup
}

func (guc *groupUpdateCommand) updateGroup(db *rrh.Database, opt *groupUpdateOptions) error {
	if !db.HasGroup(opt.target) {
		return fmt.Errorf("%s: group not found", opt.target)
	}
	var newGroup = createNewGroup(opt, db.FindGroup(opt.target))
	if !db.UpdateGroup(opt.target, newGroup) {
		return fmt.Errorf("%s: failed to update to {%s, %s, %s}", opt.target, opt.newName, opt.desc, opt.omitList)
	}
	return nil
}
