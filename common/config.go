package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mitchellh/go-homedir"
)

/*
VERSION shows the version of RRH.
*/
const VERSION = "0.3"

/*
The environment variable names.
*/
const (
	RrhAutoDeleteGroup  = "RRH_AUTO_DELETE_GROUP"
	RrhAutoCreateGroup  = "RRH_AUTO_CREATE_GROUP"
	RrhCloneDestination = "RRH_CLONE_DESTINATION"
	RrhColor            = "RRH_COLOR"
	RrhConfigPath       = "RRH_CONFIG_PATH"
	RrhDatabasePath     = "RRH_DATABASE_PATH"
	RrhDefaultGroupName = "RRH_DEFAULT_GROUP_NAME"
	RrhHome             = "RRH_HOME"
	RrhOnError          = "RRH_ON_ERROR"
	RrhSortOnUpdating   = "RRH_SORT_ON_UPDATING"
	RrhTimeFormat       = "RRH_TIME_FORMAT"
)

var availableLabels = []string{
	RrhAutoCreateGroup, RrhAutoDeleteGroup, RrhCloneDestination, RrhColor,
	RrhConfigPath, RrhDatabasePath, RrhDefaultGroupName, RrhHome, RrhOnError,
	RrhSortOnUpdating, RrhTimeFormat,
}
var boolLabels = []string{
	RrhAutoCreateGroup, RrhAutoDeleteGroup, RrhSortOnUpdating,
}

/*
The the 2nd return value of GetString method (ReadFrom).
*/
const (
	Default    = "default"
	ConfigFile = "config_file"
	Env        = "environment"
	NotFound   = "not found"
)

/*
Relative shows the default value of RrhTimeFormat.
*/
const Relative = "relative"

const (
	trueString  = "true"
	falseString = "false"
)

/*
The values of RrhOnError.
*/
const (
	Fail            = "FAIL"
	FailImmediately = "FAIL_IMMEDIATELY"
	Warn            = "WARN"
	Ignore          = "IGNORE"
)

/*
Config shows the values of configuration variables.
*/
type Config map[string]string

/*
ReadFrom shows the value of config load from.
The available values are default, config_file, environment, and not found.
*/
type ReadFrom string

var defaultValues = Config{
	RrhAutoCreateGroup:  "false",
	RrhAutoDeleteGroup:  "false",
	RrhCloneDestination: ".",
	RrhColor:            "",
	RrhConfigPath:       "${RRH_HOME}/config.json",
	RrhDatabasePath:     "${RRH_HOME}/database.json",
	RrhDefaultGroupName: "no-group",
	RrhHome:             "${HOME}/.rrh",
	RrhOnError:          Warn,
	RrhSortOnUpdating:   "false",
	RrhTimeFormat:       Relative,
}

func (config *Config) isOnErrorIgnoreOrWarn() bool {
	var onError = config.GetValue(RrhOnError)
	return onError == Ignore || onError == Warn
}

/*
PrintErrors prints errors and returns the status code by following the value of RrhOnError.
If the value of RrhOnError is Ignore or Warn, this method returns 0, otherwise, non-zero value.
*/
func (config *Config) PrintErrors(errs []error) int {
	if config.GetValue(RrhOnError) != Ignore {
		for _, err := range errs {
			fmt.Println(err.Error())
		}
	}
	if len(errs) == 0 || config.isOnErrorIgnoreOrWarn() {
		return 0
	}
	return 5
}

func trueOrFalse(value string) (string, error) {
	if strings.ToLower(value) == trueString {
		return trueString, nil
	} else if strings.ToLower(value) == falseString {
		return falseString, nil
	}
	return "", fmt.Errorf("given value is not true nor false: %s", value)
}

func availableValueOnError(value string) (string, error) {
	var newvalue = strings.ToUpper(value)
	if newvalue == Fail || newvalue == FailImmediately || newvalue == Warn || newvalue == Ignore {
		return newvalue, nil
	}
	return "", fmt.Errorf("%s: Unknown value of RRH_ON_ERROR (must be %s, %s, %s, or %s)", value, Fail, FailImmediately, Warn, Ignore)
}

