package rrh

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
const VERSION = "1.2.0"

/*
The environment variable names.
*/
const (
	AutoDeleteGroup  = "RRH_AUTO_DELETE_GROUP"
	AutoCreateGroup  = "RRH_AUTO_CREATE_GROUP"
	CloneDestination = "RRH_CLONE_DESTINATION"
	ColorSetting     = "RRH_COLOR"
	ConfigPath       = "RRH_CONFIG_PATH"
	DatabasePath     = "RRH_DATABASE_PATH"
	DefaultGroupName = "RRH_DEFAULT_GROUP_NAME"
	EnableColorized  = "RRH_ENABLE_COLORIZED"
	Home             = "RRH_HOME"
	OnError          = "RRH_ON_ERROR"
	SortOnUpdating   = "RRH_SORT_ON_UPDATING"
	TimeFormat       = "RRH_TIME_FORMAT"
)

/*
AvailableLabels represents the labels availables in the config.
*/
var AvailableLabels = []string{
	AutoCreateGroup, AutoDeleteGroup, CloneDestination, ColorSetting,
	ConfigPath, DatabasePath, DefaultGroupName, EnableColorized,
	Home, OnError, SortOnUpdating, TimeFormat,
}
var boolLabels = []string{
	AutoCreateGroup, AutoDeleteGroup, EnableColorized,
	SortOnUpdating,
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
		AutoCreateGroup:  "false",
		AutoDeleteGroup:  "false",
		CloneDestination: ".",
		ColorSetting:     "repository:fg=red+group:fg=magenta+label:op=bold+configValue:fg=green",
		ConfigPath:       "${RRH_HOME}/config.json",
		DatabasePath:     "${RRH_HOME}/database.json",
		DefaultGroupName: "no-group",
		EnableColorized:  "false",
		Home:             "${HOME}/.rrh",
		OnError:          Warn,
		SortOnUpdating:   "false",
		TimeFormat:       Relative,
	},
	Color: &Color{},
}

func (config *Config) isOnErrorIgnoreOrWarn() bool {
	var onError = config.GetValue(OnError)
	return onError == Ignore || onError == Warn
}

func printErrorImpl(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func isErrorOrIgnore(errs []error, config *Config) bool {
	return len(errs) == 0 || errs[0] == nil || config.isOnErrorIgnoreOrWarn()
}

/*
PrintErrors prints errors and returns the status code by following the value of RrhOnError.
If the value of RrhOnError is Ignore or Warn, this method returns 0, otherwise, non-zero value.
*/
func (config *Config) PrintErrors(errs ...error) int {
	if config.GetValue(OnError) != Ignore {
		for _, err := range errs {
			printErrorImpl(err)
		}
	}
	if isErrorOrIgnore(errs, config) {
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
	if label == ConfigPath {
		return fmt.Errorf("%s: cannot set in config file", ConfigPath)
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
	if label == OnError {
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
		var rrhHome = config.GetValue(Home)
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
	var configPath = config.GetValue(ConfigPath)
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
	var configPath, _ = config.getStringFromEnv(ConfigPath)
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
