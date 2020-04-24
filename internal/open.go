package internal

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/skratchdot/open-golang/open"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh"
)

/*
OpenCommand represents a command.
*/
type OpenCommand struct {
	options *openOptions
}

/*
OpenCommandFactory generates the object of AddCommand.
*/
func OpenCommandFactory() (cli.Command, error) {
	return &OpenCommand{options: new(openOptions)}, nil
}

/*
Synopsis returns the simple help message of the command.
*/
func (open *OpenCommand) Synopsis() string {
	return "open folder or web page of the given repositories."
}

/*
Help function shows the help message.
*/
func (open *OpenCommand) Help() string {
	return `rrh open [OPTIONS] <REPOSITORIES...>
OPTIONS
    -f, --folder     open the folder of the specified repository (Default).
    -w, --webpage    open the webpage of the specified repository.
    -h, --help       print this message.
ARGUMENTS
    REPOSITORIES     specifies repository names.`
}

type openOptions struct {
	helpFlag    bool
	folderFlag  bool
	browserFlag bool
	args        []string
}

func (open *OpenCommand) buildFlagSet() (*flag.FlagSet, *openOptions) {
	opts := new(openOptions)
	flags := flag.NewFlagSet("open", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(open.Help()) }
	flags.BoolVarP(&opts.helpFlag, "help", "h", false, "print this message")
	flags.BoolVarP(&opts.browserFlag, "browser", "b", false, "open the webpage of the repository")
	flags.BoolVarP(&opts.folderFlag, "folder", "f", true, "open the folder of the repository")
	return flags, opts
}

func validateArgs(flags *flag.FlagSet, opts *openOptions) (*openOptions, error) {
	if !opts.helpFlag && len(flags.Args()) == 0 {
		return nil, fmt.Errorf("no arguments are specified")
	}
	if len(flags.Args()) > 0 {
		opts.args = flags.Args()
	}
	return opts, nil
}

func (open *OpenCommand) parseOptions(args []string) (*openOptions, error) {
	flags, opts := open.buildFlagSet()
	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	return validateArgs(flags, opts)
}

func (open *OpenCommand) printHelpMessage(opts *openOptions, err error) int {
	status := 0
	if err != nil {
		fmt.Println(err.Error())
	}
	if status != 0 || opts != nil && opts.helpFlag {
		fmt.Println(open.Help())
	}
	return status
}

func convertToRepositoryURL(url string) (string, error) {
	str := strings.TrimPrefix(url, "git@")
	str = strings.TrimSuffix(str, ".git")
	index := strings.Index(str, ":")
	if index < 0 {
		return "", fmt.Errorf("%s: unrecognized git repository url", url)
	}
	host := str[0:index]
	return "https://" + host + "/" + str[index+1:], nil
}

func convertURL(url string) (string, error) {
	if strings.HasPrefix(url, "git@") {
		convertedURL, err := convertToRepositoryURL(url)
		if err != nil {
			return "", err
		}
		url = convertedURL
	}
	if strings.HasPrefix(url, "https") && strings.HasSuffix(url, ".git") {
		url = strings.TrimSuffix(url, ".git")
	}
	return url, nil
}

func generateWebPageURL(repo *rrh.Repository) (string, error) {
	if len(repo.Remotes) == 0 {
		return "", fmt.Errorf("%s: remote repository not found", repo.ID)
	}
	return convertURL(repo.Remotes[0].URL)
}

func execOpen(repo *rrh.Repository, opts *openOptions) (string, error) {
	if opts.browserFlag {
		return generateWebPageURL(repo)
	}
	return repo.Path, nil
}

func performEach(arg string, opts *openOptions, db *rrh.Database) error {
	repo := db.FindRepository(arg)
	if repo == nil {
		return fmt.Errorf("%s: repository not found", arg)
	}
	path, err := execOpen(repo, opts)
	if err != nil {
		return err
	}
	return open.Start(path)
}

func performOpen(args []string, opts *openOptions) int {
	config := rrh.OpenConfig()
	db, err := rrh.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	for _, arg := range args {
		err = performEach(arg, opts, db)
		if value := config.PrintErrors(err); value != 0 {
			return value
		}
	}
	return 0
}

/*
Run performs the command.
*/
func (open *OpenCommand) Run(args []string) int {
	opts, err := open.parseOptions(args)
	if err != nil || opts.helpFlag {
		return open.printHelpMessage(opts, err)
	}
	return performOpen(opts.args, opts)
}
