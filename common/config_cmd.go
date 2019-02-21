package common

import (
	"fmt"
	"log"

	"github.com/mitchellh/cli"
)

type ConfigCommand struct{}
type configSetCommand struct{}
type configUnsetCommand struct{}
type configListCommand struct{}

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

func (config *ConfigCommand) Help() string {
	return `rrh config <COMMAND> [ARGUMENTS]
COMMAND
    set <ENV_NAME> <VALUE>  set ENV_NAME to VALUE
    unset <ENV_NAME>        reset ENV_NAME
    list                    list all of ENVs (default)`
}

func (csc *configSetCommand) Help() string {
	return `rrh config set <ENV_NAME> <VALUE>
ARGUMENTS
    ENV_NAME   environment name.
    VALUE      the value for the given environment.`
}

func (cuc *configUnsetCommand) Help() string {
	return `rrh config unset <ENV_NAME...>
ARGUMENTS
    ENV_NAME   environment name.`
}

func (clc *configListCommand) Help() string {
	return `rrh config list`
}

func (config *ConfigCommand) Run(args []string) int {
	c := cli.NewCLI("rrh config", "1.0.0")
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
	} else {
		var exitStatus, err = c.Run()
		if err != nil {
			log.Println(err)
		}
		return exitStatus
	}
}

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

func (cuc *configUnsetCommand) Run(args []string) int {
	if len(args) != 1 {
		fmt.Println(cuc.Help())
		return 1
	}
	var config = OpenConfig()
	for _, arg := range args {
		var err = config.Update(arg, config.GetDefaultValue(arg))
		if err != nil {
			fmt.Println(err.Error())
			return 1
		}
	}
	config.StoreConfig()
	return 0
}

func (clc *configListCommand) Run(args []string) int {
	var config = OpenConfig()
	fmt.Println(config.formatVariableAndValue(RrhHome))
	fmt.Println(config.formatVariableAndValue(RrhConfigPath))
	fmt.Println(config.formatVariableAndValue(RrhDatabasePath))
	fmt.Println(config.formatVariableAndValue(RrhDefaultGroupName))
	fmt.Println(config.formatVariableAndValue(RrhOnError))
	fmt.Println(config.formatVariableAndValue(RrhTimeFormat))
	fmt.Println(config.formatVariableAndValue(RrhAutoCreateGroup))
	fmt.Println(config.formatVariableAndValue(RrhAutoDeleteGroup))
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
