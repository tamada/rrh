package rrh

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateArgumentsOnUpdate(t *testing.T) {
	var testcases = []struct {
		givesLabel  string
		givesValue  string
		wontError   bool
		wontMessage string
	}{
		{Home, "~/.rrh", false, ""},
		{ConfigPath, "~/.rrh_config", true, ConfigPath + ": cannot set in config file"},
		{"UnknownVariableName", "hoge", true, "UnknownVariableName: unknown variable name"},
	}

	for _, tc := range testcases {
		var err = validateArgumentsOnUpdate(tc.givesLabel, tc.givesValue)
		if (err == nil) == tc.wontError {
			t.Errorf("%s, %s: error flag did not match, wont %v, got %v", tc.givesLabel, tc.givesValue, tc.wontError, !tc.wontError)
		}
		if err != nil && err.Error() != tc.wontMessage {
			t.Errorf("%s, %s: error message did not match, wont %s, got %s", tc.givesLabel, tc.givesValue, tc.wontMessage, err.Error())
		}
	}
}

func TestNormalizeValueOfOnError(t *testing.T) {
	var testcases = []struct {
		givesValue  string
		wontValue   string
		wontError   bool
		wontMessage string
	}{
		{Fail, Fail, false, ""},
		{"warn", Warn, false, ""},
		{"log", "", true, fmt.Sprintf("log: Unknown value of RRH_ON_ERROR (must be %s, %s, %s, or %s)", Fail, FailImmediately, Warn, Ignore)},
	}

	for _, tc := range testcases {
		var value, err = normalizeValueOfOnError(tc.givesValue)
		if (err == nil) == tc.wontError {
			t.Errorf("%s: error flag did not match, wont %v, got %v", tc.givesValue, tc.wontError, !tc.wontError)
		}
		if err == nil && value != tc.wontValue {
			t.Errorf("%s: resultant value did not match, wont %s, got %s", tc.givesValue, tc.wontValue, value)
		}
		if err != nil && err.Error() != tc.wontMessage {
			t.Errorf("%s: error message did not match, wont %s, got %s", tc.givesValue, tc.wontMessage, err.Error())
		}
	}
}

func TestTrueOrFalse(t *testing.T) {
	var testcases = []struct {
		givesValue  string
		wontValue   string
		wontError   bool
		wontMessage string
	}{
		{"TRUE", trueString, false, ""},
		{"false", falseString, false, ""},
		{"yes", "", true, "yes: not true nor false"},
	}

	for _, tc := range testcases {
		var value, err = trueOrFalse(tc.givesValue)
		if (err == nil) == tc.wontError {
			t.Errorf("%s: error flag did not match, wont %v, got %v", tc.givesValue, tc.wontError, !tc.wontError)
		}
		if err == nil && value != tc.wontValue {
			t.Errorf("%s: resultant value did not match, wont %s, got %s", tc.givesValue, tc.wontValue, value)
		}
		if err != nil && err.Error() != tc.wontMessage {
			t.Errorf("%s: error message did not match, wont %s, got %s", tc.givesValue, tc.wontMessage, err.Error())
		}
	}
}
func TestUpdateTrueFalseValue(t *testing.T) {
	var testdata = []struct {
		key       string
		value     string
		wantError bool
		wantValue string
	}{
		{AutoDeleteGroup, "True", false, "true"},
		{AutoDeleteGroup, "FALSE", false, "false"},
		{AutoDeleteGroup, "FALSE", false, "false"},
		{AutoDeleteGroup, "YES", true, ""},
		{AutoCreateGroup, "FALSE", false, "false"},
		{AutoCreateGroup, "YES", true, ""},
		{SortOnUpdating, "FALSE", false, "false"},
		{SortOnUpdating, "YES", true, ""},
	}

	for _, data := range testdata {
		var dbfile = Rollback("testdata/test_db.json", "testdata/config.json", func(config *Config, oldDB *Database) {
			if err := config.Update(data.key, data.value); (err == nil) == data.wantError {
				t.Errorf("%s: set to \"%s\", error: %s", data.key, data.value, err.Error())
			}
			if val := config.GetValue(data.key); !data.wantError && val != data.wantValue {
				t.Errorf("%s: want: %s, got: %s", data.key, data.wantValue, val)
			}
		})
		defer os.Remove(dbfile)
	}
}

func TestUpdateValue(t *testing.T) {
	var testdata = []struct {
		label       string
		value       string
		shouldError bool
		wontValue   string
	}{
		{ConfigPath, "hogehoge", true, ""},
		{Home, "hoge1", false, "hoge1"},
		{DatabasePath, "hoge2", false, "hoge2"},
		{DefaultGroupName, "hoge3", false, "hoge3"},
		{TimeFormat, "not-relative-string", false, "not-relative-string"},
		{"unknown", "hoge4", true, ""},
	}
	for _, td := range testdata {
		var dbfile = Rollback("testdata/test_db.json", "testdata/config.json", func(unusedConfig *Config, oldDB *Database) {
			var config = NewConfig()
			var err = config.Update(td.label, td.value)
			if (err == nil) == td.shouldError {
				t.Errorf("error of Update(%s, %s) did not match, wont: %v, got: %v", td.label, td.value, td.shouldError, !td.shouldError)
			}
			if err == nil {
				var value = config.GetValue(td.label)
				if value != td.wontValue {
					t.Errorf("Value after Update(%s, %s) did not match, wont: %v, got: %v", td.label, td.value, td.wontValue, value)
				}
			}
		})
		defer os.Remove(dbfile)
	}
}

func TestConfigIsSet(t *testing.T) {
	var dbFile = Rollback("testdata/test_db.json", "testdata/config.json", func(config *Config, db *Database) {
		if config.IsSet(ConfigPath) {
			t.Errorf("not boolean variable is specified")
		}
		var home, _ = os.UserHomeDir()
		if config.GetDefaultValue(ConfigPath) != filepath.Join(home, ".config/rrh/config.json") {
			t.Errorf("RrhConfigPath did not match")
		}
		var _, from1 = config.findDefaultValue("UnknownVariable")
		if from1 != NotFound {
			t.Errorf("Unknown variable can get")
		}
		var _, from2 = config.GetString("UnknownVariable")
		if from2 != NotFound {
			t.Errorf("Unknown variable can get")
		}
		var err = config.Unset("UnknownVariable")
		if err == nil {
			t.Errorf("Unknown variable can Unset")
		}
		var beforeFlag = config.IsSet(AutoCreateGroup)
		config.Unset(AutoCreateGroup)
		var afterFlag = config.IsSet(AutoCreateGroup)
		if afterFlag || !beforeFlag {
			t.Errorf("beforeFlag should be true, and afterFlag should be false after Unset of RrhAutoCreateGroup")
		}
		config.StoreConfig()
		var config2 = OpenConfig()
		var afterFlag2 = config2.IsSet(AutoCreateGroup)
		if afterFlag2 {
			t.Errorf("afterFlag2 should be false because unset and store the config")
		}
	})
	defer os.Remove(dbFile)
}

func convertToErrors(messages []string) []error {
	var errs = []error{}
	for _, msg := range messages {
		errs = append(errs, errors.New(msg))
	}
	return errs
}

func TestOpenConfigBrokenJson(t *testing.T) {
	os.Setenv(ConfigPath, "testdata/broken.json")
	var config = OpenConfig()
	if config != nil {
		t.Error("broken json returns nil")
	}
}
