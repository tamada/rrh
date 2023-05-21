package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/common"
	"gopkg.in/src-d/go-billy.v4/osfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/cache"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

type repositoryInfo struct {
	description string
	homepage    string
	privateFlag bool
}

type newOptions struct {
	group      string
	parentPath string
	info       *repositoryInfo
	dryrunMode bool
	helpFlag   bool
}

func getHelpMessage() string {
	return `rrh new [OPTIONS] <[ORGANIZATION/]REPOSITORIES...>
OPTIONS
    -d, --description <DESC>    specifies short description of the repository.
    -D, --dry-run               performs on dry-run mode.
    -g, --groups <GROUPS>       specifies group name.
    -H, --homepage <URL>        specifies homepage url.
    -p, --private               create a private repository.
    -P, --parent-path <PATH>    specifies the destination path (default: '.').
    -h, --help                  print this message.
ARGUMENTS
    ORGANIZATION    specifies organization, if needed.
    REPOSITORY      specifies repository name, and it is directory name.`
}

func buildFlagSet(config *rrh.Config) (*flag.FlagSet, *newOptions) {
	var opt = newOptions{info: new(repositoryInfo)}
	var defaultGroup = config.GetValue(rrh.DefaultGroupName)
	flags := flag.NewFlagSet("new", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(getHelpMessage()) }
	flags.StringVarP(&opt.info.description, "description", "d", "", "specifys description of the project")
	flags.StringVarP(&opt.group, "group", "g", defaultGroup, "target group")
	flags.StringVarP(&opt.info.homepage, "homepage", "H", "", "specifies homepage url")
	flags.StringVarP(&opt.parentPath, "parent-path", "P", ".", "specifies the destination path")
	flags.BoolVarP(&opt.info.privateFlag, "private", "p", false, "create a private repository")
	flags.BoolVarP(&opt.dryrunMode, "dry-run", "D", false, "performs on dry-run mode")
	flags.BoolVarP(&opt.helpFlag, "help", "h", false, "print this message")
	return flags, &opt
}

func createArgsToHubCommand(projectName string, info *repositoryInfo) []string {
	var args = []string{"create"}
	if info.homepage != "" {
		args = append(args, "--homepage")
		args = append(args, info.homepage)
	}
	if info.privateFlag {
		args = append(args, "--private")
	}
	if info.description != "" {
		args = append(args, "--description")
		args = append(args, info.description)
	}
	args = append(args, projectName)
	return args
}

func createProjectPage(repo *repo, opts *newOptions) (string, error) {
	var argsToHub = createArgsToHubCommand(repo.givenString, opts.info)
	var cmdString = "hub " + strings.Join(argsToHub, " ")
	if opts.dryrunMode {
		return cmdString, nil
	}
	var cmd = exec.Command("hub", argsToHub...)
	cmd.Dir = repo.dest
	var _, err = cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s: %s", repo.givenString, err.Error())
	}
	return cmdString, nil
}

func createReadme(dest, projectName string) {
	var path = filepath.Join(dest, "README.md")
	var file, err = os.OpenFile(path, os.O_CREATE, 0644)
	defer file.Close()
	if err == nil {
		file.WriteString(fmt.Sprintf("# %s", projectName))
	}
}

func makeGitDirectory(config *rrh.Config, repo *repo, opts *newOptions) error {
	if opts.dryrunMode {
		return nil
	}
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
	return convertToAbsolutePath(dest, opts)
}

func convertToAbsolutePath(dest string, opts *newOptions) (string, error) {
	var abs, err = filepath.Abs(filepath.Join(opts.parentPath, dest))
	if err != nil {
		return "", err
	}
	return abs, nil
}

func findRepoName(arg string) string {
	var terms = strings.Split(arg, "/")
	if len(terms) == 1 {
		return arg
	}
	return terms[1]
}

func createRepo(config *rrh.Config, arg string, opts *newOptions) (*repo, error) {
	var dest, err = findDirectoryName(arg, opts)
	if err != nil {
		return nil, err
	}
	var repoName = findRepoName(arg)
	return &repo{givenString: arg, dest: dest, repoName: repoName}, nil
}

func registerToGroup(db *rrh.Database, repo *repo, opts *newOptions) error {
	if opts.dryrunMode {
		return nil
	}
	var remotes, _ = rrh.FindRemotes(repo.dest)
	var _, err1 = db.CreateRepository(repo.repoName, repo.dest, opts.info.description, remotes)
	if err1 != nil {
		return err1
	}
	var err2 = db.Relate(opts.group, repo.repoName)
	if err2 != nil {
		return err2
	}
	return nil
}

func createRepository(db *rrh.Database, arg string, opts *newOptions) error {
	var repo, err = createRepo(db.Config, arg, opts)
	if err == nil {
		err = makeGitDirectory(db.Config, repo, opts)
		fmt.Printf("1/3 create git directory on \"%s\"\n", repo.dest)
	}
	if err == nil {
		var cmd string
		cmd, err = createProjectPage(repo, opts)
		fmt.Printf("2/3 create remote repository of %s by \"%s\"\n", repo.repoName, cmd)
	}
	if err == nil {
		err = registerToGroup(db, repo, opts)
		fmt.Printf("3/3 add repository \"%s\" to group \"%s\"\n", repo.repoName, opts.group)
	}
	return err
}

func storeDbWhenSucceeded(db *rrh.Database, errors common.ErrorList) {
	if errors.IsNil() {
		db.StoreAndClose()
	}
}

func createRepositories(config *rrh.Config, args []string, opts *newOptions) error {
	var errors = common.NewErrorList()
	var db, err = rrh.Open(config)
	defer storeDbWhenSucceeded(db, errors)
	errors = errors.Append(err)
	if errors.IsErr() {
		return errors
	}
	for _, arg := range args[1:] {
		if err := createRepository(db, arg, opts); err != nil {
			errors = errors.Append(err)
		}
	}
	return errors
}

func perform(config *rrh.Config, args []string, opts *newOptions) int {
	if opts.helpFlag {
		fmt.Println(getHelpMessage())
		return 0
	}
	var errs = createRepositories(config, args, opts)
	fmt.Println(errs.Error())
	if errs != nil {
		return 1
	}
	return 0
}

func goMain(args []string) int {
	var config = rrh.OpenConfig()
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
