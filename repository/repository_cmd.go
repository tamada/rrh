package repository

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
type listCommand struct {
	options *listOptions
}
type infoCommand struct {
	options *infoOptions
}
type updateCommand struct {
	options *updateOptions
}

type infoOptions struct {
	color   bool
	csv     bool
	noColor bool
	args    []string
}

type listOptions struct {
	id    bool
	path  bool
	group bool
	args  []string
}

type updateOptions struct {
	repositoryID string
	newID        string
	description  string
	newPath      string
}

func printErrors(config *common.Config, errs []error) int {
	var onError = config.GetValue(common.RrhOnError)
	for _, err := range errs {
		fmt.Println(err.Error())
	}
	if onError == common.Fail || onError == common.FailImmediately {
		return 1
	}
	return 0
}

func infoCommandFactory() (cli.Command, error) {
	return &infoCommand{}, nil
}

func listCommandFactory() (cli.Command, error) {
	return &listCommand{}, nil
}

func updateCommandFactory() (cli.Command, error) {
	return &updateCommand{}, nil
}

/*
CommandFactory returns an instance of the PruneCommand.
*/
func CommandFactory() (cli.Command, error) {
	return &Command{}, nil
}

func (info *infoCommand) buildFlagSet() (*flag.FlagSet, *infoOptions) {
	var options = infoOptions{}
	flags := flag.NewFlagSet("info", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(info.Help()) }
	flags.BoolVar(&options.csv, "c", false, "prints in the csv format.")
	flags.BoolVar(&options.csv, "csv", false, "prints in the csv format.")
	flags.BoolVar(&options.color, "G", false, "enables colorized output.")
	flags.BoolVar(&options.color, "color", false, "enables colorized output.")
	flags.BoolVar(&options.noColor, "no-color", false, "no colorized output.")
	return flags, &options
}

func (info *infoCommand) parseOptions(args []string) error {
	var flags, options = info.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return err
	}
	options.args = flags.Args()
	info.options = options
	if len(options.args) == 0 {
		return fmt.Errorf("missing arguments")
	}
	return nil
}

func printInfo(result common.Repository, options *infoOptions) {
	fmt.Printf("%-12s %s\n", common.ColorizedLabel("ID:"), common.ColorizedRepositoryID(result.ID))
	fmt.Printf("%-12s %s\n", common.ColorizedLabel("Description:"), result.Description)
	fmt.Printf("%-12s %s\n", common.ColorizedLabel("Path:"), result.Path)
	if len(result.Remotes) > 0 {
		fmt.Printf("%-12s\n", common.ColorizedLabel("Remote:"))
		for _, remote := range result.Remotes {
			fmt.Printf("    %s: %s\n", common.ColorizedLabel(remote.Name), remote.URL)
		}
	}
}

func printInfoResult(result common.Repository, options *infoOptions) {
	if options.csv {
		fmt.Printf("%s,%s,%s\n", common.ColorizedRepositoryID(result.ID), result.Description, result.Path)
	} else {
		printInfo(result, options)
	}
}

func (info *infoCommand) perform(db *common.Database, args []string) int {
	var results, errs = findResults(db, args)
	var onError = db.Config.GetValue(common.RrhOnError)
	for _, result := range results {
		printInfoResult(result, info.options)
	}
	if len(errs) > 0 && onError != common.Ignore {
		return printErrors(db.Config, errs)
	}
	return 0
}

func (info *infoCommand) Run(args []string) int {
	var err = info.parseOptions(args)
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
	common.SetColorize(info.options.color || !info.options.noColor)
	return info.perform(db, info.options.args)
}

func printListGroup(db *common.Database, result common.Repository) {
	var groups = db.FindRelationsOfRepository(result.ID)
	for _, group := range groups {
		fmt.Printf("%s/%s\n", group, result.ID)
	}
}

func printListResult(db *common.Database, result common.Repository, options *listOptions) {
	if options.group {
		printListGroup(db, result)
	}
	if options.id {
		fmt.Println(result.ID)
	}
	if options.path {
		fmt.Println(result.Path)
	}
}

func (list *listCommand) perform(db *common.Database, args []string) int {
	var results, errs = findAll(db, args)
	var onError = db.Config.GetValue(common.RrhOnError)
	for _, result := range results {
		printListResult(db, result, list.options)
	}
	if len(errs) > 0 && onError != common.Ignore {
		return printErrors(db.Config, errs)
	}
	return 0
}

