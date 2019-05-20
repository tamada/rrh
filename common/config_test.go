package common

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func assert(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Errorf("expected: %s, however actually %s", expected, actual)
	}
}

func TestHelps(t *testing.T) {
	var command, _ = CommandFactory()
	if command.Help() != `rrh config <COMMAND> [ARGUMENTS]
COMMAND
    set <ENV_NAME> <VALUE>  set ENV_NAME to VALUE
    unset <ENV_NAME>        reset ENV_NAME
    list                    list all of ENVs (default)` {
		t.Errorf("help message did not match")
	}
	var clc, _ = listCommandFactory()
	if clc.Help() != `rrh config list` {
		t.Errorf("help message did not match")
	}
	var cuc, _ = unsetCommandFactory()
	if cuc.Help() != `rrh config unset <ENV_NAME...>
ARGUMENTS
    ENV_NAME   environment name.` {
		t.Errorf("help message did not match")
	}
	var csc, _ = setCommandFactory()
	if csc.Help() != `rrh config set <ENV_NAME> <VALUE>
ARGUMENTS
    ENV_NAME   environment name.
    VALUE      the value for the given environment.` {
		t.Errorf("help message did not match")
	}
}

func TestSynopsises(t *testing.T) {
	var command, _ = CommandFactory()
	if command.Synopsis() != "set/unset and list configuration of RRH." {
		t.Errorf("synopsis did not match")
	}

	var clc, _ = listCommandFactory()
	if clc.Synopsis() != "list the environment and its value." {
		t.Errorf("synopsis did not match")
	}

	var cuc, _ = unsetCommandFactory()
	if cuc.Synopsis() != "reset the given environment." {
		t.Errorf("synopsis did not match")
	}

	var csc, _ = setCommandFactory()
	if csc.Synopsis() != "set the environment with the given value." {
		t.Errorf("synopsis did not match")
	}
}

func TestConfigUnset(t *testing.T) {
	os.Setenv(RrhOnError, Fail)
	var testcases = []struct {
		args      []string
		status    int
		wontValue string
		wontFrom  ReadFrom
	}{
		{[]string{RrhAutoCreateGroup}, 0, "false", Default},
		{[]string{"unknown"}, 5, "", NotFound},
		{[]string{RrhAutoCreateGroup, "tooManyArgs"}, 1, "", ""},
	}
	for _, tc := range testcases {
		var dbfile = Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
			var cuc, _ = unsetCommandFactory()
			var statusCode = cuc.Run(tc.args)
			if statusCode != tc.status {
				t.Errorf("%v: status code did not match, wont: %d, got: %d", tc, tc.status, statusCode)
			}
			if statusCode == 0 {
				var config = OpenConfig()
				var value, from = config.GetString(tc.args[0])
				if value != tc.wontValue || from != tc.wontFrom {
					t.Errorf("%v: did not match: wont: (%s, %s), got: (%s, %s)", tc, tc.wontValue, tc.wontFrom, value, from)
				}
			}
		})
		defer os.Remove(dbfile)
	}
	os.Unsetenv(RrhOnError)
}