func contains(slice []string, label string) bool {
	for _, item := range slice {
		if label == item {
			return true
		}
	}
	return false
}

/*
Unset method deletes the specified config value.
*/
func (config *Config) Unset(label string) error {
	if !contains(availableLabels, label) {
		return fmt.Errorf("%s: unknown variable name", label)
	}
	delete((*config), label)
	return nil
}

func validateArgumentsOnUpdate(label string, value string) error {
	if !contains(availableLabels, label) {
		return fmt.Errorf("%s: unknown variable name", label)
	}
	if label == RrhConfigPath {
		return fmt.Errorf("%s: cannot set in config file", RrhConfigPath)
	}
	return nil
}

/*
Update method updates the config value with the given `value`.
*/
func (config *Config) Update(label string, value string) error {
	if err := validateArgumentsOnUpdate(label, value); err != nil {
		return err
	}
	if contains(boolLabels, label) {
		var flag, err = trueOrFalse(value)
		if err == nil {
			(*config)[label] = string(flag)
		}
		return err
	}
	if label == RrhOnError {
		var newValue, err = availableValueOnError(value)
		if err != nil {
			return err
		}
		value = newValue
	}
	(*config)[label] = value
	return nil
}

/*
IsSet returns the bool value of the given label.
If the label is not RrhAutoCreateGroup, RrhAutoDeleteGroup, and RrhSortOnUpdating, this method always returns false.
*/
func (config *Config) IsSet(label string) bool {
	if contains(boolLabels, label) {
		return strings.ToLower((*config)[label]) == trueString
	}
	return false
}

func (config *Config) replaceHome(value string) string {
	if strings.Contains(value, "${HOME}") {
		var home, _ = homedir.Dir()
		value = strings.Replace(value, "${HOME}", home, 1)
	}
	if strings.Contains(value, "${RRH_HOME}") {
		var rrhHome = config.GetValue(RrhHome)
		value = strings.Replace(value, "${RRH_HOME}", strings.TrimRight(rrhHome, "/"), -1)
	}
	return value
}

/*
GetValue returns the value of the given variable name.
*/
func (config *Config) GetValue(label string) string {
	var value, _ = config.GetString(label)
	return config.replaceHome(value)
}

/*
GetString returns the value of the given variable name and the definition part of the value.
*/
func (config *Config) GetString(label string) (string, ReadFrom) {
	if !contains(availableLabels, label) {
		return "", NotFound
	}
	var value, ok = (*config)[label]
	if !ok {
		return config.getStringFromEnv(label)
	}
	return config.replaceHome(value), ConfigFile
}

/*
GetDefaultValue returns the default value of the given variable name.
*/
func (config *Config) GetDefaultValue(label string) string {
	var value, _ = config.findDefaultValue(label)
	return value
}

func (config *Config) getStringFromEnv(label string) (string, ReadFrom) {
	var valueFromEnv = os.Getenv(label)
	if valueFromEnv != "" {
		valueFromEnv = config.replaceHome(valueFromEnv)
		return valueFromEnv, Env
	}
	return config.findDefaultValue(label)
}

func (config *Config) findDefaultValue(label string) (string, ReadFrom) {
	if !contains(availableLabels, label) {
		return "", NotFound
	}
	var value = defaultValues[label]
	value = config.replaceHome(value)
	return value, Default
}

func configPath() string {
	var config = new(Config)
	var configPath, _ = config.getStringFromEnv(RrhConfigPath)
	return configPath
}

/*
StoreConfig saves the store.
*/
func (config *Config) StoreConfig() error {
	var configPath = config.GetValue(RrhConfigPath)
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

/*
OpenConfig reads the config file and returns it.
The load path is based on `RrhConfigPath` of the environment variables.
*/
func OpenConfig() *Config {
	bytes, err := ioutil.ReadFile(configPath())
	if err != nil {
		return new(Config)
	}
	var config Config
	if err := json.Unmarshal(bytes, &config); err != nil {
		return nil
	}
	return &config
}

func (config *Config) formatVariableAndValue(label string) string {
	var value, readFrom = config.GetString(label)
	return fmt.Sprintf("%s: %s (%s)", label, value, readFrom)
}
