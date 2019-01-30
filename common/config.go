package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mitchellh/cli"
)

const (
	GrimHome             = "GRIM_HOME"
	GrimConfigPath       = "GRIM_CONFIG_PATH"
	GrimDatabasePath     = "GRIM_DATABASE_PATH"
	GrimDefaultGroupName = "GRIM_DEFAULT_GROUP_NAME"
	GrimOnError          = "GRIM_ON_ERROR"
	GrimAutoDeleteGroup  = "GRIM_AUTO_DELETE_GROUP"
	GrimAutoCreateGroup  = "GRIM_AUTO_CREATE_GROUP"
)

const (
	Default    = "default"
	ConfigFile = "config_file"
	Env        = "environment"
	NotFound   = "not found"
)

const (
	Fail            = "FAIL"
	FailImmediately = "FAIL_IMMEDIATELY"
	Warn            = "WARN"
	Ignore          = "IGNORE"
)

type Config struct {
	Home             string `json:"grim_home"`
	AutoDeleteGroup  string `json:"grim_auto_delete_group"`
	AutoCreateGroup  string `json:"grim_auto_create_group"`
	ConfigPath       string `json:"grim_config_path"`
	DatabasePath     string `json:"grim_database_path"`
	DefaultGroupName string `json:"grim_default_group_name"`
	OnError          string `json:"grim_on_error"`
}

func trueOrFalse(value string) (string, error) {
	if strings.ToLower(value) == "true" {
		return "true", nil
	} else if strings.ToLower(value) == "false" {
		return "false", nil
	}
	return "", fmt.Errorf("given value is not true nor false: %s", value)
}

func availableValueOnError(value string) (string, error) {
	var newvalue = strings.ToUpper(value)
	if newvalue == Fail || newvalue == FailImmediately || newvalue == Warn || newvalue == Ignore {
		return newvalue, nil
	}
	return "", fmt.Errorf("%s: Unknown value of GRIM_ON_ERROR (must be %s, %s, %s, or %s)", value, Fail, FailImmediately, Warn, Ignore)
}

func (config *Config) Update(label string, value string) error {
	switch label {
	case GrimAutoDeleteGroup:
		var flag, err = trueOrFalse(value)
		if err == nil {
			config.AutoDeleteGroup = flag
		}
		return err
	case GrimAutoCreateGroup:
		var flag, err = trueOrFalse(value)
		if err == nil {
			config.AutoCreateGroup = flag
		}
		return err
	case GrimHome:
		config.Home = value
		return nil
	case GrimDatabasePath:
		config.DatabasePath = value
		return nil
	case GrimDefaultGroupName:
		config.DefaultGroupName = value
		return nil
	case GrimOnError:
		var newValue, err = availableValueOnError(value)
		if err == nil {
			config.OnError = newValue
		}
		return err
	case GrimConfigPath:
		return fmt.Errorf("%s: does not set on config file. set on environment.", label)
	}
	return fmt.Errorf("%s: Unknown variable name", label)
}

func (config *Config) GetValue(label string) string {
	var value, _ = config.GetString(label)
	return value
}

func (config *Config) GetDefaultValue(label string) string {
	var value, _ = config.findDefaultValue(label, "")
	return value
}

func (config *Config) GetString(label string) (value string, readFrom string) {
	switch label {
	case GrimAutoDeleteGroup:
		return config.getStringFromEnv(GrimAutoDeleteGroup, config.AutoDeleteGroup)
	case GrimAutoCreateGroup:
		return config.getStringFromEnv(GrimAutoCreateGroup, config.AutoCreateGroup)
	case GrimHome:
		return config.getStringFromEnv(GrimHome, config.Home)
	case GrimConfigPath:
		return config.getStringFromEnv(GrimConfigPath, config.ConfigPath)
	case GrimDefaultGroupName:
		return config.getStringFromEnv(GrimDefaultGroupName, config.DefaultGroupName)
	case GrimDatabasePath:
		return config.getStringFromEnv(GrimDatabasePath, config.DatabasePath)
	case GrimOnError:
		return config.getStringFromEnv(GrimOnError, config.OnError)
	default:
		return "", NotFound
	}
}

func (config *Config) getStringFromEnv(label string, valueFromConfigFile string) (value string, readFrom string) {
	if valueFromConfigFile != "" {
		return valueFromConfigFile, ConfigFile
	}
	return config.findDefaultValue(label, os.Getenv(label))
}

