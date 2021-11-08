package sfg

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh/common"
)

type options struct {
	shell          string
	withoutCdRrh   bool
	withoutRrhPeco bool
	withoutRrhFzf  bool
}

var opts = &options{}

//go:embed scripts/**
var scripts embed.FS

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "generating shell functions",
		RunE:  perform,
	}
	flags := cmd.Flags()
	flags.StringVarP(&opts.shell, "shell", "s", "bash", "specify the target shell (default: bash)")
	flags.BoolVarP(&opts.withoutCdRrh, "without-cdrrh", "", false, "generate no cdrrh function")
	flags.BoolVarP(&opts.withoutRrhPeco, "without-rrhpeco", "", false, "generate no rrhpeco function")
	flags.BoolVarP(&opts.withoutRrhFzf, "without-rrhfzf", "", false, "generate no rrhfzf function")
	return cmd
}

func printScript(c *cobra.Command, path string) error {
	data, err := scripts.ReadFile(path)
	if err == nil {
		c.Println(string(data))
	}
	return err
}

func printScripts(c *cobra.Command, shellName string) error {
	errs := common.NewErrorList()
	targets := []struct {
		withoutFlag bool
		path        string
	}{
		{opts.withoutCdRrh, "cdrrh"},
		{opts.withoutRrhPeco, "rrhpeco"},
		{opts.withoutRrhFzf, "rrhfzf"},
	}
	for _, target := range targets {
		if !target.withoutFlag {
			err := printScript(c, filepath.Join("scripts", shellName, target.path))
			errs = errs.Append(err)
		}
	}
	return errs.NilOrThis()
}

func perform(c *cobra.Command, args []string) error {
	lowerShellName := strings.ToLower(opts.shell)
	switch lowerShellName {
	case "bash", "zsh":
		return printScripts(c, "bash")
	default:
		return fmt.Errorf("%s: sorry, unsupported shell", opts.shell)
	}
}
