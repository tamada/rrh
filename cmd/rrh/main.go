package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/cmd/rrh/commands/add"
	"github.com/tamada/rrh/cmd/rrh/commands/alias"
	"github.com/tamada/rrh/cmd/rrh/commands/clone"
	"github.com/tamada/rrh/cmd/rrh/commands/config"
	"github.com/tamada/rrh/cmd/rrh/commands/execcmd"
	"github.com/tamada/rrh/cmd/rrh/commands/group"
	"github.com/tamada/rrh/cmd/rrh/commands/list"
	"github.com/tamada/rrh/cmd/rrh/commands/migrate"
	"github.com/tamada/rrh/cmd/rrh/commands/open"
	"github.com/tamada/rrh/cmd/rrh/commands/prune"
	"github.com/tamada/rrh/cmd/rrh/commands/repository"
	"github.com/tamada/rrh/cmd/rrh/commands/sfg"
)

var (
	verboseMode bool
	configFile  string
)

func rootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Short:   "Remote Repositories Head/Repositories, Ready to Hack",
		Use:     "rrh",
		Version: rrh.VERSION,
		Args:    cobra.ArbitraryArgs,
		RunE: func(c *cobra.Command, args []string) error {
			config := rrh.OpenConfig()
			if len(args) == 0 {
				return fmt.Errorf("subcommand not found")
			} else if done, err := findAndExecuteAlias(c, args, config); done {
				return err
			} else if done, err := findAndExecuteExternalCommand(c, args); done {
				return err
			}
			return fmt.Errorf("%s: not found internal commands, external commands and aliases", args[0])
		},
	}
	rootCmd.SetOut(os.Stdout)

	flags := rootCmd.PersistentFlags()
	flags.BoolVarP(&verboseMode, "verbose", "v", false, "verbose mode")
	flags.StringVarP(&configFile, "config-file", "c", "${HOME}/.config/rrh/config.json",
		"specifies the config file path.")

	registerSubCommands(rootCmd)

	return rootCmd
}

func registerSubCommands(c *cobra.Command) {
	c.AddCommand(alias.New())
	c.AddCommand(add.New())
	c.AddCommand(clone.New())
	c.AddCommand(config.New())
	c.AddCommand(execcmd.New())
	c.AddCommand(group.New())
	c.AddCommand(list.New())
	c.AddCommand(open.New())
	c.AddCommand(prune.New())
	c.AddCommand(repository.New())
	c.AddCommand(migrate.New())
	c.AddCommand(sfg.New())
}

func loadAndFindAlias(c *cobra.Command, args []string, config *rrh.Config) (*alias.Command, error) {
	aliases, err := alias.LoadAliases(config)
	if err != nil {
		return nil, err
	}
	alias := alias.FindAlias(args[0], aliases)
	if alias == nil {
		return nil, fmt.Errorf("%s: alias not found", args[0])
	}
	return alias, nil
}

func findAndExecuteAlias(c *cobra.Command, args []string, config *rrh.Config) (bool, error) {
	command, err := loadAndFindAlias(c, args, config)
	if err != nil {
		return false, err
	}
	return true, command.Execute(c, args)
}

func executeCommand(commandPath string, c *cobra.Command, args []string) error {
	cmd := exec.Command(commandPath, args...)
	output, err := cmd.CombinedOutput()
	c.Print(string(output))
	return err
}

func findAndExecuteExternalCommand(c *cobra.Command, args []string) (bool, error) {
	commandName := fmt.Sprintf("rrh-%s", args[0])
	command, err := findExecutableFromPathEnv(commandName)
	if err != nil {
		return false, err
	}
	return true, executeCommand(command, c, args[1:])
}

func findExecutableFromPathEnv(commandName string) (string, error) {
	var pathEnv = os.Getenv("PATH")
	for _, env := range strings.Split(pathEnv, ":") {
		if findExecutableFromDir(env, commandName) {
			return filepath.Join(env, commandName), nil
		}
	}
	return "", fmt.Errorf("%s: command not found", commandName)
}

func findExecutableFromDir(dir, commandName string) bool {
	path := filepath.Join(dir, commandName)
	finfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	if finfo.Mode().IsRegular() && (finfo.Mode().Perm()&0555) == 0555 {
		return true
	}
	return false
}

func main() {
	err := rootCommand().Execute()
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
