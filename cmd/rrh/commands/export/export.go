package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/utils"
)

type exportOptions struct {
	noIndent   bool
	noHideHome bool
}

var exportOpts = &exportOptions{}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "export rrh database to stdout",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return utils.PerformRrhCommand(c, args, perform)
		},
	}
	flags := cmd.Flags()
	flags.BoolVarP(&exportOpts.noIndent, "no-indent", "", false, "print result as no indented json")
	flags.BoolVarP(&exportOpts.noHideHome, "no-hide-home", "", false, "not replace home directory to '${HOME}' keyword")
	return cmd
}

func perform(c *cobra.Command, args []string, db *rrh.Database) error {
	var result, _ = json.Marshal(db)
	var stringResult = string(result)
	if !exportOpts.noHideHome {
		stringResult = hideHome(stringResult)
	}

	if !exportOpts.noIndent {
		var result, err = indentJSON(stringResult)
		if err != nil {
			return err
		}
		stringResult = result
	}
	c.Println(stringResult)
	return nil
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
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: chould not get home directory")
	}
	var absPath, _ = filepath.Abs(home)
	return strings.Replace(result, absPath, "${HOME}", -1)
}