func ExampleCommand() {
	os.Setenv(RrhConfigPath, "../testdata/config.json")
	os.Setenv(RrhHome, "../testdata/")
	os.Setenv(RrhDatabasePath, "${RRH_HOME}/tmp.json")
	var command, _ = CommandFactory()
	command.Run([]string{}) // the output of no arguments are same as list subcommand.
	// Output:
	// RRH_AUTO_CREATE_GROUP: true (config_file)
	// RRH_AUTO_DELETE_GROUP: false (config_file)
	// RRH_CLONE_DESTINATION: . (default)
	// RRH_COLOR: repository:fg=red+group:fg=magenta+label:op=bold+configValue:fg=green (default)
	// RRH_CONFIG_PATH: ../testdata/config.json (environment)
	// RRH_DATABASE_PATH: ../testdata/tmp.json (environment)
	// RRH_DEFAULT_GROUP_NAME: no-group (default)
	// RRH_ENABLE_COLORIZED: false (default)
	// RRH_HOME: ../testdata/ (environment)
	// RRH_ON_ERROR: WARN (default)
	// RRH_SORT_ON_UPDATING: true (config_file)
	// RRH_TIME_FORMAT: relative (default)
}
func ExampleCommand_Run() {
	os.Setenv(RrhConfigPath, "../testdata/config.json")
	os.Setenv(RrhHome, "../testdata/")
	os.Setenv(RrhDatabasePath, "${RRH_HOME}/database.json")
	var command, _ = CommandFactory()
	command.Run([]string{"list"}) // the output of no arguments are same as list subcommand.
	// Output:
	// RRH_AUTO_CREATE_GROUP: true (config_file)
	// RRH_AUTO_DELETE_GROUP: false (config_file)
	// RRH_CLONE_DESTINATION: . (default)
	// RRH_COLOR: repository:fg=red+group:fg=magenta+label:op=bold+configValue:fg=green (default)
	// RRH_CONFIG_PATH: ../testdata/config.json (environment)
	// RRH_DATABASE_PATH: ../testdata/database.json (environment)
	// RRH_DEFAULT_GROUP_NAME: no-group (default)
	// RRH_ENABLE_COLORIZED: false (default)
	// RRH_HOME: ../testdata/ (environment)
	// RRH_ON_ERROR: WARN (default)
	// RRH_SORT_ON_UPDATING: true (config_file)
	// RRH_TIME_FORMAT: relative (default)
}
func Example_listCommand_Run() {
	os.Setenv(RrhConfigPath, "../testdata/config.json")
	os.Setenv(RrhHome, "../testdata/")
	os.Unsetenv(RrhDatabasePath)
	var clc, _ = listCommandFactory()
	clc.Run([]string{})
	// Output:
	// RRH_AUTO_CREATE_GROUP: true (config_file)
	// RRH_AUTO_DELETE_GROUP: false (config_file)
	// RRH_CLONE_DESTINATION: . (default)
	// RRH_COLOR: repository:fg=red+group:fg=magenta+label:op=bold+configValue:fg=green (default)
	// RRH_CONFIG_PATH: ../testdata/config.json (environment)
	// RRH_DATABASE_PATH: ../testdata/database.json (default)
	// RRH_DEFAULT_GROUP_NAME: no-group (default)
	// RRH_ENABLE_COLORIZED: false (default)
	// RRH_HOME: ../testdata/ (environment)
	// RRH_ON_ERROR: WARN (default)
	// RRH_SORT_ON_UPDATING: true (config_file)
	// RRH_TIME_FORMAT: relative (default)
}

func TestOpenConfigBrokenJson(t *testing.T) {
	os.Setenv(RrhConfigPath, "../testdata/broken.json")
	var config = OpenConfig()
	if config != nil {
		t.Error("broken json returns nil")
	}
}

