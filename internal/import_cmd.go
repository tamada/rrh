package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh"
)

type importOptions struct {
	overwrite bool
	autoClone bool
	verbose   bool
	database  string
}

/*
ImportCommand represents a command.
*/
type ImportCommand struct {
	options *importOptions
}

/*
ImportCommandFactory generate the command struct.
*/
func ImportCommandFactory() (cli.Command, error) {
	return &ImportCommand{&importOptions{}}, nil
}

func (options *importOptions) printIfNeeded(message string) {
	if options.verbose {
		fmt.Println(message)
	}
}

func eraseDatabase(db *rrh.Database, command *ImportCommand) {
	db.Groups = []*rrh.Group{}
	db.Repositories = []*rrh.Repository{}
	db.Relations = []*rrh.Relation{}
	command.options.printIfNeeded("The local database is cleared")
}

func perform(db *rrh.Database, command *ImportCommand) int {
	if command.options.overwrite {
		eraseDatabase(db, command)
	}
	var db2, err = readNewDB(command.options.database, db.Config)
	if err != nil {
		fmt.Printf(err.Error())
		return 4
	}
	var errs = command.copyDB(db2, db)
	var statusCode = db.Config.PrintErrors(errs...)
	if statusCode == 0 {
		db.StoreAndClose()
	}
	return statusCode
}

/*
Run peforms the command.
*/
func (command *ImportCommand) Run(args []string) int {
	var err1 = parse(args, command)
	if err1 != nil {
		fmt.Println(err1)
		return 1
	}
	var config = rrh.OpenConfig()
	var db, err2 = rrh.Open(config)
	if err2 != nil {
		return 2
	}
	return perform(db, command)
}

func (command *ImportCommand) buildFlagSet() (*flag.FlagSet, *importOptions) {
	var options = importOptions{false, false, false, ""}
	var flags = flag.NewFlagSet("import", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(command.Help()) }
	flags.BoolVar(&options.overwrite, "overwrite", false, "overwrite mode")
	flags.BoolVar(&options.autoClone, "auto-clone", false, "auto clone mode")
	flags.BoolVarP(&options.verbose, "verbose", "v", false, "verbose mode")
	return flags, &options
}

func parse(args []string, command *ImportCommand) error {
	var flags, options = command.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return err
	}
	var arguments = flags.Args()
	if len(arguments) == 0 {
		return fmt.Errorf("too few arguments")
	} else if len(arguments) > 1 {
		return fmt.Errorf("too many arguments: %v", arguments)
	}
	options.database = arguments[0]
	command.options = options
	return nil
}

/*
Synopsis returns the simple help message of the command.
*/
func (command *ImportCommand) Synopsis() string {
	return "import the given database."
}

/*
Help returns the help message of the command.
*/
func (command *ImportCommand) Help() string {
	return `rrh import [OPTIONS] <DATABASE_JSON>
OPTIONS
    --auto-clone    clone the repository, if paths do not exist.
    --overwrite     replace the local RRH database to the given database.
    -v, --verbose   verbose mode.
ARGUMENTS
    DATABASE_JSON   the exported RRH database.`
}

func readNewDB(path string, config *rrh.Config) (*rrh.Database, error) {
	var db = rrh.Database{Timestamp: rrh.Now(), Repositories: []*rrh.Repository{}, Groups: []*rrh.Group{}, Relations: []*rrh.Relation{}, Config: config}
	var bytes, err = ioutil.ReadFile(path)
	if err != nil {
		return &db, nil
	}
	var homeReplacedString = replaceHome(bytes)

	if err := json.Unmarshal([]byte(homeReplacedString), &db); err != nil {
		return nil, err
	}
	return &db, nil
}

func (command *ImportCommand) copyDB(from *rrh.Database, to *rrh.Database) []error {
	var errs = []error{}
	var errs1 = command.copyGroups(from, to)
	var errs2 = command.copyRepositories(from, to)
	var errs3 = command.copyRelations(from, to)
	errs = append(errs, errs1...)
	errs = append(errs, errs2...)
	return append(errs, errs3...)
}

