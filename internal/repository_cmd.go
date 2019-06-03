package internal

import (
	"fmt"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/lib"
)

/*
RepositoryCommand represents a command.
*/
type RepositoryCommand struct{}
type repositoryListCommand struct {
	options *repositoryListOptions
}
type repositoryInfoCommand struct {
	options *repositoryInfoOptions
}
type repositoryUpdateCommand struct {
	options *repositoryUpdateOptions
}

type repositoryInfoOptions struct {
	color   bool
	csv     bool
	noColor bool
	args    []string
}

type repositoryListOptions struct {
	id    bool
	path  bool
	group bool
	args  []string
}

type repositoryUpdateOptions struct {
	repositoryID string
	newID        string
	description  string
	newPath      string
}

func repositoryInfoCommandFactory() (cli.Command, error) {
	return &repositoryInfoCommand{}, nil
}

func repositoryListCommandFactory() (cli.Command, error) {
	return &repositoryListCommand{}, nil
}

func repositoryUpdateCommandFactory() (cli.Command, error) {
	return &repositoryUpdateCommand{}, nil
}

/*
RepositoryCommandFactory returns an instance of the PruneCommand.
*/
func RepositoryCommandFactory() (cli.Command, error) {
	return &RepositoryCommand{}, nil
}

func (info *repositoryInfoCommand) buildFlagSet() (*flag.FlagSet, *repositoryInfoOptions) {
	var options = repositoryInfoOptions{}
	flags := flag.NewFlagSet("info", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(info.Help()) }
	flags.BoolVarP(&options.csv, "csv", "c", false, "prints in the csv format.")
	flags.BoolVarP(&options.color, "color", "G", false, "enables colorized output.")
	flags.BoolVar(&options.noColor, "no-color", false, "no colorized output.")
	return flags, &options
}