func TestLoadConfigFile(t *testing.T) {
	os.Setenv(RrhConfigPath, "../testdata/config.json")

	var testdata = []struct {
		key   string
		value string
		from  ReadFrom
	}{
		{RrhAutoDeleteGroup, "false", ConfigFile},
		{RrhAutoCreateGroup, "true", ConfigFile},
		{RrhSortOnUpdating, "true", ConfigFile},
		{RrhConfigPath, "../testdata/config.json", Env},
		{RrhTimeFormat, Relative, Default},
		{RrhOnError, Warn, Default},
		{RrhEnableColorized, "false", Default},
		{"unknown", "", NotFound},
	}

	var config = OpenConfig()
	for _, data := range testdata {
		if val, from := config.GetString(data.key); val != data.value || from != data.from {
			t.Errorf("%s: want: (%s, %s), got: (%s, %s)", data.key, data.value, data.from, val, from)
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
		var dbfile = Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
			var config = OpenConfig()
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
		var dbfile = Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
			var config = OpenConfig()
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
		var dbfile = Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
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

func TestOpenConfig(t *testing.T) {
	os.Unsetenv(RrhHome)
	os.Unsetenv(RrhDatabasePath)
	os.Unsetenv(RrhConfigPath)
	var home, _ = homedir.Dir()
	var testdata = []struct {
		key  string
		want string
	}{
		{RrhHome, fmt.Sprintf("%s/.rrh", home)},
		{RrhConfigPath, fmt.Sprintf("%s/.rrh/config.json", home)},
		{RrhDatabasePath, fmt.Sprintf("%s/.rrh/database.json", home)},
		{RrhDefaultGroupName, "no-group"},
		{RrhCloneDestination, "."},
		{RrhOnError, Warn},
		{RrhAutoCreateGroup, "false"},
		{RrhAutoDeleteGroup, "false"},
		{RrhSortOnUpdating, "false"},
		{RrhTimeFormat, Relative},
		{"unknown", ""},
	}
	// os.Unsetenv(RrhConfigPath)
	// os.Unsetenv(RrhHome)

	var config = OpenConfig()
	for _, data := range testdata {
		if value := config.GetDefaultValue(data.key); value != data.want {
			t.Errorf("%s: want: %s, got: %s", data.key, data.want, value)
		}
	}
	assert(t, config.GetDefaultValue("unknown"), "")
}

func TestPrintErrors(t *testing.T) {
	var testcases = []struct {
		onError    string
		error      []error
		wontStatus int
		someOutput bool
	}{
		{Ignore, []error{}, 0, false},
		{Ignore, []error{errors.New("error")}, 0, false},
		{Warn, []error{}, 0, false},
		{Warn, []error{errors.New("error")}, 0, true},
		{Fail, []error{}, 0, false},
		{Fail, []error{errors.New("error")}, 5, true},
		{FailImmediately, []error{}, 0, false},
		{FailImmediately, []error{errors.New("error")}, 5, true},
	}

	var config = NewConfig()
	for _, tc := range testcases {
		config.Update(RrhOnError, tc.onError)
		var output = CaptureStdout(func() {
			var statusCode = config.PrintErrors(tc.error)
			if statusCode != tc.wontStatus {
				t.Errorf("%v: status code did not match, wont: %d, got: %d", tc, tc.wontStatus, statusCode)
			}
		})
		output = strings.TrimSpace(output)
		if (output == "") == tc.someOutput {
			t.Errorf("%v: output did not match, wont: %v, got: %v (%s)", tc, tc.someOutput, !tc.someOutput, output)
		}
	}
}

func TestConfigSet(t *testing.T) {
	var testdata = []struct {
		args       []string
		statusCode int
		value      string
		location   ReadFrom
	}{
		{[]string{"RRH_DEFAULT_GROUP_NAME", "newgroup"}, 0, "newgroup", ConfigFile},
		{[]string{"RRH_DEFAULT_GROUP_NAME"}, 1, "", ""},
		{[]string{"RRH_AUTO_DELETE_GROUP", "yes"}, 2, "", ""},
		{[]string{RrhConfigPath, "../testdata/broken.json"}, 2, "", ""},
	}
	for _, td := range testdata {
		var dbfile = Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
			var set, _ = setCommandFactory()
			var status = set.Run(td.args)
			if status != td.statusCode {
				t.Errorf("%v: status code did not match, wont: %d, got: %d", td.args, td.statusCode, status)
			}
			if status == 0 {
				var config = OpenConfig()
				var value, from = config.GetString(td.args[0])
				if value != td.value {
					t.Errorf("%v: set value did not match, wont: %s, got: %s", td.args, td.value, value)
				}
				if from != td.location {
					t.Errorf("%v: read from did not match, wont: %s, got: %s", td.args, td.location, from)
				}
			}
		})
		defer os.Remove(dbfile)
	}
}

func TestFormatVariableAndValue(t *testing.T) {
	os.Setenv(RrhConfigPath, "../testdata/config.json")
	var config = OpenConfig()
	assert(t, config.formatVariableAndValue(RrhDefaultGroupName), "RRH_DEFAULT_GROUP_NAME: no-group (default)")
	if config.IsSet(RrhOnError) {
		t.Errorf("IsSet accepts only bool variable")
	}
}
