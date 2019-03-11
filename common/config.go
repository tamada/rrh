package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const VERSION = "0.1"

const (
	RrhHome             = "RRH_HOME"
	RrhConfigPath       = "RRH_CONFIG_PATH"
	RrhDatabasePath     = "RRH_DATABASE_PATH"
	RrhDefaultGroupName = "RRH_DEFAULT_GROUP_NAME"
	RrhOnError          = "RRH_ON_ERROR"
	RrhAutoDeleteGroup  = "RRH_AUTO_DELETE_GROUP"
	RrhAutoCreateGroup  = "RRH_AUTO_CREATE_GROUP"
	RrhTimeFormat       = "RRH_TIME_FORMAT"
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

type Config struct {
	Home             string `json:"rrh_home"`
	AutoDeleteGroup  string `json:"rrh_auto_delete_group"`
	AutoCreateGroup  string `json:"rrh_auto_create_group"`
	ConfigPath       string `json:"rrh_config_path"`
	DatabasePath     string `json:"rrh_database_path"`
	DefaultGroupName string `json:"rrh_default_group_name"`
	TimeFormat       string `json:"rrh_time_format"`
	OnError          string `json:"rrh_on_error"`
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
		return fmt.Errorf("%s: does not set on config file. set on environment.", label)
	}
	return fmt.Errorf("%s: Unknown variable name", label)
}

func (config *Config) GetValue(label string) string {
	var value, _ = config.GetString(label)
	return value
}

func (config *Config) GetDefaultValue(label string) string {
	var value, _ = config.findDefaultValue(label)
	return value
}

func (config *Config) GetString(label string) (value string, readFrom string) {
	switch label {
	case RrhAutoDeleteGroup:
		return config.getStringFromEnv(RrhAutoDeleteGroup, config.AutoDeleteGroup)
	case RrhAutoCreateGroup:
		return config.getStringFromEnv(RrhAutoCreateGroup, config.AutoCreateGroup)
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
	switch label {
	case RrhHome:
		return fmt.Sprintf("%s/.rrh", os.Getenv("HOME")), Default
	case RrhConfigPath:
		return fmt.Sprintf("%s/.rrh/config.json", os.Getenv("HOME")), Default
	case RrhDatabasePath:
		return fmt.Sprintf("%s/.rrh/database.json", os.Getenv("HOME")), Default
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
	default:
		return "", NotFound
	}
}

func configPath() string {
	var config = new(Config)
	var home, _ = config.getStringFromEnv(RrhConfigPath, "")
	return home
}

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
