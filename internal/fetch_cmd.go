package internal

import (
	"fmt"
	"os/exec"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/lib"
)

/*
Progress represents a fetching progress.
*/
type Progress struct {
	current int
	total   int
}

/*
NewProgress genereate an instance of Progress.
If the given value is negative value, this method treats the given value is 0.
*/
func NewProgress(total int) *Progress {
	if total < 0 {
		total = 0
	}
	return &Progress{total: total}
}

/*
String returns the string representation of the progress.
*/
func (progress *Progress) String() string {
	return fmt.Sprintf("%3d/%3d", progress.current, progress.total)
}

/*
Increment adds 1 to the progress.
However, if current value is equals to total, this method does nothing.
*/
func (progress *Progress) Increment() {
	if progress.current != progress.total {
		progress.current++
	}
}

/*
FetchCommand represents a command.
*/
type FetchCommand struct {
	options *options
}

/*
FetchCommandFactory returns an instance of command.
*/
func FetchCommandFactory() (cli.Command, error) {
	return &FetchCommand{&options{}}, nil
}

/*
Help returns the help message of the command.
*/
func (fetch *FetchCommand) Help() string {
	return `rrh fetch [OPTIONS] [GROUPS...]
OPTIONS
    -r, --remote <REMOTE>   specify the remote name. Default is "origin."
ARGUMENTS
    GROUPS                  run "git fetch" command on each repository on the group.
                            if no value is specified, run on the default group.`
}

/*
Synopsis returns the help message of the command.
*/
func (fetch *FetchCommand) Synopsis() string {
	return "run \"git fetch\" on the given groups."
}

/*
Run performs the command.
*/
func (fetch *FetchCommand) Run(args []string) int {
	var options, err = fetch.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var config = lib.OpenConfig()
	if len(options.args) == 0 {
		options.args = []string{config.GetValue(lib.RrhDefaultGroupName)}
	}
	var db, err2 = lib.Open(config)
	if err2 != nil {
		fmt.Println(err2.Error())
		return 1
	}
	return printErrors(config, fetch.perform(db))
}

func (fetch *FetchCommand) perform(db *lib.Database) []error {
	var errorlist = []error{}
	var onError = db.Config.GetValue(lib.RrhOnError)
	var relations = lib.FindTargets(db, fetch.options.args)
	var progress = NewProgress(len(relations))

	for _, relation := range relations {
		var err = fetch.FetchRepository(db, &relation, progress)
		if err != nil {
			if onError == lib.FailImmediately {
				return []error{err}
			}
			errorlist = append(errorlist, err)
		}
	}
	return errorlist
}

type options struct {
	remote string
	// key      string
	// userName string
	// password string
	args []string
}

func (fetch *FetchCommand) parse(args []string) (*options, error) {
	var options = options{"origin", []string{}}
	flags := flag.NewFlagSet("fetch", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(fetch.Help()) }
	flags.StringVarP(&options.remote, "remote", "r", "origin", "remote name")
	// flags.StringVar(&options.key, "k", "", "private key path")
	// flags.StringVar(&options.userName, "u", "", "user name")
	// flags.StringVar(&options.password, "p", "", "password")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	options.args = flags.Args()
	fetch.options = &options
	return &options, nil
}

/*
DoFetch executes fetch operation of git.
Currently, fetch is conducted by the system call.
Ideally, fetch is performed by using go-git.
*/
func (fetch *FetchCommand) DoFetch(repo *lib.Repository, relation *lib.Relation, progress *Progress) error {
	var cmd = exec.Command("git", "fetch", fetch.options.remote)
	cmd.Dir = repo.Path
	progress.Increment()
	fmt.Printf("%s fetching %s....", progress, relation)
	var output, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("%s,%s", relation, err.Error())
	}
	fmt.Printf("done\n%s", output)
	return nil
}

/*
FetchRepository execute `git fetch` on the given repository.
*/
func (fetch *FetchCommand) FetchRepository(db *lib.Database, relation *lib.Relation, progress *Progress) error {
	var repository = db.FindRepository(relation.RepositoryID)
	if repository == nil {
		return fmt.Errorf("%s: repository not found", relation)
	}
	return fetch.DoFetch(repository, relation, progress)
}
