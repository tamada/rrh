package internal

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
	"github.com/tamada/rrh"
)

/*
ExportCommand represents a command.
*/
type ExportCommand struct {
	options *exportOptions
}

type exportOptions struct {
	noIndent   bool
	noHideHome bool
}

/*
ExportCommandFactory generate the command struct.
*/
func ExportCommandFactory() (cli.Command, error) {
	return &ExportCommand{}, nil
}

/*
Help returns the help message of the command.
*/
func (export *ExportCommand) Help() string {
	return `rrh export [OPTIONS]
OPTIONS
    --no-indent      print result as no indented json
    --no-hide-home   not replace home directory to '${HOME}' keyword`
}

/*
Run peforms the command.
*/
func (export *ExportCommand) Run(args []string) int {
	var _, err = export.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var config = rrh.OpenConfig()
	db, err := rrh.Open(config)
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

func hideHome(result string) string {
	var home, err = homedir.Dir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: chould not get home directory")
	}
	var absPath, _ = filepath.Abs(home)
	return strings.Replace(result, absPath, "${HOME}", -1)
}

func (export *ExportCommand) perform(db *rrh.Database) int {
	var result, _ = json.Marshal(db)
	var stringResult = string(result)
	if !export.options.noHideHome {
		stringResult = hideHome(stringResult)
	}

	if !export.options.noIndent {
		var result, err = indentJSON(stringResult)
		if err != nil {
			return printErrors(db.Config, []error{err})
		}
		stringResult = result
	}
	fmt.Println(stringResult)
	return 0
}

func (export *ExportCommand) parse(args []string) (*exportOptions, error) {
	var options = exportOptions{false, false}
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
func (export *ExportCommand) Synopsis() string {
	return "export rrh database to stdout."
}
