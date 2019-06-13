package internal

import (
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh/lib"
)

/*
ConfigCommand represents a command.
*/
type ConfigCommand struct{}
type setCommand struct{}
type unsetCommand struct{}
type listCommand struct{}

/*
ConfigCommandFactory returns an instance of the ConfigCommand.
*/
func ConfigCommandFactory() (cli.Command, error) {
	return &ConfigCommand{}, nil
}

func setCommandFactory() (cli.Command, error) {
	return &setCommand{}, nil
}

func unsetCommandFactory() (cli.Command, error) {
	return &unsetCommand{}, nil
}

func listCommandFactory() (cli.Command, error) {
	return &listCommand{}, nil
}

/*
Help returns the help message.
*/
func (config *ConfigCommand) Help() string {
	return `rrh config <COMMAND> [ARGUMENTS]
COMMAND
    set <ENV_NAME> <VALUE>  set ENV_NAME to VALUE
    unset <ENV_NAME>        reset ENV_NAME
    list                    list all of ENVs (default)`
}

/*
Help returns the help message.
*/
func (csc *setCommand) Help() string {
	return `rrh config set <ENV_NAME> <VALUE>
ARGUMENTS
    ENV_NAME   environment name.
    VALUE      the value for the given environment.`
}

/*
Help returns the help message.
*/
func (cuc *unsetCommand) Help() string {
	return `rrh config unset <ENV_NAME...>
ARGUMENTS
    ENV_NAME   environment name.`
}

/*
Help returns the help message.
*/
func (clc *listCommand) Help() string {
	return `rrh config list`
}

/*
Run performs the command.
*/
func (config *ConfigCommand) Run(args []string) int {
	c := cli.NewCLI("rrh config", lib.VERSION)
	c.Args = args
	c.Autocomplete = true
	c.Commands = map[string]cli.CommandFactory{
		"set":   setCommandFactory,
		"unset": unsetCommandFactory,
		"list":  listCommandFactory,
	}
	if len(args) == 0 {
		new(listCommand).Run([]string{})
		return 0
	}
	var exitStatus, _ = c.Run()
	return exitStatus
}

/*
Run performs the command.
*/
func (csc *setCommand) Run(args []string) int {
	if len(args) != 2 {
		fmt.Println(csc.Help())
		return 1
	}
	var config = lib.OpenConfig()
	var err = config.Update(args[0], args[1])
	if err != nil {
		fmt.Println(err.Error())
		return 2
	}
	config.StoreConfig()
	return 0
}

/*
Run performs the command.
*/
func (cuc *unsetCommand) Run(args []string) int {
	if len(args) != 1 {
		fmt.Println(cuc.Help())
		return 1
	}
	var config = lib.OpenConfig()
	var err = config.Unset(args[0])
	if err != nil {
		var status = config.PrintErrors([]error{err})
		if status != 0 {
			return status
		}
	}
	config.StoreConfig()
	return 0
}

/*
Run performs the command.
*/
func (clc *listCommand) Run(args []string) int {
	var config = lib.OpenConfig()
	for _, label := range lib.AvailableLabels {
		fmt.Println(formatVariableAndValue(config, label))
	}
	return 0
}

func formatVariableAndValue(config *lib.Config, label string) string {
	var value, readFrom = config.GetString(label)
	return fmt.Sprintf("%s: %s (%s)",
		config.Color.ColorizedLabel(label), config.Color.ColorizeConfigValue(value), readFrom)
}

/*
Synopsis returns the help message of the command.
*/
func (csc *setCommand) Synopsis() string {
	return "set the environment with the given value."
}

/*
Synopsis returns the help message of the command.
*/
func (cuc *unsetCommand) Synopsis() string {
	return "reset the given environment."
}

/*
Synopsis returns the help message of the command.
*/
func (clc *listCommand) Synopsis() string {
	return "list the environment and its value."
}

/*
Synopsis returns the help message of the command.
*/
func (config *ConfigCommand) Synopsis() string {
	return "set/unset and list configuration of RRH."
}