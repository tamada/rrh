package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mitchellh/go-homedir"
	"github.com/tamada/rrh/lib"
)

func assert(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Errorf("expected: %s, however actually %s", expected, actual)
	}
}

func TestHelps(t *testing.T) {
	var command, _ = ConfigCommandFactory()
	if command.Help() != `rrh config <COMMAND> [ARGUMENTS]
COMMAND
    set <ENV_NAME> <VALUE>  set ENV_NAME to VALUE
    unset <ENV_NAME>        reset ENV_NAME
    list                    list all of ENVs (default)` {
		t.Errorf("help message did not match")
	}
	var clc, _ = configListCommandFactory()
	if clc.Help() != `rrh config list` {
		t.Errorf("help message did not match")
	}
	var cuc, _ = configUnsetCommandFactory()
	if cuc.Help() != `rrh config unset <ENV_NAME...>
ARGUMENTS
    ENV_NAME   environment name.` {
		t.Errorf("help message did not match")
	}
	var csc, _ = configSetCommandFactory()
	if csc.Help() != `rrh config set <ENV_NAME> <VALUE>
ARGUMENTS
    ENV_NAME   environment name.
    VALUE      the value for the given environment.` {
		t.Errorf("help message did not match")
	}
}

func TestSynopsises(t *testing.T) {
	var command, _ = ConfigCommandFactory()
	if command.Synopsis() != "set/unset and list configuration of RRH." {
		t.Errorf("synopsis did not match")
	}

	var clc, _ = configListCommandFactory()
	if clc.Synopsis() != "list the environment and its value." {
		t.Errorf("synopsis did not match")
	}

	var cuc, _ = configUnsetCommandFactory()
	if cuc.Synopsis() != "reset the given environment." {
		t.Errorf("synopsis did not match")
	}

	var csc, _ = configSetCommandFactory()
	if csc.Synopsis() != "set the environment with the given value." {
		t.Errorf("synopsis did not match")
	}
}

func TestConfigUnset(t *testing.T) {
	os.Setenv(lib.RrhOnError, lib.Fail)
	var testcases = []struct {
		args      []string
		status    int
		wontValue string
		wontFrom  lib.ReadFrom
	}{
		{[]string{lib.RrhAutoCreateGroup}, 0, "false", lib.Default},
		{[]string{"unknown"}, 5, "", lib.NotFound},
		{[]string{lib.RrhAutoCreateGroup, "tooManyArgs"}, 1, "", ""},
	}
	for _, tc := range testcases {
		var dbfile = lib.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *lib.Config, oldDB *lib.Database) {
			var cuc, _ = configUnsetCommandFactory()
			var statusCode = cuc.Run(tc.args)
			if statusCode != tc.status {
				t.Errorf("%v: status code did not match, wont: %d, got: %d", tc, tc.status, statusCode)
			}
			if statusCode == 0 {
				var config = lib.OpenConfig()
				var value, from = config.GetString(tc.args[0])
				if value != tc.wontValue || from != tc.wontFrom {
					t.Errorf("%v: did not match: wont: (%s, %s), got: (%s, %s)", tc, tc.wontValue, tc.wontFrom, value, from)
				}
			}
		})
		defer os.Remove(dbfile)
	}
	os.Unsetenv(lib.RrhOnError)
}

