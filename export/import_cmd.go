package export

import (
	"fmt"

	"github.com/mitchellh/cli"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/common"
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

func eraseDatabase(db *common.Database, command *ImportCommand) {
	db.Groups = []common.Group{}
	db.Repositories = []common.Repository{}
	db.Relations = []common.Relation{}
	command.options.printIfNeeded("The local database is cleared")
}

func perform(db *common.Database, command *ImportCommand) int {
	if command.options.overwrite {
		eraseDatabase(db, command)
	}
	var db2, err = readNewDB(command.options.database, db.Config)
	if err != nil {
		fmt.Printf(err.Error())
		return 4
	}
	var errs = command.copyDB(db2, db)
	var statusCode = db.Config.PrintErrors(errs)
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
	var config = common.OpenConfig()
	var db, err2 = common.Open(config)
	if err2 != nil {
		return 2
	}
	return perform(db, command)
}

func buildFlagSet(command *ImportCommand) (*flag.FlagSet, *importOptions) {
	var options = importOptions{false, false, false, ""}
	var flags = flag.NewFlagSet("import", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(command.Help()) }
	flags.BoolVar(&options.overwrite, "overwrite", false, "overwrite mode")
	flags.BoolVar(&options.autoClone, "auto-clone", false, "auto clone mode")
	flags.BoolVarP(&options.verbose, "verbose", "v", false, "verbose mode")
	return flags, &options
}

func parse(args []string, command *ImportCommand) error {
	var flags, options = buildFlagSet(command)
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
