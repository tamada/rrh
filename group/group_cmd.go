package group

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

/*
Command represents a command.
*/
type Command struct{}
type addCommand struct{}
type listCommand struct{}
type ofCommand struct{}
type updateCommand struct{}
type removeCommand struct {
	options *removeOptions
}

/*
CommandFactory returns an instance of command.
*/
func CommandFactory() (cli.Command, error) {
	return &Command{}, nil
}

func addCommandFactory() (cli.Command, error) {
	return &addCommand{}, nil
}

func ofCommandFactory() (cli.Command, error) {
	return &ofCommand{}, nil
}

func listCommandFactory() (cli.Command, error) {
	return &listCommand{}, nil
}

func updateCommandFactory() (cli.Command, error) {
	return &updateCommand{}, nil
}

func removeCommandFactory() (cli.Command, error) {
	return &removeCommand{&removeOptions{}}, nil
}

func (gac *addCommand) Help() string {
	return `rrh group add [OPTIONS] <GROUPS...>
OPTIONS
    -d, --desc <DESC>        gives the description of the group.
    -o, --omit-list <FLAG>   gives the omit list flag of the group.
ARGUMENTS
    GROUPS                   gives group names.`
}

func (glc *listCommand) Help() string {
	return `rrh group list [OPTIONS]
OPTIONS
    -d, --desc             show description.
    -r, --repository       show repositories in the group.
    -o, --only-groupname   show only group name. This option is prioritized.`
}

func (goc *ofCommand) Help() string {
	return `rrh group of <REPOSITORY_ID>
ARGUMENTS
    REPOSITORY_ID     show the groups of the repository.`
}

func (grc *removeCommand) Help() string {
	return `rrh group rm [OPTIONS] <GROUPS...>
OPTIONS
    -f, --force      force remove.
    -i, --inquiry    inquiry mode.
    -v, --verbose    verbose mode.
ARGUMENTS
    GROUPS           target group names.`
}

func (guc *updateCommand) Help() string {
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
func (group *Command) Help() string {
	return `rrh group <SUBCOMMAND>
SUBCOMMAND
    add       add new group.
    list      list groups (default).
    of        shows groups of the specified repository.
    rm        remove group.
    update    update group.`
}

/*
Run peforms the command.
*/
func (group *Command) Run(args []string) int {
	c := cli.NewCLI("rrh group", common.VERSION)
	c.Args = args
	c.Autocomplete = true
	c.Commands = map[string]cli.CommandFactory{
		"add":    addCommandFactory,
		"update": updateCommandFactory,
		"of":     ofCommandFactory,
		"rm":     removeCommandFactory,
		"list":   listCommandFactory,
	}
	if len(args) == 0 {
		new(listCommand).Run([]string{})
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

func (gac *addCommand) buildFlagSet() (*flag.FlagSet, *addOptions) {
	var opt = addOptions{}
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(gac.Help()) }
	flags.StringVar(&opt.desc, "d", "", "description")
	flags.StringVar(&opt.desc, "desc", "", "description")
	flags.StringVar(&opt.omit, "o", "", "omit list flag")
	flags.StringVar(&opt.omit, "omit-list", "", "omit list flag")
	return flags, &opt
}

func (gac *addCommand) parse(args []string) (*addOptions, error) {
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
func (gac *addCommand) Run(args []string) int {
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

func (glc *listCommand) buildFlagSet() (*flag.FlagSet, *listOptions) {
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

func (glc *listCommand) parse(args []string) (*listOptions, error) {
	var flags, opt = glc.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	return opt, nil
}

func (goc *ofCommand) perform(db *common.Database, repositoryID string) int {
	if !db.HasRepository(repositoryID) {
		fmt.Printf("%s: repository not found\n", repositoryID)
		return 3
	}
	var groups = db.FindRelationsOfRepository(repositoryID)
	fmt.Printf("%s, %v\n", repositoryID, groups)
	return 0
}

func (goc *ofCommand) Run(args []string) int {
	if len(args) != 1 {
		fmt.Println(goc.Help())
		return 1
	}
	var config = common.OpenConfig()
	var db, err = common.Open(config)
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

func findGroupName(name string, nameOnlyFlag bool) string {
	if nameOnlyFlag {
		return name
	}
	return common.ColorrizedGroupName(name)
}

func (glc *listCommand) printResult(result Result, options *listOptions) {
	fmt.Print(findGroupName(result.Name, options.nameOnly))
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

func (glc *listCommand) printAll(results []Result, options *listOptions) {
	for _, result := range results {
		glc.printResult(result, options)
	}
}

/*
Run performs the command.
*/
func (glc *listCommand) Run(args []string) int {
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

func (grc *removeCommand) printIfVerbose(message string) {
	if grc.options.verbose {
		fmt.Println(message)
	}
}

func (grc *removeCommand) Inquiry(groupName string) bool {
	// no inquiry option, do remove group.
	if !grc.options.inquiry {
		return true
	}
	return common.IsInputYes(fmt.Sprintf("%s: remove group? [yN]", groupName))
}

func (grc *removeCommand) buildFlagSet() (*flag.FlagSet, *removeOptions) {
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

func (grc *removeCommand) parse(args []string) (*removeOptions, error) {
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
func (grc *removeCommand) Run(args []string) int {
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
func (guc *updateCommand) Run(args []string) int {
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

func (guc *updateCommand) buildFlagSet() (*flag.FlagSet, *updateOptions) {
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

func (guc *updateCommand) parse(args []string) (*updateOptions, error) {
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
func (group *Command) Synopsis() string {
	return "add/list/update/remove groups and show groups of the repository."
}

/*
Synopsis returns the help message of the command.
*/
func (gac *addCommand) Synopsis() string {
	return "add group."
}

/*
Synopsis returns the help message of the command.
*/
func (glc *listCommand) Synopsis() string {
	return "list groups."
}

func (goc *ofCommand) Synopsis() string {
	return "show groups of the repository."
}

/*
Synopsis returns the help message of the command.
*/
func (grc *removeCommand) Synopsis() string {
	return "remove given group."
}

/*
Synopsis returns the help message of the command.
*/
func (guc *updateCommand) Synopsis() string {
	return "update group."
}
