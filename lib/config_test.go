package lib

import (
	"fmt"
	"os"
	"testing"
)

func TestValidateArgumentsOnUpdate(t *testing.T) {
	var testcases = []struct {
		givesLabel  string
		givesValue  string
		wontError   bool
		wontMessage string
	}{
		{RrhHome, "~/.rrh", false, ""},
		{RrhConfigPath, "~/.rrh_config", true, RrhConfigPath + ": cannot set in config file"},
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
		{RrhAutoDeleteGroup, "True", false, "true"},
		{RrhAutoDeleteGroup, "FALSE", false, "false"},
		{RrhAutoDeleteGroup, "FALSE", false, "false"},
		{RrhAutoDeleteGroup, "YES", true, ""},
		{RrhAutoCreateGroup, "FALSE", false, "false"},
		{RrhAutoCreateGroup, "YES", true, ""},
		{RrhSortOnUpdating, "FALSE", false, "false"},
		{RrhSortOnUpdating, "YES", true, ""},
	}

	for _, data := range testdata {
		var dbfile = Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *Config, oldDB *Database) {
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

func TestUpdateOnError(t *testing.T) {
	var testdata = []struct {
		key     string
		success bool
	}{
		{Ignore, true},
		{Fail, true},
		{FailImmediately, true},
		{Warn, true},
		{"unknown", false},
	}

	for _, data := range testdata {
		var dbfile = Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *Config, oldDB *Database) {
			if err := config.Update(RrhOnError, data.key); (err == nil) != data.success {
				t.Errorf("%s: set to \"%s\", success: %v", RrhOnError, data.key, data.success)
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
		{RrhConfigPath, "hogehoge", true, ""},
		{RrhHome, "hoge1", false, "hoge1"},
		{RrhDatabasePath, "hoge2", false, "hoge2"},
		{RrhDefaultGroupName, "hoge3", false, "hoge3"},
		{RrhTimeFormat, "not-relative-string", false, "not-relative-string"},
		{"unknown", "hoge4", true, ""},
	}
	for _, td := range testdata {
		var dbfile = Rollback("../testdata/test_db.json", "../testdata/config.json", func(unusedConfig *Config, oldDB *Database) {
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
