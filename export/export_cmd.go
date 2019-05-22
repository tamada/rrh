package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/mitchellh/go-homedir"
	flag "github.com/spf13/pflag"
	"github.com/tamada/rrh/common"
)

/*
Command represents a command.
*/
type Command struct {
	options *options
}

type options struct {
	noIndent   bool
	noHideHome bool
}

/*
CommandFactory generate the command struct.
*/
func CommandFactory() (cli.Command, error) {
	return &Command{}, nil
}

/*
Help returns the help message of the command.
*/
func (export *Command) Help() string {
	return `rrh export [OPTIONS]
OPTIONS
    --no-indent      print result as no indented json
    --no-hide-home   not replace home directory to '${HOME}' keyword`
}

/*
Run peforms the command.
*/
func (export *Command) Run(args []string) int {
	var _, err = export.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var config = common.OpenConfig()
	db, err := common.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 2
	}
	return export.perform(db)
}

func indentJSON(result string) (string, error) {
	var buffer bytes.Buffer
	var err = json.Indent(&buffer, []byte(result), "", "  ")
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func printError(err error) int {
	fmt.Println(err.Error())
	return 5
}

func hideHome(result string) string {
	var home, err = homedir.Dir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: chould not get home directory")
	}
	var absPath, _ = filepath.Abs(home)
	return strings.Replace(result, absPath, "${HOME}", -1)
}

func (export *Command) perform(db *common.Database) int {
	var result, _ = json.Marshal(db)
	var stringResult = string(result)
	if !export.options.noHideHome {
		stringResult = hideHome(stringResult)
	}

	if !export.options.noIndent {
		var result, err = indentJSON(stringResult)
		if err != nil {
			return printError(err)
		}
		stringResult = result
	}
	fmt.Println(stringResult)
	return 0
}

func (export *Command) parse(args []string) (*options, error) {
	var options = options{false, false}
	flags := flag.NewFlagSet("export", flag.ContinueOnError)
	flags.Usage = func() { fmt.Println(export.Help()) }
	flags.BoolVar(&options.noIndent, "no-indent", false, "print not indented result")
	flags.BoolVar(&options.noHideHome, "no-hide-home", false, "not hide home directory")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	export.options = &options
	return &options, nil
}

/*
Synopsis returns the simple help message of the command.
*/
func (export *Command) Synopsis() string {
	return "export RRH database to stdout."
}
