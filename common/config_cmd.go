package common

import (
	"fmt"
	"log"

	"github.com/mitchellh/cli"
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
	c := cli.NewCLI("rrh config", VERSION)
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
	var exitStatus, err = c.Run()
	if err != nil {
		log.Println(err)
	}
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
	var config = OpenConfig()
	var err = config.Update(args[0], args[1])
	if err != nil {
		fmt.Println(err.Error())
		return 1
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
	var config = OpenConfig()
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
func (clc *configListCommand) Run(args []string) int {
	var config = OpenConfig()
	fmt.Println(config.formatVariableAndValue(RrhHome))
	fmt.Println(config.formatVariableAndValue(RrhConfigPath))
	fmt.Println(config.formatVariableAndValue(RrhDatabasePath))
	fmt.Println(config.formatVariableAndValue(RrhDefaultGroupName))
	fmt.Println(config.formatVariableAndValue(RrhOnError))
	fmt.Println(config.formatVariableAndValue(RrhAutoCreateGroup))
	fmt.Println(config.formatVariableAndValue(RrhAutoDeleteGroup))
	fmt.Println(config.formatVariableAndValue(RrhTimeFormat))
	fmt.Println(config.formatVariableAndValue(RrhCloneDestination))
	fmt.Println(config.formatVariableAndValue(RrhSortOnUpdating))
	return 0
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
