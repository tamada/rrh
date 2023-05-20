package sfg

import (
	"embed"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
)

type sfgOptions struct {
	shellName      string
	withoutCdrrh   bool
	withoutRrhPeco bool
	withoutRrhFzf  bool
}

var sfgOpts = &sfgOptions{}

//go:embed data
var functions embed.FS

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "generate shell functions for shell",
		Args:  validateArgs,
		RunE:  perform,
	}
	flags := cmd.Flags()
	flags.StringVarP(&sfgOpts.shellName, "shell", "s", "", "specify the shell name. this option is mandatory. availables are: bash, and zsh")
	flags.BoolVarP(&sfgOpts.withoutCdrrh, "without-cdrrh", "", false, "not generate the cdrrh function")
	flags.BoolVarP(&sfgOpts.withoutRrhPeco, "without-rrhpeco", "", false, "not generate the rrhpeco function")
	flags.BoolVarP(&sfgOpts.withoutRrhFzf, "without-rrhfzf", "", false, "not generate the rrhfzf function")
	return cmd
}

func validateArgs(c *cobra.Command, args []string) error {
	if sfgOpts.shellName == "" {
		return errors.New("shell name option is mandatory")
	}
	switch strings.ToLower(sfgOpts.shellName) {
	case "sh", "bash", "zsh":
		return nil
	default:
		return fmt.Errorf("%s: unsupported shell", sfgOpts.shellName)
	}
}

func perform(c *cobra.Command, args []string) error {
	c.SilenceUsage = true
	shell := strings.ToLower(sfgOpts.shellName)
	switch shell {
	case "sh", "bash", "zsh":
		return printFunctions(c, shell, "bash")
	}
	return fmt.Errorf("%s: unsupported shell", sfgOpts.shellName)
}

func printFunctions(c *cobra.Command, shellName, defaultShellName string) error {
	executors := []struct {
		withoutFlag bool
		scriptName  string
	}{
		{sfgOpts.withoutCdrrh, "cdrrh"},
		{sfgOpts.withoutRrhPeco, "rrhpeco"},
		{sfgOpts.withoutRrhFzf, "rrhfzf"},
	}
	for _, executor := range executors {
		if !executor.withoutFlag {
			if err := printFunction(c, shellName, executor.scriptName); err == nil {
				continue
			}
			if err := printFunction(c, defaultShellName, executor.scriptName); err != nil {
				return err
			}
		}
	}
	return nil
}

func printFunction(c *cobra.Command, shellName, scriptName string) error {
	in, err := functions.Open(fmt.Sprintf("data/%s/%s", shellName, scriptName))
	if err != nil {
		return err
	}
	defer in.Close()
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return err
	}
	c.Println(string(data))
	return nil
}
