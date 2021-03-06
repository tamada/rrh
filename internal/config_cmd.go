package internal

import (
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/tamada/rrh"
)

/*
ConfigCommand represents a command.
*/
type ConfigCommand struct{}
type configSetCommand struct{}
type configUnsetCommand struct{}
type configListCommand struct{}

/*
ConfigCommandFactory returns an instance of the ConfigCommand.
*/
func ConfigCommandFactory() (cli.Command, error) {
	return &ConfigCommand{}, nil
}

func configSetCommandFactory() (cli.Command, error) {
	return &configSetCommand{}, nil
}

func configUnsetCommandFactory() (cli.Command, error) {
	return &configUnsetCommand{}, nil
}

func configListCommandFactory() (cli.Command, error) {
	return &configListCommand{}, nil
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
func (csc *configSetCommand) Help() string {
	return `rrh config set <ENV_NAME> <VALUE>
ARGUMENTS
    ENV_NAME   environment name.
    VALUE      the value for the given environment.`
}

/*
Help returns the help message.
*/
func (cuc *configUnsetCommand) Help() string {
	return `rrh config unset <ENV_NAME...>
ARGUMENTS
    ENV_NAME   environment name.`
}

/*
Help returns the help message.
*/
func (clc *configListCommand) Help() string {
	return `rrh config list`
}

/*
Run performs the command.
*/
func (config *ConfigCommand) Run(args []string) int {
	c := cli.NewCLI("rrh config", rrh.VERSION)
	c.Args = args
	c.Autocomplete = true
	c.Commands = map[string]cli.CommandFactory{
		"set":   configSetCommandFactory,
		"unset": configUnsetCommandFactory,
		"list":  configListCommandFactory,
	}
	if len(args) == 0 {
		new(configListCommand).Run([]string{})
		return 0
	}
	var exitStatus, _ = c.Run()
	return exitStatus
}

/*
Run performs the command.
*/
func (csc *configSetCommand) Run(args []string) int {
	if len(args) != 2 {
		fmt.Println(csc.Help())
		return 1
	}
	var config = rrh.OpenConfig()
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
func (cuc *configUnsetCommand) Run(args []string) int {
	if len(args) != 1 {
		fmt.Println(cuc.Help())
		return 1
	}
	var config = rrh.OpenConfig()
	var err = config.Unset(args[0])
	if err != nil {
		var status = config.PrintErrors(err)
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
func (clc *configListCommand) Run(args []string) int {
	var config = rrh.OpenConfig()
	for _, label := range rrh.AvailableLabels {
		fmt.Println(formatVariableAndValue(config, label))
	}
	return 0
}

func formatVariableAndValue(config *rrh.Config, label string) string {
	var value, readFrom = config.GetString(label)
	return fmt.Sprintf("%s: %s (%s)",
		config.Color.ColorizedLabel(label), config.Color.ColorizeConfigValue(value), readFrom)
}

/*
Synopsis returns the help message of the command.
*/
func (csc *configSetCommand) Synopsis() string {
	return "set the environment with the given value."
}

/*
Synopsis returns the help message of the command.
*/
func (cuc *configUnsetCommand) Synopsis() string {
	return "reset the given environment."
}

/*
Synopsis returns the help message of the command.
*/
func (clc *configListCommand) Synopsis() string {
	return "list the environment and its value."
}

/*
Synopsis returns the help message of the command.
*/
func (config *ConfigCommand) Synopsis() string {
	return "set/unset and list configuration of RRH."
}
