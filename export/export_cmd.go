package export

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/common"
)

type ExportCommand struct {
}

type exportOptions struct {
	NoIndent bool
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
    --no-indent    print result as no indented json (Default indented json)`
}

/*
Run peforms the command.
*/
func (export *ExportCommand) Run(args []string) int {
	options, err := export.parse(args)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	var config = common.OpenConfig()
	db, err := common.Open(config)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}

	var result, _ = json.Marshal(db)
	if options.NoIndent {
		fmt.Println(string(result))
	} else {
		var buffer bytes.Buffer
		err := json.Indent(&buffer, result, "", "  ")
		if err != nil {
			fmt.Println(err.Error())
			return 1
		}
		fmt.Println(buffer.String())
	}
	return 0
}

func (export *ExportCommand) parse(args []string) (*exportOptions, error) {
	var options = exportOptions{false}
	flags := flag.NewFlagSet("export", flag.ExitOnError)
	flags.Usage = func() { fmt.Println(export.Help()) }
	flags.BoolVar(&options.NoIndent, "no-indent", false, "print not indented result")

	if err := flags.Parse(args); err != nil {
		return nil, err
	}
	return &options, nil
}

/*
Synopsis returns the simple help message of the command.
*/
func (export *ExportCommand) Synopsis() string {
	return "export RRH database to stdout."
}
