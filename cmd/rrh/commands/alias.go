package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tamada/rrh"
)

func AliasCommand() *cobra.Command {
	aliasCommand := &cobra.Command{
		Use:   "alias",
		Short: "manages alias",
		RunE: func(c *cobra.Command, args []string) error {
			fmt.Printf("alias args: %v\n", args)
			alias, err := loadAliases()
			if err != nil {
				return err
			} else if len(args) == 0 {
				return listAlias(c, alias)
			} else if v, err := c.Flags().GetBool("remove"); v && err != nil {
				return removeAliases(c, args, alias)
			} else {
				return registerAlias(c, args, alias)
			}
		},
	}
	flags := aliasCommand.Flags()
	flags.BoolP("remove", "r", false, "remove the specified alias name")
	flags.BoolP("dry-run", "d", false, "dry-run mode")

	return aliasCommand
}

type Alias struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func loadAliases() ([]*Alias, error) {
	path := viper.GetString("alias_file_path")
	alias := []*Alias{}
	err := rrh.LoadJson(path, &alias)
	return alias, err
}

func storeAliases(aliasList []*Alias) error {
	path := viper.GetString("alias_file_path")
	rrh.StoreJson(path, aliasList)
	return nil
}

func listAlias(cmd *cobra.Command, aliasList []*Alias) error {
	for _, a := range aliasList {
		cmd.Printf("%s=%s\n", a.Name, strings.Join(a.Values, " "))
	}
	return nil
}

func removeAliases(cmd *cobra.Command, args []string, aliasList []*Alias) error {
	notFoundNames := []string{}
	foundNames := []string{}
	resultAliases := aliasList
	for _, arg := range args {
		r, err := removeIt(arg, aliasList)
		if err != nil {
			notFoundNames = append(notFoundNames, arg)
		} else {
			foundNames = append(foundNames, arg)
			rrh.PrintIfVerbose(cmd, fmt.Sprintf("remove %s from alias list", arg))
		}
		resultAliases = r
	}
	dryRunMode, err := cmd.Flags().GetBool("dry-run")
	if len(resultAliases) < len(aliasList) && err != nil && dryRunMode {
		storeAliases(resultAliases)
	}
	printDryRun(cmd, fmt.Sprintf("%s: remove alias (dry run mode)", strings.Join(foundNames, ",")))
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

func removeIt(aliasName string, aliasList []*Alias) ([]*Alias, error) {
	foundFlag := false
	resultList := []*Alias{}
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

func findAlias(name string, aliasList []*Alias) *Alias {
	for _, alias := range aliasList {
		if alias.Name == name {
			return alias
		}
	}
	return nil
}

func registerAlias(cmd *cobra.Command, args []string, aliasList []*Alias) error {
	if item := findAlias(args[0], aliasList); item != nil {
		return fmt.Errorf("%s: already registered alias", args[0])
	}
	alias := &Alias{Name: args[0], Values: args[1:]}
	newList := append(aliasList, alias)
	dryRunMode, err := cmd.Flags().GetBool("dry-run")
	if err == nil || !dryRunMode {
		return storeAliases(newList)
	}
	return nil
}

func executeAlias(cmd *cobra.Command, args []string, alias *Alias) error {
	newArgs := alias.Values
	newArgs = append(newArgs, args[1:]...)
	root := cmd.Root()
	root.SetArgs(newArgs)
	return root.Execute()
}
