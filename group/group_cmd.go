package group

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

type GroupCommand struct{}
type groupAddCommand struct{}
type groupListCommand struct{}
type groupUpdateCommand struct{}
type groupRemoveCommand struct{}

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
	return &groupRemoveCommand{}, nil
}

func (group *groupAddCommand) Help() string {
	return `rrh group add [OPTIONS] <GROUPS...>
OPTIONS
    -d, --desc <DESC>    give the description of the group
ARGUMENTS
    GROUP                gives group names.`
}

func (group *groupListCommand) Help() string {
	return `rrh group list [OPTIONS]
OPTIONS
    -d, --desc          show description.
    -r, --repository    show repositories in the group.`
}

func (group *groupRemoveCommand) Help() string {
	return `rrh group rm [OPTIONS] <GROUPS...>
OPTIONS
    -f, --force      force remove
	-i, --inquery    inquiry mode
	-v, --verbose    verbose mode
ARGUMENTS
    GROUPS           target group names.`
}

func (group *groupUpdateCommand) Help() string {
	return `rrh group update [OPTIONS] <GROUP>
OPTIONS
    -n, --name <NAME>   change group name to NAME.
    -d, --desc <DESC>   change description to DESC.
ARGUMENTS
    GROUP               update target group names.`
}

func (group *GroupCommand) Help() string {
	return `rrh group <SUBCOMMAND>
SUBCOMMAND
    add       add new group.
    list      list groups (default).
    rm        remove group.
    update    update group`
}

func (group *GroupCommand) Run(args []string) int {
	c := cli.NewCLI("rrh group", "1.0.0")
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
	} else {
		var exitStatus, err = c.Run()
		if err != nil {
			fmt.Println(err.Error())
		}
		return exitStatus
	}
}

type addOptions struct {
	desc string
	args []string
}

func (group *groupAddCommand) parse(args []string) (*addOptions, error) {
	var opt = addOptions{}
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(group.Help()) }
	flags.StringVar(&opt.desc, "d", "", "description")
	flags.StringVar(&opt.desc, "desc", "", "description")
	opt.args = flags.Args()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	return &opt, nil
}

/*
Run performs the command.
*/
func (group *groupAddCommand) Run(args []string) int {
	var options, err = group.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(group.Help())
		return 1
	}
	var config = common.OpenConfig()
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	if err := group.addGroups(db, options); err != nil {
		fmt.Println(err.Error())
		return 3
	}
	db.StoreAndClose()

	return 0
}

type listOptions struct {
	desc         bool
	repositories bool
}

func (group *groupListCommand) parse(args []string) (*listOptions, error) {
	var opt = listOptions{}
	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(group.Help()) }
	flags.BoolVar(&opt.desc, "d", false, "show description")
	flags.BoolVar(&opt.desc, "desc", false, "show description")
	flags.BoolVar(&opt.repositories, "r", false, "show repositories")
	flags.BoolVar(&opt.repositories, "repository", false, "show repositories")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	return &opt, nil
}

func (group *groupListCommand) printAll(results []GroupResult, options *listOptions) {
	for _, result := range results {
		fmt.Printf("%s,", result.Name)
		if options.desc {
			fmt.Printf("%s,", result.Description)
		}
		if options.repositories {
			fmt.Printf("%v,", result.Repos)
		}
		if len(result.Repos) == 1 {
			fmt.Println("(1 repository)")
		} else {
			fmt.Printf("(%d repositories)\n", len(result.Repos))
		}
	}
}

/*
Run performs the command.
*/
func (group *groupListCommand) Run(args []string) int {
	var listOption, err = group.parse(args)
	if err != nil {
		return 1
	}
	var config = common.OpenConfig()
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	var results, err3 = group.listGroups(db, listOption)
	if err3 != nil {
		fmt.Println(err3.Error())
	}
	group.printAll(results, listOption)

	return 0
}

type removeOptions struct {
	inquiry bool
	verbose bool
	force   bool
	args    []string
}

func (options *removeOptions) printIfVerbose(message string) {
	if options.verbose {
		fmt.Println(message)
	}
}

func (options *removeOptions) Inquiry(groupName string) bool {
	// no inquiry option, do remove group.
	if !options.inquiry {
		return true
	}
	return common.IsInputYes(fmt.Sprintf("%s: remove group? [yN]", groupName))
}

func (group *groupRemoveCommand) parse(args []string) (*removeOptions, error) {
	var opt = removeOptions{}
	flags := flag.NewFlagSet("rm", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(group.Help()) }
	flags.BoolVar(&opt.inquiry, "i", false, "inquiry mode")
	flags.BoolVar(&opt.verbose, "v", false, "verbose mode")
	flags.BoolVar(&opt.force, "f", false, "force remove")
	opt.args = flags.Args()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if len(opt.args) == 0 {
		return nil, fmt.Errorf("no arguments are specified")
	}
	return &opt, nil
}

/*
Run performs the command.
*/
func (group *groupRemoveCommand) Run(args []string) int {
	var options, err = group.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(group.Help())
		return 1
	}
	var config = common.OpenConfig()
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	if err := group.removeGroups(db, options); err != nil {
		fmt.Println(err.Error())
		return 3
	}
	db.StoreAndClose()

	return 0
}

type updateOptions struct {
	newName string
	desc    string
	target  string
}

/*
Run performs the command.
*/
func (group *groupUpdateCommand) Run(args []string) int {
	var updateOption, err = group.parse(args)
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
	var err3 = group.updateGroup(db, updateOption)
	if err3 != nil {
		fmt.Println(err3.Error())
		return 3
	}
	db.StoreAndClose()
	return 0
}

func (group *groupUpdateCommand) parse(args []string) (*updateOptions, error) {
	var opt = updateOptions{}
	flags := flag.NewFlagSet("update", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(group.Help()) }
	flags.StringVar(&opt.newName, "n", "", "show description")
	flags.StringVar(&opt.newName, "name", "", "show description")
	flags.StringVar(&opt.desc, "d", "", "show repositories")
	flags.StringVar(&opt.desc, "desc", "", "show repositories")
	var arguments = flags.Args()

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	if len(arguments) == 0 {
		return nil, fmt.Errorf("no arguments are specified")
	}
	if len(arguments) > 1 {
		return nil, fmt.Errorf("could not accept multiple arguments")
	}
	opt.target = arguments[0]
	return &opt, nil
}

func (group *GroupCommand) Synopsis() string {
	return "add/list/update/remove groups."
}

func (group *groupAddCommand) Synopsis() string {
	return "add group."
}

func (group *groupListCommand) Synopsis() string {
	return "list groups."
}

func (group *groupRemoveCommand) Synopsis() string {
	return "remove given group."
}

func (group *groupUpdateCommand) Synopsis() string {
	return "update group."
}
