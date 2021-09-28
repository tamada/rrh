package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh/cmd/rrh/commands/add"
	"github.com/tamada/rrh/cmd/rrh/commands/group"
	"github.com/tamada/rrh/cmd/rrh/commands/list"
	"github.com/tamada/rrh/cmd/rrh/commands/prune"
)

var (
	verboseMode bool
	configFile  string
)

func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Short: "Remote Repositories Head/Repositories, Ready to Hack",
		Use:   "rrh",
		Args:  cobra.ArbitraryArgs,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("subcommand not found")
			} else if done, err := findAndExecuteAlias(c, args); done {
				return err
			} else if done, err := findAndExecuteExternalCommand(c, args); done {
				return err
			}
			return fmt.Errorf("%s: not found internal commands, external commands and aliases", args[0])
		},
	}
	rootCmd.SetOut(os.Stdout)

	flags := rootCmd.PersistentFlags()
	flags.BoolVarP(&verboseMode, "verbose", "v", false, "verbose mod")
	flags.StringVarP(&configFile, "config-file", "c", "${HOME}/.config/rrh/config.json",
		"specifies the config file path.")

	registerSubCommands(rootCmd)

	return rootCmd
}

func registerSubCommands(c *cobra.Command) {
	c.AddCommand(AliasCommand())
	c.AddCommand(prune.New())
	c.AddCommand(group.New())
	c.AddCommand(add.New())
	c.AddCommand(list.New())
}

func loadAndFindAlias(c *cobra.Command, args []string) (*Alias, error) {
	aliases, err := loadAliases()
	if err != nil {
		return nil, err
	}
	alias := findAlias(args[0], aliases)
	if alias == nil {
		return nil, fmt.Errorf("%s: alias not found", args[0])
	}
	return alias, nil
}

func findAndExecuteAlias(c *cobra.Command, args []string) (bool, error) {
	alias, err := loadAndFindAlias(c, args)
	if err != nil {
		return false, err
	}
	return true, executeAlias(c, args, alias)
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

func Execute() error {
	return RootCommand().Execute()
}
