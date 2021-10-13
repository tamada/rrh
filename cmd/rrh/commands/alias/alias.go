package alias

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
)

func New() *cobra.Command {
	aliasCommand := &cobra.Command{
		Use:   "alias",
		Short: "manage alias (different names of the commands)",
		PreRunE: func(c *cobra.Command, args []string) error {
			if aliasOpts.removeFlag && len(args) == 0 {
				return errors.New("remove flag requires arguments")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			config := rrh.OpenConfig()
			alias, err := LoadAliases(config)
			if err != nil {
				return err
			} else if len(args) == 0 {
				return listAlias(c, alias)
			} else if aliasOpts.removeFlag {
				return removeAliases(c, args, alias, config)
			} else {
				return registerAlias(c, args, alias, config)
			}
		},
	}
	flags := aliasCommand.Flags()
	flags.BoolVarP(&aliasOpts.removeFlag, "remove", "r", false, "remove the specified alias name")
	flags.BoolVarP(&aliasOpts.dryRunFlag, "dry-run", "D", false, "dry-run mode")

	return aliasCommand
}

var aliasOpts = &aliasOptions{}

type aliasOptions struct {
	removeFlag bool
	dryRunFlag bool
}

type Command struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func LoadAliases(config *rrh.Config) ([]*Command, error) {
	path := config.GetValue(rrh.AliasPath)
	alias := []*Command{}
	err := rrh.LoadJson(path, &alias)
	return alias, err
}

func storeAliases(aliasList []*Command, config *rrh.Config) error {
	path := config.GetValue(rrh.AliasPath)
	rrh.StoreJson(path, aliasList)
	return nil
}

func listAlias(cmd *cobra.Command, aliasList []*Command) error {
	for _, a := range aliasList {
		cmd.Printf("%s=%s\n", a.Name, strings.Join(a.Values, " "))
	}
	return nil
}

func removeAliases(cmd *cobra.Command, args []string, aliasList []*Command, config *rrh.Config) error {
	notFoundNames := []string{}
	foundNames := []string{}
	resultAliases := aliasList
	for _, arg := range args {
		r, err := removeIt(arg, aliasList)
		if err != nil {
			notFoundNames = append(notFoundNames, arg)
		} else {
			foundNames = append(foundNames, arg)
			cmd.Printf("remove %s from alias list", arg)
		}
		resultAliases = r
	}
	dryRunMode, err := cmd.Flags().GetBool("dry-run")
	if len(resultAliases) < len(aliasList) && err != nil && dryRunMode {
		storeAliases(resultAliases, config)
	}
	printDryRun(cmd, fmt.Sprintf("%s: remove alias (dry-run mode)", strings.Join(foundNames, ",")))
	return createError(notFoundNames)
}

func printDryRun(cmd *cobra.Command, message string) {
	dryRunMode, err := cmd.Flags().GetBool("dry-run")
	if dryRunMode && err != nil {
		cmd.Printf(message)
	}
}

func createError(names []string) error {
	switch len(names) {
	case 0:
		return nil
	case 1:
		return fmt.Errorf("%s: alias name not found", names[0])
	default:
		return fmt.Errorf("%s: alias names not found", strings.Join(names, ", "))
	}
}

func removeIt(aliasName string, aliasList []*Command) ([]*Command, error) {
	foundFlag := false
	resultList := []*Command{}
	for _, alias := range aliasList {
		if alias.Name == aliasName {
			foundFlag = true
		} else {
			resultList = append(resultList, alias)
		}
	}
	if !foundFlag {
		return resultList, fmt.Errorf("%s: alias not found", aliasName)
	}
	return resultList, nil
}

func FindAlias(name string, aliasList []*Command) *Command {
	for _, alias := range aliasList {
		if alias.Name == name {
			return alias
		}
	}
	return nil
}

func registerAlias(cmd *cobra.Command, args []string, aliasList []*Command, config *rrh.Config) error {
	if item := FindAlias(args[0], aliasList); item != nil {
		return fmt.Errorf("%s: already registered alias", args[0])
	}
	alias := &Command{Name: args[0], Values: args[1:]}
	newList := append(aliasList, alias)
	dryRunMode, err := cmd.Flags().GetBool("dry-run")
	if err == nil || !dryRunMode {
		return storeAliases(newList, config)
	}
	return nil
}

func (cmd *Command) Execute(c *cobra.Command, args []string) error {
	newArgs := cmd.Values
	newArgs = append(newArgs, args[1:]...)
	root := c.Root()
	root.SetArgs(newArgs)
	return root.Execute()
}
