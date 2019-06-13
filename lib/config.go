package lib

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
const VERSION = "0.5"

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
	RrhEnableColorized  = "RRH_ENABLE_COLORIZED"
	RrhHome             = "RRH_HOME"
	RrhOnError          = "RRH_ON_ERROR"
	RrhSortOnUpdating   = "RRH_SORT_ON_UPDATING"
	RrhTimeFormat       = "RRH_TIME_FORMAT"
)

/*
AvailableLabels represents the labels availables in the config.
*/
var AvailableLabels = []string{
	RrhAutoCreateGroup, RrhAutoDeleteGroup, RrhCloneDestination, RrhColor,
	RrhConfigPath, RrhDatabasePath, RrhDefaultGroupName, RrhEnableColorized,
	RrhHome, RrhOnError, RrhSortOnUpdating, RrhTimeFormat,
}
var boolLabels = []string{
	RrhAutoCreateGroup, RrhAutoDeleteGroup, RrhEnableColorized,
	RrhSortOnUpdating,
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
type Config struct {
	values map[string]string
	Color  *Color
}

/*
ReadFrom shows the value of config load from.
The available values are default, config_file, environment, and not found.
*/
type ReadFrom string

var defaultValues = Config{
	values: map[string]string{
		RrhAutoCreateGroup:  "false",
		RrhAutoDeleteGroup:  "false",
		RrhCloneDestination: ".",
		RrhColor:            "repository:fg=red+group:fg=magenta+label:op=bold+configValue:fg=green",
		RrhConfigPath:       "${RRH_HOME}/config.json",
		RrhDatabasePath:     "${RRH_HOME}/database.json",
		RrhDefaultGroupName: "no-group",
		RrhEnableColorized:  "false",
		RrhHome:             "${HOME}/.rrh",
		RrhOnError:          Warn,
		RrhSortOnUpdating:   "false",
		RrhTimeFormat:       Relative,
	},
	Color: &Color{},
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
	return "", fmt.Errorf("%s: not true nor false", value)
}

func normalizeValueOfOnError(value string) (string, error) {
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
	if !contains(AvailableLabels, label) {
		return fmt.Errorf("%s: unknown variable name", label)
	}
	delete(config.values, label)
	return nil
}

func validateArgumentsOnUpdate(label string, value string) error {
	if !contains(AvailableLabels, label) {
		return fmt.Errorf("%s: unknown variable name", label)
	}
	if label == RrhConfigPath {
		return fmt.Errorf("%s: cannot set in config file", RrhConfigPath)
	}
	return nil
}

func (config *Config) updateBoolValue(label string, value string) error {
	var flag, err = trueOrFalse(value)
	if err == nil {
		config.values[label] = string(flag)
	}
	return err
}

/*
Update method updates the config value with the given `value`.
*/
func (config *Config) Update(label string, value string) error {
	if err := validateArgumentsOnUpdate(label, value); err != nil {
		return err
	}
	if contains(boolLabels, label) {
		return config.updateBoolValue(label, value)
	}
	if label == RrhOnError {
		var newValue, err = normalizeValueOfOnError(value)
		if err != nil {
			return err
		}
		value = newValue
	}
	config.values[label] = value
	return nil
}

/*
IsSet returns the bool value of the given label.
If the label is not RrhAutoCreateGroup, RrhAutoDeleteGroup, and RrhSortOnUpdating, this method always returns false.
*/
func (config *Config) IsSet(label string) bool {
	if contains(boolLabels, label) {
		return strings.ToLower(config.values[label]) == trueString
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
	if !contains(AvailableLabels, label) {
		return "", NotFound
	}
	var value, ok = config.values[label]
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
	if !contains(AvailableLabels, label) {
		return "", NotFound
	}
	var value = defaultValues.values[label]
	value = config.replaceHome(value)
	return value, Default
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
	var bytes, err2 = json.Marshal(config.values)
	if err2 != nil {
		return err2
	}
	return ioutil.WriteFile(configPath, bytes, 0644)
}

/*
NewConfig generates the new Config instance.
*/
func NewConfig() *Config {
	return &Config{values: map[string]string{}, Color: &Color{colorSettings{}, colorFuncs{}}}
}

/*
OpenConfig reads the config file and returns it.
The load path is based on `RrhConfigPath` of the environment variables.
*/
func OpenConfig() *Config {
	var config = NewConfig()
	var configPath, _ = config.getStringFromEnv(RrhConfigPath)
	bytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config
	}
	var values = map[string]string{}
	if err := json.Unmarshal(bytes, &values); err != nil {
		return nil
	}
	config.values = values
	config.Color = InitializeColor(config)
	return config
}