func (command *ImportCommand) copyGroup(group *rrh.Group, to *rrh.Database) []error {
	var list = []error{}
	if to.HasGroup(group.Name) {
		var successFlag = to.UpdateGroup(group.Name, group)
		if !successFlag {
			list = append(list, fmt.Errorf("%s: update failed", group.Name))
		}
	} else {
		var _, err = to.CreateGroup(group.Name, group.Description, group.OmitList)
		if err != nil {
			list = append(list, err)
		}
		command.options.printIfNeeded(fmt.Sprintf("%s: create group", group.Name))
	}
	return list
}

func (command *ImportCommand) copyGroups(from *rrh.Database, to *rrh.Database) []error {
	var list = []error{}
	for _, group := range from.Groups {
		var errs = command.copyGroup(group, to)
		list = append(list, errs...)
		if len(errs) != 0 && isFailImmediately(from.Config) {
			return list
		}
	}
	return list
}

func findOrigin(remotes []*rrh.Remote) *rrh.Remote {
	for _, remote := range remotes {
		if remote.Name == "origin" {
			return remote
		}
	}
	return remotes[0]
}

func doClone(repository *rrh.Repository, remote *rrh.Remote) error {
	var cmd = exec.Command("git", "clone", remote.URL, repository.Path)
	var err = cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: clone error (%s)", remote.URL, err.Error())
	}
	return nil
}

func (command *ImportCommand) cloneRepository(repository *rrh.Repository) error {
	if len(repository.Remotes) == 0 {
		return fmt.Errorf("%s: could not clone, did not have remotes", repository.ID)
	}
	var remote = findOrigin(repository.Remotes)
	var err = doClone(repository, remote)
	command.options.printIfNeeded(fmt.Sprintf("%s: clone repository from %s", repository.ID, remote.URL))
	return err
}

func (command *ImportCommand) cloneIfNeeded(repository *rrh.Repository) error {
	if !command.options.autoClone {
		return fmt.Errorf("%s: repository path did not exist at %s", repository.ID, repository.Path)
	}
	command.cloneRepository(repository)
	return nil
}

func (command *ImportCommand) copyRepository(repository *rrh.Repository, to *rrh.Database) []error {
	if to.HasRepository(repository.ID) {
		return []error{}
	}
	var _, err = os.Stat(repository.Path)
	if err != nil {
		var err1 = command.cloneIfNeeded(repository)
		if err1 != nil {
			return []error{err1}
		}
	}
	return command.copyRepositoryImpl(repository, to)
}

func (command *ImportCommand) copyRepositoryImpl(repository *rrh.Repository, to *rrh.Database) []error {
	if err := rrh.IsExistAndGitRepository(repository.Path, repository.ID); err != nil {
		return []error{err}
	}
	to.CreateRepository(repository.ID, repository.Path, repository.Description, repository.Remotes)
	command.options.printIfNeeded(fmt.Sprintf("%s: create repository", repository.ID))
	return []error{}
}

func (command *ImportCommand) copyRepositories(from *rrh.Database, to *rrh.Database) []error {
	var list = []error{}
	for _, repository := range from.Repositories {
		var errs = command.copyRepository(repository, to)
		list = append(list, errs...)
		if len(errs) > 0 && isFailImmediately(from.Config) {
			return list
		}
	}
	return list
}

func (command *ImportCommand) copyRelation(rel *rrh.Relation, to *rrh.Database) []error {
	var list = []error{}
	if to.HasGroup(rel.GroupName) && to.HasRepository(rel.RepositoryID) {
		to.Relate(rel.GroupName, rel.RepositoryID)
		command.options.printIfNeeded(fmt.Sprintf("%s, %s: create relation", rel.GroupName, rel.RepositoryID))
	} else {
		list = append(list, fmt.Errorf("group %s and repository %s: could not relate", rel.GroupName, rel.RepositoryID))
	}
	return list
}

func (command *ImportCommand) copyRelations(from *rrh.Database, to *rrh.Database) []error {
	var list = []error{}
	for _, rel := range from.Relations {
		var errs = command.copyRelation(rel, to)
		list = append(list, errs...)
		if len(errs) > 0 && isFailImmediately(from.Config) {
			return list
		}
	}
	return list
}

func replaceHome(bytes []byte) string {
	var home, err = os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: could not get home directory")
	}
	var absPath, _ = filepath.Abs(home)
	return strings.Replace(string(bytes), "${HOME}", absPath, -1)
}
