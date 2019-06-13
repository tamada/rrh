package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/lib"
)

/*
CloneCommand represents a command.
*/
type CloneCommand struct {
	options *cloneOptions
}

/*
CloneCommandFactory returns an instance of the CloneCommand.
*/
func CloneCommandFactory() (cli.Command, error) {
	return &CloneCommand{&cloneOptions{}}, nil
}

type cloneOptions struct {
	group   string
	dest    string
	verbose bool
}

/*
Help function shows the help message.
*/
func (clone *CloneCommand) Help() string {
	return `rrh clone [OPTIONS] <REMOTE_REPOS...>
OPTIONS
    -g, --group <GROUP>   print managed repositories categorized in the group.
    -d, --dest <DEST>     specify the destination.
    -v, --verbose         verbose mode.
ARGUMENTS
    REMOTE_REPOS          repository urls`
}

/*
Synopsis returns the help message of the command.
*/
func (clone *CloneCommand) Synopsis() string {
	return "run \"git clone\" and register it to a group."
}

func (clone *CloneCommand) printIfVerbose(message string) {
	if clone.options.verbose {
		fmt.Println(message)
	}
}

func (options *cloneOptions) showError(list []error) {
	for _, err := range list {
		fmt.Println(err.Error())
	}
}

/*
Run performs the command.
*/
func (clone *CloneCommand) Run(args []string) int {
	var config = lib.OpenConfig()
	arguments, err := clone.parse(args, config)
	if err != nil || len(arguments) == 0 {
		fmt.Printf(clone.Help())
		return 1
	}
	db, err := lib.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 2
	}
	return clone.perform(db, arguments)
}

func (clone *CloneCommand) perform(db *lib.Database, arguments []string) int {
	var count, list = clone.DoClone(db, arguments)
	if len(list) != 0 {
		clone.options.showError(list)
		var onError = db.Config.GetValue(lib.RrhOnError)
		if onError == lib.Fail || onError == lib.FailImmediately {
			return 1
		}
	}
	db.StoreAndClose()
	printResult(count, clone.options.dest, clone.options.group)
	return 0
}

func printResult(count int, dest string, group string) {
	switch count {
	case 0:
		fmt.Println("no repositories cloned")
	case 1:
		fmt.Printf("a repository cloned into %s and registered to group %s\n", dest, group)
	default:
		fmt.Printf("%d repositories cloned into %s and registered to group %s\n", count, dest, group)
	}
}

func (clone *CloneCommand) buildFlagSets(config *lib.Config) (*flag.FlagSet, *cloneOptions) {
	var defaultGroup = config.GetValue(lib.RrhDefaultGroupName)
	var destination = config.GetValue(lib.RrhCloneDestination)
	var options = cloneOptions{defaultGroup, ".", false}
	flags := flag.NewFlagSet("clone", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(clone.Help()) }
	flags.StringVarP(&options.group, "group", "g", defaultGroup, "belonging group")
	flags.StringVarP(&options.dest, "dest", "d", destination, "destination")
	flags.BoolVarP(&options.verbose, "verbose", "v", false, "verbose mode")
	return flags, &options
}

func (clone *CloneCommand) parse(args []string, config *lib.Config) ([]string, error) {
	var flags, options = clone.buildFlagSets(config)
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	clone.options = options

	return flags.Args(), nil
}

func registerPath(db *lib.Database, dest string, repoID string) (*lib.Repository, error) {
	var path, err = filepath.Abs(dest)
	if err != nil {
		return nil, err
	}
	var remotes, err2 = lib.FindRemotes(path)
	if err2 != nil {
		return nil, err2
	}
	fmt.Printf("createRepository(%s, %s)\n", repoID, path)
	var repo, err3 = db.CreateRepository(repoID, path, "", remotes)
	if err3 != nil {
		return nil, err3
	}
	return repo, nil
}

func (clone *CloneCommand) toDir(db *lib.Database, URL string, dest string, repoID string) (*lib.Repository, error) {
	clone.printIfVerbose(fmt.Sprintf("git clone %s %s (%s)", URL, dest, repoID))
	var cmd = exec.Command("git", "clone", URL, dest)
	var err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("%s: clone error (%s)", URL, err.Error())
	}
	return registerPath(db, dest, repoID)
}

func isExistDir(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	stat, err := os.Stat(abs)
	return !os.IsNotExist(err) && stat.IsDir()
}

/*
DoClone performs `git clone` command and register the cloned repositories to RRH database.
*/
func (clone *CloneCommand) DoClone(db *lib.Database, arguments []string) (int, []error) {
	if len(arguments) == 1 {
		var err = clone.doCloneARepository(db, arguments[0])
		if err != nil {
			return 0, []error{err}
		}
		return 1, []error{}
	}
	return clone.doCloneRepositories(db, arguments)
}

func (clone CloneCommand) doCloneRepositories(db *lib.Database, arguments []string) (int, []error) {
	var errorlist = []error{}
	var count = 0
	for _, url := range arguments {
		var increment, err = clone.doCloneEachRepository(db, url)
		if err != nil {
			errorlist = append(errorlist, err)
			if db.Config.GetValue(lib.RrhOnError) == lib.FailImmediately {
				return count, errorlist
			}
		}
		count += increment
	}
	return count, errorlist
}

func (clone *CloneCommand) relateTo(db *lib.Database, groupID string, repoID string) error {
	var _, err = db.AutoCreateGroup(groupID, "", false)
	if err != nil {
		return fmt.Errorf("%s: group not found", groupID)
	}
	db.Relate(groupID, repoID)
	return nil
}

/*
doCloneEachRepository performes `git clone` for each repository.
This function is called repeatedly.
*/
func (clone *CloneCommand) doCloneEachRepository(db *lib.Database, URL string) (int, error) {
	var count int
	var id = findIDFromURL(URL)
	var path = filepath.Join(clone.options.dest, id)
	var _, err = clone.toDir(db, URL, path, id)
	if err == nil {
		if err := clone.relateTo(db, clone.options.group, id); err != nil {
			return count, err
		}
		count++
	}
	return count, err
}

func (clone *CloneCommand) doCloneARepository(db *lib.Database, URL string) error {
	var id, path string

	if isExistDir(clone.options.dest) {
		id = findIDFromURL(URL)
		path = filepath.Join(clone.options.dest, id)
	} else {
		var _, newid = filepath.Split(clone.options.dest)
		path = clone.options.dest
		id = newid
	}
	var _, err = clone.toDir(db, URL, path, id)
	if err != nil {
		return err
	}
	return clone.relateTo(db, clone.options.group, id)
}

func findIDFromURL(URL string) string {
	var _, dir = path.Split(URL)
	if strings.HasSuffix(dir, ".git") {
		return strings.TrimSuffix(dir, ".git")
	}
	return dir
}
