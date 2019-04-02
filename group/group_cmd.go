package group

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

/*
GroupCommand represents a command.
*/
type GroupCommand struct{}
type groupAddCommand struct{}
type groupListCommand struct{}
type groupUpdateCommand struct{}
type groupRemoveCommand struct {
	Options *removeOptions
}

/*
GroupCommandFactory returns an instance of command.
*/
func GroupCommandFactory() (cli.Command, error) {
	return &GroupCommand{}, nil
}

func groupAddCommandFactory() (cli.Command, error) {
	return &groupAddCommand{}, nil
}

func groupListCommandFactory() (cli.Command, error) {
	return &groupListCommand{}, nil
}

func groupUpdateCommandFactory() (cli.Command, error) {
	return &groupUpdateCommand{}, nil
}

func groupRemoveCommandFactory() (cli.Command, error) {
	return &groupRemoveCommand{&removeOptions{}}, nil
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

func (grc *groupRemoveCommand) Help() string {
	return `rrh group rm [OPTIONS] <GROUPS...>
OPTIONS
    -f, --force      force remove.
    -i, --inquiry    inquiry mode.
    -v, --verbose    verbose mode.
ARGUMENTS
    GROUPS           target group names.`
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
    list      list groups (default).
    rm        remove group.
    update    update group.`
}

/*
Run peforms the command.
*/
func (group *GroupCommand) Run(args []string) int {
	c := cli.NewCLI("rrh group", common.VERSION)
	c.Args = args
	c.Autocomplete = true
	c.Commands = map[string]cli.CommandFactory{
		"add":    groupAddCommandFactory,
		"update": groupUpdateCommandFactory,
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

type addOptions struct {
	desc string
	omit string
	args []string
}

func (gac *groupAddCommand) buildFlagSet() (*flag.FlagSet, *addOptions) {
	var opt = addOptions{}
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(gac.Help()) }
	flags.StringVar(&opt.desc, "d", "", "description")
	flags.StringVar(&opt.desc, "desc", "", "description")
	flags.StringVar(&opt.omit, "o", "", "omit list flag")
	flags.StringVar(&opt.omit, "omit-list", "", "omit list flag")
	return flags, &opt
}

func (gac *groupAddCommand) parse(args []string) (*addOptions, error) {
	var flags, opt = gac.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	opt.args = flags.Args()
	return opt, nil
}

/*
Run performs the command.
*/
func (gac *groupAddCommand) Run(args []string) int {
	var options, err = gac.parse(args)
	if err != nil {
		return 1
	}
	var config = common.OpenConfig()
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	if len(options.args) == 0 {
		fmt.Println(gac.Help())
		return 3
	}
	if err := gac.addGroups(db, options); err != nil {
		fmt.Println(err.Error())
		return 4
	}
	db.StoreAndClose()

	return 0
}

type listOptions struct {
	desc         bool
	repositories bool
	nameOnly     bool
}

func (glc *groupListCommand) buildFlagSet() (*flag.FlagSet, *listOptions) {
	var opt = listOptions{}
	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(glc.Help()) }
	flags.BoolVar(&opt.desc, "d", false, "show description")
	flags.BoolVar(&opt.desc, "desc", false, "show description")
	flags.BoolVar(&opt.repositories, "r", false, "show repositories")
	flags.BoolVar(&opt.repositories, "repository", false, "show repositories")
	flags.BoolVar(&opt.nameOnly, "o", false, "show only group names")
	flags.BoolVar(&opt.nameOnly, "only-groupname", false, "show only group names")
	return flags, &opt
}

func (glc *groupListCommand) parse(args []string) (*listOptions, error) {
	var flags, opt = glc.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	return opt, nil
}

func printRepositoryCount(count int) {
	if count == 1 {
		fmt.Print(",1 repository")
	} else {
		fmt.Printf(",%d repositories", count)
	}
}

func (glc *groupListCommand) printResult(result GroupResult, options *listOptions) {
	fmt.Print(result.Name)
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

func (glc *groupListCommand) printAll(results []GroupResult, options *listOptions) {
	for _, result := range results {
		glc.printResult(result, options)
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
	var config = common.OpenConfig()
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	var results = glc.listGroups(db, listOption)
	glc.printAll(results, listOption)

	return 0
}

type removeOptions struct {
	inquiry bool
	verbose bool
	force   bool
	args    []string
}

func (grc *groupRemoveCommand) printIfVerbose(message string) {
	if grc.Options.verbose {
		fmt.Println(message)
	}
}

func (grc *groupRemoveCommand) Inquiry(groupName string) bool {
	// no inquiry option, do remove group.
	if !grc.Options.inquiry {
		return true
	}
	return common.IsInputYes(fmt.Sprintf("%s: remove group? [yN]", groupName))
}

func (grc *groupRemoveCommand) buildFlagSet() (*flag.FlagSet, *removeOptions) {
	var opt = removeOptions{}
	flags := flag.NewFlagSet("rm", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(grc.Help()) }
	flags.BoolVar(&opt.inquiry, "i", false, "inquiry mode")
	flags.BoolVar(&opt.verbose, "v", false, "verbose mode")
	flags.BoolVar(&opt.force, "f", false, "force remove")
	flags.BoolVar(&opt.inquiry, "inquiry", false, "inquiry mode")
	flags.BoolVar(&opt.verbose, "verbose", false, "verbose mode")
	flags.BoolVar(&opt.force, "force", false, "force remove")
	return flags, &opt
}

func (grc *groupRemoveCommand) parse(args []string) (*removeOptions, error) {
	var flags, opt = grc.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	opt.args = flags.Args()
	if len(opt.args) == 0 {
		return nil, fmt.Errorf("no arguments are specified")
	}
	grc.Options = opt
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
	var config = common.OpenConfig()
	var db, err2 = common.Open(config)
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

type updateOptions struct {
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
	var config = common.OpenConfig()
	var db, err2 = common.Open(config)
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

func (guc *groupUpdateCommand) buildFlagSet() (*flag.FlagSet, *updateOptions) {
	var opt = updateOptions{}
	flags := flag.NewFlagSet("update", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(guc.Help()) }
	flags.StringVar(&opt.newName, "n", "", "specify new group name")
	flags.StringVar(&opt.newName, "name", "", "specify new group name")
	flags.StringVar(&opt.desc, "d", "", "specify the description")
	flags.StringVar(&opt.desc, "desc", "", "specify the description")
	flags.StringVar(&opt.omitList, "omit-list", "", "set the omit list flag.")
	flags.StringVar(&opt.omitList, "o", "", "set the omit list flag.")
	return flags, &opt
}

func (guc *groupUpdateCommand) parse(args []string) (*updateOptions, error) {
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
	return "add/list/update/remove groups."
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

/*
Synopsis returns the help message of the command.
*/
func (grc *groupRemoveCommand) Synopsis() string {
	return "remove given group."
}

/*
Synopsis returns the help message of the command.
*/
func (guc *groupUpdateCommand) Synopsis() string {
	return "update group."
}
