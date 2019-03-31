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
const VERSION = "0.2"

const (
	RrhHome             = "RRH_HOME"
	RrhConfigPath       = "RRH_CONFIG_PATH"
	RrhDatabasePath     = "RRH_DATABASE_PATH"
	RrhDefaultGroupName = "RRH_DEFAULT_GROUP_NAME"
	RrhOnError          = "RRH_ON_ERROR"
	RrhAutoDeleteGroup  = "RRH_AUTO_DELETE_GROUP"
	RrhAutoCreateGroup  = "RRH_AUTO_CREATE_GROUP"
	RrhTimeFormat       = "RRH_TIME_FORMAT"
	RrhSortOnUpdating   = "RRH_SORT_ON_UPDATING"
)

const (
	Default    = "default"
	ConfigFile = "config_file"
	Env        = "environment"
	NotFound   = "not found"
	Relative   = "relative"
)

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
	Home             string `json:"rrh_home"`
	AutoDeleteGroup  string `json:"rrh_auto_delete_group"`
	AutoCreateGroup  string `json:"rrh_auto_create_group"`
	ConfigPath       string `json:"rrh_config_path"`
	DatabasePath     string `json:"rrh_database_path"`
	DefaultGroupName string `json:"rrh_default_group_name"`
	TimeFormat       string `json:"rrh_time_format"`
	OnError          string `json:"rrh_on_error"`
	SortOnUpdating   string `json:"rrh_sort_on_updating"`
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
	return "", fmt.Errorf("%s: Unknown value of RRH_ON_ERROR (must be %s, %s, %s, or %s)", value, Fail, FailImmediately, Warn, Ignore)
}

/*
Update method updates the config value with the given `value`.
*/
func (config *Config) Update(label string, value string) error {
	switch label {
	case RrhAutoDeleteGroup:
		var flag, err = trueOrFalse(value)
		if err == nil {
			config.AutoDeleteGroup = flag
		}
		return err
	case RrhAutoCreateGroup:
		var flag, err = trueOrFalse(value)
		if err == nil {
			config.AutoCreateGroup = flag
		}
		return err
	case RrhSortOnUpdating:
		var flag, err = trueOrFalse(value)
		if err == nil {
			config.SortOnUpdating = flag
		}
		return err
	case RrhHome:
		config.Home = value
		return nil
	case RrhTimeFormat:
		config.TimeFormat = value
		return nil
	case RrhDatabasePath:
		config.DatabasePath = value
		return nil
	case RrhDefaultGroupName:
		config.DefaultGroupName = value
		return nil
	case RrhOnError:
		var newValue, err = availableValueOnError(value)
		if err == nil {
			config.OnError = newValue
		}
		return err
	case RrhConfigPath:
		return fmt.Errorf("%s: does not set on config file. set on environment", label)
	}
	return fmt.Errorf("%s: Unknown variable name", label)
}

/*
IsSet returns the bool value of the given label.
If the label is not RrhAutoCreateGroup, RrhAutoDeleteGroup, and RrhSortOnUpdating, this method always returns false.
*/
func (config *Config) IsSet(label string) bool {
	var value = config.GetValue(label)
	if label != RrhAutoCreateGroup && label != RrhAutoDeleteGroup && label != RrhSortOnUpdating {
		return false
	}
	return strings.ToLower(value) == "true"
}

/*
GetValue returns the value of the given variable name.
*/
func (config *Config) GetValue(label string) string {
	var value, _ = config.GetString(label)
	return value
}

/*
GetDefaultValue returns the default value of the given variable name.
*/
func (config *Config) GetDefaultValue(label string) string {
	var value, _ = config.findDefaultValue(label)
	return value
}

/*
GetString returns the value of the given variable name and the definition part of the value.
*/
func (config *Config) GetString(label string) (value string, readFrom string) {
	switch label {
	case RrhAutoDeleteGroup:
		return config.getStringFromEnv(RrhAutoDeleteGroup, config.AutoDeleteGroup)
	case RrhAutoCreateGroup:
		return config.getStringFromEnv(RrhAutoCreateGroup, config.AutoCreateGroup)
	case RrhSortOnUpdating:
		return config.getStringFromEnv(RrhSortOnUpdating, config.SortOnUpdating)
	case RrhHome:
		return config.getStringFromEnv(RrhHome, config.Home)
	case RrhConfigPath:
		return config.getStringFromEnv(RrhConfigPath, config.ConfigPath)
	case RrhDefaultGroupName:
		return config.getStringFromEnv(RrhDefaultGroupName, config.DefaultGroupName)
	case RrhDatabasePath:
		return config.getStringFromEnv(RrhDatabasePath, config.DatabasePath)
	case RrhTimeFormat:
		return config.getStringFromEnv(RrhTimeFormat, config.TimeFormat)
	case RrhOnError:
		return config.getStringFromEnv(RrhOnError, config.OnError)
	default:
		return "", NotFound
	}
}

func (config *Config) getStringFromEnv(label string, valueFromConfigFile string) (value string, readFrom string) {
	if valueFromConfigFile != "" {
		return valueFromConfigFile, ConfigFile
	}
	var valueFromEnv = os.Getenv(label)
	if valueFromEnv != "" {
		return valueFromEnv, Env
	}
	return config.findDefaultValue(label)
}

func (config *Config) findDefaultValue(label string) (value string, readFrom string) {
	var home, _ = homedir.Dir()
	switch label {
	case RrhHome:
		return fmt.Sprintf("%s/.rrh", home), Default
	case RrhConfigPath:
		return fmt.Sprintf("%s/.rrh/config.json", home), Default
	case RrhDatabasePath:
		return fmt.Sprintf("%s/.rrh/database.json", home), Default
	case RrhDefaultGroupName:
		return "no-group", Default
	case RrhOnError:
		return Warn, Default
	case RrhTimeFormat:
		return Relative, Default
	case RrhAutoDeleteGroup:
		return "false", Default
	case RrhAutoCreateGroup:
		return "false", Default
	case RrhSortOnUpdating:
		return "false", Default
	default:
		return "", NotFound
	}
}

func configPath() string {
	var config = new(Config)
	var home, _ = config.getStringFromEnv(RrhConfigPath, "")
	return home
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