func (info *repositoryInfoCommand) parseOptions(args []string) error {
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

func (options *repositoryInfoOptions) printInfo(result lib.Repository, config *lib.Config) {
	fmt.Printf("%-12s %s\n", config.Color.ColorizedLabel("ID:"), config.Color.ColorizedRepositoryID(result.ID))
	fmt.Printf("%-12s %s\n", config.Color.ColorizedLabel("Description:"), result.Description)
	fmt.Printf("%-12s %s\n", config.Color.ColorizedLabel("Path:"), result.Path)
	if len(result.Remotes) > 0 {
		printRemoteInfo(result.Remotes, config)
	}
}

func printRemoteInfo(remotes []lib.Remote, config *lib.Config) {
	fmt.Printf("%-12s\n", config.Color.ColorizedLabel("Remote:"))
	for _, remote := range remotes {
		fmt.Printf("    %s: %s\n", config.Color.ColorizedLabel(remote.Name), remote.URL)
	}
}

func (options *repositoryInfoOptions) printInfoResult(result lib.Repository, config *lib.Config) {
	if options.csv {
		fmt.Printf("%s,%s,%s\n", config.Color.ColorizedRepositoryID(result.ID), result.Description, result.Path)
	} else {
		options.printInfo(result, config)
	}
}

func (info *repositoryInfoCommand) perform(db *lib.Database, args []string) int {
	var results, errs = findResults(db, args)
	var onError = db.Config.GetValue(lib.RrhOnError)
	for _, result := range results {
		info.options.printInfoResult(result, db.Config)
	}
	if len(errs) > 0 && onError != lib.Ignore {
		return printErrors(db.Config, errs)
	}
	return 0
}

func (info *repositoryInfoCommand) Run(args []string) int {
	var err = info.parseOptions(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var config = lib.OpenConfig()
	var db, err2 = lib.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	config.Color.SetColorize(info.options.color || !info.options.noColor)
	return info.perform(db, info.options.args)
}

func printListGroup(db *lib.Database, result lib.Repository) {
	var groups = db.FindRelationsOfRepository(result.ID)
	for _, group := range groups {
		fmt.Printf("%s/%s\n", group, result.ID)
	}
}

func printListResult(db *lib.Database, result lib.Repository, options *repositoryListOptions) {
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

func (list *repositoryListCommand) perform(db *lib.Database, args []string) int {
	var results, errs = findAll(db, args)
	var onError = db.Config.GetValue(lib.RrhOnError)
	for _, result := range results {
		printListResult(db, result, list.options)
	}
	if len(errs) > 0 && onError != lib.Ignore {
		return printErrors(db.Config, errs)
	}
	return 0
}

func (list *repositoryListCommand) buildFlagSet() (*flag.FlagSet, *repositoryListOptions) {
	var options = repositoryListOptions{}
	flags := flag.NewFlagSet("list", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(list.Help()) }
	flags.BoolVar(&options.id, "id", false, "prints id of the repository.")
	flags.BoolVar(&options.path, "path", false, "prints path of the repository.")
	flags.BoolVar(&options.group, "with-group", false, "prints group of the repository.")
	return flags, &options
}

func (list *repositoryListCommand) parseOptions(args []string) error {
	var flags, options = list.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return err
	}
	options.args = flags.Args()
	list.options = options
	return nil
}

func (list *repositoryListCommand) Run(args []string) int {
	var err = list.parseOptions(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var config = lib.OpenConfig()
	var db, err2 = lib.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 2
	}
	return list.perform(db, list.options.args)
}

func (update *repositoryUpdateCommand) buildFlagSet() (*flag.FlagSet, *repositoryUpdateOptions) {
	var options = repositoryUpdateOptions{}
	flags := flag.NewFlagSet("update", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(update.Help()) }
	flags.StringVarP(&options.newID, "id", "i", "", "specifies new repository id")
	flags.StringVarP(&options.description, "desc", "d", "", "specifies description")
	flags.StringVarP(&options.newPath, "path", "p", "", "specifies new path")
	return flags, &options
}

func (update *repositoryUpdateCommand) parseOptions(args []string) error {
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

func (update *repositoryUpdateCommand) Run(args []string) int {
	var err = update.parseOptions(args)
	if err != nil {
		fmt.Printf(update.Help())
		return 1
	}
	var config = lib.OpenConfig()
	var db, err2 = lib.Open(config)
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
func (repository *RepositoryCommand) Run(args []string) int {
	c := cli.NewCLI("rrh repository", lib.VERSION)
	c.Args = args
	c.Autocomplete = true
	c.Commands = map[string]cli.CommandFactory{
		"list":   repositoryListCommandFactory,
		"info":   repositoryInfoCommandFactory,
		"update": repositoryUpdateCommandFactory,
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
func (repository *RepositoryCommand) Help() string {
	return `rrh repository <SUBCOMMAND>
SUBCOMMAND
    info [OPTIONS] <REPO...>     shows repository information.
    update [OPTIONS] <REPO...>   updates repository information.`
}

func (info *repositoryInfoCommand) Help() string {
	return `rrh repository info [OPTIONS] [REPOSITORIES...]
    -G, --color     prints the results with color.
    -c, --csv       prints the results in the csv format.
ARGUMENTS
    REPOSITORIES    target repositories.  If no repositories are specified,
                    this sub command failed.`
}

func (list *repositoryListCommand) Help() string {
	return `rrh repository list [OPTIONS] [ARGUMENTS...]
OPTIONS
    --id            prints ids in the results.
    --path          prints paths in the results.
    --with-group    prints the results in "GROUP/REPOSITORY" format.
Note:
    This sub command is used for a completion target generation.`
}

func (update *repositoryUpdateCommand) Help() string {
	return `rrh repository update [OPTIONS] <REPOSITORY>
OPTIONS
    -i, --id <NEWID>     specifies new repository id.
    -d, --desc <DESC>    specifies new description.
    -p, --path <PATH>    specifies new path.
ARGUMENTS
    REPOSITORY           specifies the repository id.`
}

func (info *repositoryInfoCommand) Synopsis() string {
	return "prints information of the specified repositories."
}

func (list *repositoryListCommand) Synopsis() string {
	return "lists repositories."
}

func (update *repositoryUpdateCommand) Synopsis() string {
	return "update information of the specified repository."
}

/*
Synopsis returns the help message of the command.
*/
func (repository *RepositoryCommand) Synopsis() string {
	return "manages repositories."
}

func findAll(db *lib.Database, args []string) ([]lib.Repository, []error) {
	if len(args) > 0 {
		return findResults(db, args)
	}
	return db.Repositories, []error{}
}

func findResults(db *lib.Database, args []string) ([]lib.Repository, []error) {
	var results = []lib.Repository{}
	var errs = []error{}
	for _, arg := range args {
		var repo = db.FindRepository(arg)
		if repo == nil {
			errs = append(errs, fmt.Errorf("%s: repository not found", arg))
			if db.Config.GetValue(lib.RrhOnError) == lib.FailImmediately {
				return []lib.Repository{}, errs
			}
		} else {
			results = append(results, *repo)
		}
	}
	return results, errs
}

func (update *repositoryUpdateCommand) perform(db *lib.Database, targetRepoID string) error {
	var repo = db.FindRepository(targetRepoID)
	if repo == nil {
		return fmt.Errorf("%s: repository not found", targetRepoID)
	}
	var newRepo = buildNewRepo(update.options, repo)
	if !db.UpdateRepository(targetRepoID, newRepo) {
		return fmt.Errorf("%s: repository update failed", targetRepoID)
	}
	return nil
}

func buildNewRepo(options *repositoryUpdateOptions, repo *lib.Repository) lib.Repository {
	var newRepo = lib.Repository{ID: repo.ID, Path: repo.Path, Description: repo.Description}
	if options.description != "" {
		newRepo.Description = options.description
	}
	if options.newID != "" {
		newRepo.ID = options.newID
	}
	if options.newPath != "" {
		newRepo.Path = options.newPath
	}
	return newRepo
}