func (list *listCommand) buildFlagSet() (*flag.FlagSet, *listOptions) {
	var options = listOptions{}
	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(list.Help()) }
	flags.BoolVar(&options.id, "id", false, "prints id of the repository.")
	flags.BoolVar(&options.group, "group", false, "prints group of the repository.")
	flags.BoolVar(&options.path, "path", false, "prints path of the repository.")
	return flags, &options
}

func (list *listCommand) parseOptions(args []string) error {
	var flags, options = list.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return err
	}
	options.args = flags.Args()
	list.options = options
	return nil
}

func (list *listCommand) Run(args []string) int {
	var err = list.parseOptions(args)
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
	return list.perform(db, list.options.args)
}

func (update *updateCommand) buildFlagSet() (*flag.FlagSet, *updateOptions) {
	var options = updateOptions{}
	flags := flag.NewFlagSet("update", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(update.Help()) }
	flags.StringVar(&options.newID, "i", "", "specifies new repository id")
	flags.StringVar(&options.description, "d", "", "specifies description")
	flags.StringVar(&options.newPath, "p", "", "specifies new path")
	flags.StringVar(&options.newID, "id", "", "specifies new repository id")
	flags.StringVar(&options.description, "desc", "", "specifies description")
	flags.StringVar(&options.newPath, "path", "", "specifies new path")
	return flags, &options
}

func (update *updateCommand) parseOptions(args []string) error {
	var flags, options = update.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return err
	}
	var arguments = flags.Args()
	if len(arguments) == 0 {
		return fmt.Errorf("missing arguments")
	} else if len(arguments) > 1 {
		return fmt.Errorf("too many arguments: %v", arguments)
	}
	options.repositoryID = arguments[0]
	update.options = options
	return nil
}

func (update *updateCommand) Run(args []string) int {
	var err = update.parseOptions(args)
	if err != nil {
		fmt.Printf(update.Help())
		return 1
	}
	var config = common.OpenConfig()
	var db, err2 = common.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	var err3 = update.perform(db, update.options.repositoryID)
	if err3 != nil {
		fmt.Println(err3.Error())
		return 3
	}
	db.StoreAndClose()
	return 0
}

/*
Run performs the command.
*/
func (repository *Command) Run(args []string) int {
	c := cli.NewCLI("rrh repository", common.VERSION)
	c.Args = args
	c.Autocomplete = true
	c.Commands = map[string]cli.CommandFactory{
		"list":   listCommandFactory,
		"info":   infoCommandFactory,
		"update": updateCommandFactory,
	}
	if len(args) == 0 {
		fmt.Println(repository.Help())
		return 0
	}
	var exitStatus, err = c.Run()
	if err != nil {
		fmt.Println(err.Error())
	}
	return exitStatus
}

/*
Help function shows the help message.
*/
func (repository *Command) Help() string {
	return `rrh repository <SUBCOMMAND>
SUBCOMMAND
    info [OPTIONS] <REPO...>     shows repository information.
    update [OPTIONS] <REPO...>   updates repository information.`
}

func (info *infoCommand) Help() string {
	return `rrh repository info [OPTIONS] [REPOSITORIES...]
    -G, --color     prints the results with color.
    -c, --csv       prints the results in the csv format.
ARGUMENTS
    REPOSITORIES    target repositories.  If no repositories are specified,
                    this sub command failed.`
}

func (list *listCommand) Help() string {
	return `rrh repository list [OPTIONS] [ARGUMENTS...]
OPTIONS
    -i, --id       prints ids in the results.
    -p, --path     prints paths in the results.
    -g, --group    prints the results in "GROUP/REPOSITORY" format.
Note:
    This sub command is used for a completion target generation.`
}

func (update *updateCommand) Help() string {
	return `rrh repository update [OPTIONS] <REPOSITORY>
OPTIONS
    -i, --id <NEWID>     specifies new repository id.
    -d, --desc <DESC>    specifies new description.
    -p, --path <PATH>    specifies new path.
ARGUMENTS
    REPOSITORY           specifies the repository id.`
}

func (info *infoCommand) Synopsis() string {
	return "prints information of the specified repositories."
}

func (list *listCommand) Synopsis() string {
	return "lists repositories."
}

func (update *updateCommand) Synopsis() string {
	return "update information of the specified repository."
}

/*
Synopsis returns the help message of the command.
*/
func (repository *Command) Synopsis() string {
	return "manages repositories."
}