func ExampleConfigCommand() {
	os.Setenv(lib.RrhConfigPath, "../testdata/config.json")
	os.Setenv(lib.RrhHome, "../testdata/")
	os.Setenv(lib.RrhDatabasePath, "${RRH_HOME}/test_db.json")
	var command, _ = ConfigCommandFactory()
	command.Run([]string{}) // the output of no arguments are same as list subcommand.
	// Output:
	// RRH_AUTO_CREATE_GROUP: true (config_file)
	// RRH_AUTO_DELETE_GROUP: false (config_file)
	// RRH_CLONE_DESTINATION: . (default)
	// RRH_COLOR: repository:fg=red+group:fg=magenta+label:op=bold+configValue:fg=green (default)
	// RRH_CONFIG_PATH: ../testdata/config.json (environment)
	// RRH_DATABASE_PATH: ../testdata/test_db.json (environment)
	// RRH_DEFAULT_GROUP_NAME: no-group (default)
	// RRH_ENABLE_COLORIZED: false (default)
	// RRH_HOME: ../testdata/ (environment)
	// RRH_ON_ERROR: WARN (default)
	// RRH_SORT_ON_UPDATING: true (config_file)
	// RRH_TIME_FORMAT: relative (default)
}
func ExampleConfigCommand_Run() {
	os.Setenv(lib.RrhConfigPath, "../testdata/config.json")
	os.Setenv(lib.RrhHome, "../testdata/")
	os.Setenv(lib.RrhDatabasePath, "${RRH_HOME}/database.json")
	var command, _ = ConfigCommandFactory()
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
	os.Setenv(lib.RrhConfigPath, "../testdata/config.json")
	os.Setenv(lib.RrhHome, "../testdata/")
	os.Unsetenv(lib.RrhDatabasePath)
	var clc, _ = configListCommandFactory()
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

func TestLoadConfigFile(t *testing.T) {
	os.Setenv(lib.RrhConfigPath, "../testdata/config.json")

	var testdata = []struct {
		key   string
		value string
		from  lib.ReadFrom
	}{
		{lib.RrhAutoDeleteGroup, "false", lib.ConfigFile},
		{lib.RrhAutoCreateGroup, "true", lib.ConfigFile},
		{lib.RrhSortOnUpdating, "true", lib.ConfigFile},
		{lib.RrhConfigPath, "../testdata/config.json", lib.Env},
		{lib.RrhTimeFormat, lib.Relative, lib.Default},
		{lib.RrhOnError, lib.Warn, lib.Default},
		{lib.RrhEnableColorized, "false", lib.Default},
		{"unknown", "", lib.NotFound},
	}

	var config = lib.OpenConfig()
	for _, data := range testdata {
		if val, from := config.GetString(data.key); val != data.value || from != data.from {
			t.Errorf("%s: want: (%s, %s), got: (%s, %s)", data.key, data.value, data.from, val, from)
		}
	}
}

func TestOpenConfig(t *testing.T) {
	os.Unsetenv(lib.RrhHome)
	os.Unsetenv(lib.RrhDatabasePath)
	os.Unsetenv(lib.RrhConfigPath)
	var home, _ = homedir.Dir()
	var testdata = []struct {
		key  string
		want string
	}{
		{lib.RrhHome, fmt.Sprintf("%s/.rrh", home)},
		{lib.RrhConfigPath, fmt.Sprintf("%s/.rrh/config.json", home)},
		{lib.RrhDatabasePath, fmt.Sprintf("%s/.rrh/database.json", home)},
		{lib.RrhDefaultGroupName, "no-group"},
		{lib.RrhCloneDestination, "."},
		{lib.RrhOnError, lib.Warn},
		{lib.RrhAutoCreateGroup, "false"},
		{lib.RrhAutoDeleteGroup, "false"},
		{lib.RrhSortOnUpdating, "false"},
		{lib.RrhTimeFormat, lib.Relative},
		{"unknown", ""},
	}
	// os.Unsetenv(RrhConfigPath)
	// os.Unsetenv(RrhHome)

	var config = lib.OpenConfig()
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
		{lib.Ignore, []error{}, 0, false},
		{lib.Ignore, []error{errors.New("error")}, 0, false},
		{lib.Warn, []error{}, 0, false},
		{lib.Warn, []error{errors.New("error")}, 0, true},
		{lib.Fail, []error{}, 0, false},
		{lib.Fail, []error{errors.New("error")}, 5, true},
		{lib.FailImmediately, []error{}, 0, false},
		{lib.FailImmediately, []error{errors.New("error")}, 5, true},
	}

	var config = lib.NewConfig()
	for _, tc := range testcases {
		config.Update(lib.RrhOnError, tc.onError)
		var output = lib.CaptureStdout(func() {
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
		location   lib.ReadFrom
	}{
		{[]string{"RRH_DEFAULT_GROUP_NAME", "newgroup"}, 0, "newgroup", lib.ConfigFile},
		{[]string{"RRH_DEFAULT_GROUP_NAME"}, 1, "", ""},
		{[]string{"RRH_AUTO_DELETE_GROUP", "yes"}, 2, "", ""},
		{[]string{lib.RrhConfigPath, "../testdata/broken.json"}, 2, "", ""},
	}
	for _, td := range testdata {
		var dbfile = lib.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *lib.Config, oldDB *lib.Database) {
			var set, _ = configSetCommandFactory()
			var status = set.Run(td.args)
			if status != td.statusCode {
				t.Errorf("%v: status code did not match, wont: %d, got: %d", td.args, td.statusCode, status)
			}
			if status == 0 {
				var config = lib.OpenConfig()
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
	os.Setenv(lib.RrhConfigPath, "../testdata/config.json")
	var config = lib.OpenConfig()
	assert(t, formatVariableAndValue(config, lib.RrhDefaultGroupName), "RRH_DEFAULT_GROUP_NAME: no-group (default)")
	if config.IsSet(lib.RrhOnError) {
		t.Errorf("IsSet accepts only bool variable")
	}
}
