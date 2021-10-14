package config

import (
	"github.com/spf13/cobra"
	"github.com/tamada/rrh"
	"github.com/tamada/rrh/common"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "set/unset and list configuration of RRH",
		Args:  cobra.ArbitraryArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return listConfig(c)
		},
	}
	cmd.AddCommand(newListCommand())
	cmd.AddCommand(newSetCommand())
	cmd.AddCommand(newUnsetCommand())
	return cmd
}

func newUnsetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unset <KEYS...>",
		Short: "unset the environment values of the given keys",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			return unsetConfig(c, args)
		},
	}
	return cmd
}

func newSetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set <KEY> <VALUE>",
		Short: "set the environment value with the given value",
		Args:  cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			return setConfig(c, args)
		},
	}
	return cmd
}

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list the environment and its value",
		Args:  cobra.NoArgs,
		RunE: func(c *cobra.Command, args []string) error {
			return listConfig(c)
		},
	}
	return cmd
}

func setConfig(c *cobra.Command, args []string) error {
	config := rrh.OpenConfig()
	err := config.Update(args[0], args[1])
	if err != nil {
		return err
	}
	config.StoreConfig()
	return nil
}

func unsetConfig(c *cobra.Command, args []string) error {
	config := rrh.OpenConfig()
	el := common.NewErrorList()
	for _, key := range args {
		err := config.Unset(key)
		el = el.Append(err)
	}
	config.StoreConfig()
	return el.NilOrThis()
}

func listConfig(c *cobra.Command) error {
	var config = rrh.OpenConfig()
	for _, label := range rrh.AvailableLabels {
		value, readFrom := config.GetString(label)
		c.Printf("%s: %s (%s)\n", label, value, readFrom)
	}
	return nil
}
