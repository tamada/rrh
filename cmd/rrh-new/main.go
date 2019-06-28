package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/lib"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

type newOptions struct {
	description string
	group       string
	homepage    string
	parentPath  string
	privateFlag bool
	helpFlag    bool
}

func getHelpMessage() string {
	return `rrh new [OPTIONS] <[ORGANIZATION/]REPOSITORIES...>
OPTIONS
    -g, --group <GROUP>         specifies group name.
    -H, --homepage <URL>        specifies homepage url.
    -d, --description <DESC>    specifies short description of the repository.
    -p, --private               create a private repository.
    -p, --parent-path <PATH>    specifies the destination path (default: '.').
    -h, --help                  print this message.
ARGUMENTS
    ORGANIZATION    specifies organization, if needed.
    REPOSITORY      specifies repository name, and it is directory name.
`
}

func buildFlagSet(config *lib.Config) (*flag.FlagSet, *newOptions) {
	var opt = newOptions{}
	var defaultGroup = config.GetValue(lib.RrhDefaultGroupName)
	flags := flag.NewFlagSet("new", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(getHelpMessage()) }
	flags.StringVarP(&opt.description, "description", "d", "", "specifys description of the project")
	flags.StringVarP(&opt.group, "group", "g", defaultGroup, "target group")
	flags.StringVarP(&opt.homepage, "homepage", "H", "", "specifies homepage url")
	flags.StringVarP(&opt.parentPath, "parent-path", "P", ".", "specifies the destination path")
	flags.BoolVarP(&opt.privateFlag, "private", "p", false, "create a private repository")
	flags.BoolVarP(&opt.helpFlag, "help", "h", false, "print this message")
	return flags, &opt
}

func createArgsToHubCommand(projectName string, opts *newOptions) []string {
	var args = []string{"create"}
	if opts.homepage != "" {
		args = append(args, "--homepage")
		args = append(args, opts.homepage)
	}
	if opts.privateFlag {
		args = append(args, "--private")
	}
	if opts.description != "" {
		args = append(args, "--description")
		args = append(args, opts.description)
	}
	args = append(args, projectName)
	return args
}

func createProjectPage(repo *repo, opts *newOptions) error {
	var argsToHub = createArgsToHubCommand(repo.givenString, opts)
	var cmd = exec.Command("hub", argsToHub...)
	cmd.Dir = repo.dest
	var _, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("%s: %s", repo.givenString, err.Error())
	}
	return nil
}

func createReadme(dest, projectName string) {
	var path = filepath.Join(dest, "README.md")
	var file, err = os.OpenFile(path, os.O_CREATE, 644)
	defer file.Close()
	if err == nil {
		file.WriteString(fmt.Sprintf("# %s", projectName))
	}
}

func makeGitDirectory(config *lib.Config, repo *repo, opts *newOptions) error {
	os.MkdirAll(repo.dest, 0755)
	var gitDir = filepath.Join(repo.dest, ".git")
	var fsys = osfs.New(gitDir)
	var _, err = git.Init(filesystem.NewStorage(fsys, cache.NewObjectLRUDefault()), fsys)
	if err != nil {
		return err
	}
	createReadme(repo.dest, repo.repoName)
	return nil
}

type repo struct {
	givenString string
	dest        string
	remoteURL   string
	githubRepo  string
	repoName    string
}

func availableDir(opts *newOptions, dir string) bool {
	var path = filepath.Join(opts.parentPath, dir)
	var _, err = os.Stat(path)
	return err == nil
}

func findDirectoryName(arg string, opts *newOptions) (string, error) {
	var terms = strings.Split(arg, "/")
	var dest = arg
	if len(terms) == 2 {
		dest = terms[1]
	} else if len(terms) > 2 {
		return "", fmt.Errorf("%s: illegal format for specifying project", arg)
	}
	if availableDir(opts, dest) {
		return "", fmt.Errorf("%s/%s: directory exist", opts.parentPath, dest)
	}
	return dest, nil
}

func findRepoName(arg string) string {
	var terms = strings.Split(arg, "/")
	if len(terms) == 1 {
		return arg
	}
	return terms[1]
}

func createRepo(config *lib.Config, arg string, opts *newOptions) (*repo, error) {
	var dest, err = findDirectoryName(arg, opts)
	if err != nil {
		return nil, err
	}
	var repoName = findRepoName(arg)
	return &repo{givenString: arg, dest: filepath.Join(opts.parentPath, dest), repoName: repoName}, nil
}

func registerToGroup(db *lib.Database, repo *repo, opts *newOptions) error {
	var remotes, _ = lib.FindRemotes(repo.dest)
	var _, err1 = db.CreateRepository(repo.repoName, repo.dest, opts.description, remotes)
	if err1 != nil {
		return err1
	}
	var err2 = db.Relate(opts.group, repo.repoName)
	if err2 != nil {
		return err2
	}
	return nil
}

func createRepository(db *lib.Database, arg string, opts *newOptions) error {
	var repo, err = createRepo(db.Config, arg, opts)
	if err == nil {
		err = makeGitDirectory(db.Config, repo, opts)
		fmt.Printf("1/3 create git directory on %s\n", repo.dest)
	}
	if err == nil {
		err = createProjectPage(repo, opts)
		fmt.Printf("2/3 create remote repository of %s\n", repo.repoName)
	}
	if err == nil {
		err = registerToGroup(db, repo, opts)
		fmt.Printf("3/3 add %s to rrh database with group %s\n", repo.repoName, opts.group)
	}
	return err
}

func isFailImmediately(config *lib.Config) bool {
	var onError = config.GetValue(lib.RrhOnError)
	return onError == lib.FailImmediately
}

func createRepositories(config *lib.Config, args []string, opts *newOptions) []error {
	var errors = []error{}
	var db, err = lib.Open(config)
	if err != nil {
		return []error{err}
	}

	for _, arg := range args[1:] {
		var err = createRepository(db, arg, opts)
		if err != nil {
			if isFailImmediately(config) {
				return []error{err}
			}
			errors = append(errors, err)
		}
	}
	if len(errors) == 0 || config.GetValue(lib.RrhOnError) == lib.Ignore {
		db.StoreAndClose()
	}
	return errors
}

func perform(config *lib.Config, args []string, opts *newOptions) int {
	if opts.helpFlag {
		fmt.Println(getHelpMessage())
		return 0
	}
	var errs = createRepositories(config, args, opts)
	return config.PrintErrors(errs)
}

func goMain(args []string) int {
	var config = lib.OpenConfig()
	var flag, opts = buildFlagSet(config)
	if err := flag.Parse(args); err != nil {
		fmt.Println(err.Error())
		return 1
	}
	return perform(config, flag.Args(), opts)
}

func main() {
	var status = goMain(os.Args)
	os.Exit(status)
}
