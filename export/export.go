package export

import (
	"encoding/json"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamadalab/grim/common"
)

type ExportCommand struct {
}

func ExportCommandFactory() (cli.Command, error) {
	return &ExportCommand{}, nil
}

func (export *ExportCommand) Help() string {
	return `grim export`
}

func (export *ExportCommand) Run(args []string) int {
	var config = common.OpenConfig()
	var db = common.Open(config)
	var bytes, _ = json.Marshal(db)
	fmt.Println(string(bytes))
	return 0
}

func (export *ExportCommand) Synopsis() string {
	return "export GRIM database to stdout."
}