func (config *Config) findDefaultValue(label string, valueFromEnv string) (value string, readFrom string) {
	if valueFromEnv != "" {
		return valueFromEnv, Env
	}
	switch label {
	case GrimHome:
		return fmt.Sprintf("%s/.grim", os.Getenv("HOME")), Default
	case GrimConfigPath:
		var home, _ = config.GetString(GrimHome)
		return fmt.Sprintf("%s/config.json", home), Default
	case GrimDatabasePath:
		var home, _ = config.GetString(GrimHome)
		return fmt.Sprintf("%s/database.json", home), Default
	case GrimDefaultGroupName:
		return "no-group", Default
	case GrimOnError:
		return Warn, Default
	case GrimAutoDeleteGroup:
		return "false", Default
	case GrimAutoCreateGroup:
		return "false", Default
	default:
		return "", NotFound
	}
}

func configPath() string {
	var config = new(Config)
	var home, _ = config.getStringFromEnv(GrimConfigPath, "")
	return home
}

func (config *Config) StoreConfig() error {
	var configPath = config.GetValue(GrimConfigPath)
	var err1 = CreateParentDir(configPath)
	if err1 != nil {
		return err1
	}
	var bytes, err2 = json.Marshal(*config)
	if err2 == nil {
		return ioutil.WriteFile(configPath, bytes, 0644)
	}
	return err2
}

func OpenConfig() *Config {
	bytes, err := ioutil.ReadFile(configPath())
	if err != nil {
		return new(Config)
	}
	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		log.Fatal(err)
	}
	return &config
}

func printVariableAndValue(label string, config *Config) {
	var value, readFrom = config.GetString(label)
	fmt.Printf("%s: %s (%s)\n", label, value, readFrom)
}

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
	return `grim config <COMMAND> [ARGUMENTS]
COMMAND
	set <ENV_NAME> <VALUE>  set ENV_NAME to VALUE
	unset <ENV_NAME>        reset ENV_NAME
	list                    list all of ENVs`
}

func (csc *configSetCommand) Help() string {
	return `grim config set <ENV_NAME> <VALUE>
ARGUMENTS
	ENV_NAME   environment name.
	VALUE      the value for the given environment.`
}

func (cuc *configUnsetCommand) Help() string {
	return `grim config unset <ENV_NAME>
ARGUMENTS
	ENV_NAME   environment name.`
}

func (clc *configListCommand) Help() string {
	return `grim config list`
}

func (config *ConfigCommand) Run(args []string) int {
	c := cli.NewCLI("grim config", "1.0.0")
	c.Args = args
	c.Autocomplete = true
	c.Commands = map[string]cli.CommandFactory{
		"set":   configSetCommandFactory,
		"unset": configUnsetCommandFactory,
		"list":  configListCommandFactory,
	}
	var exitStatus, err = c.Run()
	if err != nil {
		log.Println(err)
	}
	return exitStatus
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
	var err = config.Update(args[0], config.GetDefaultValue(args[0]))
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	config.StoreConfig()
	return 0
}

func (clc *configListCommand) Run(args []string) int {
	var config = OpenConfig()
	printVariableAndValue(GrimHome, config)
	printVariableAndValue(GrimConfigPath, config)
	printVariableAndValue("GRIM_DATABASE_PATH", config)
	printVariableAndValue("GRIM_DEFAULT_GROUP_NAME", config)
	printVariableAndValue("GRIM_ON_ERROR", config)
	printVariableAndValue("GRIM_AUTO_CREATE_GROUP", config)
	printVariableAndValue("GRIM_AUTO_DELETE_GROUP", config)
	return 0
}

/*
Synopsis returns the messages for help on `grim bconfig set`.
*/
func (csc *configSetCommand) Synopsis() string {
	return "set the environment with the given value."
}

/*
Synopsis returns the messages for help on `grim bconfig unset`.
*/
func (cuc *configUnsetCommand) Synopsis() string {
	return "reset the given environment."
}

/*
Synopsis returns the messages for help on `grim bconfig list`.
*/
func (clc *configListCommand) Synopsis() string {
	return "list the environment and its value."
}

func (config *ConfigCommand) Synopsis() string {
	return "set/unset and list configuration of grim."
}
