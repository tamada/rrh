package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/tamada/rrh"
)

func assert(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Errorf("expected: %s, however actually %s", expected, actual)
	}
}

func TestConfigUnset(t *testing.T) {
	var testcases = []struct {
		args      []string
		wontError bool
	}{
		{[]string{rrh.AutoCreateGroup}, false},
		{[]string{"unknown"}, true},
		{[]string{rrh.AutoCreateGroup, rrh.AutoDeleteGroup}, false},
	}
	for _, tc := range testcases {
		var dbfile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			cuc := newUnsetCommand()
			cuc.SetArgs(tc.args)
			cuc.SetOut(os.Stdout)
			err := cuc.Execute()

			if err == nil && tc.wontError || err != nil && !tc.wontError {
				t.Errorf("%v: error wont %v, but got: %v", tc.args, tc.wontError, err)
			}
		})
		defer os.Remove(dbfile)
	}
}

func ExampleConfigCommand() {
	os.Setenv(rrh.ConfigPath, "../../../../testdata/config.json")
	os.Setenv(rrh.Home, "../../../../testdata/")
	os.Setenv(rrh.DatabasePath, "${RRH_HOME}/test_db.json")
	cmd := New()
	cmd.SetOut(os.Stdout)
	cmd.Execute()
	// Output:
	// RRH_ALIAS_PATH: ../../../../testdata/alias.json (default)
	// RRH_AUTO_CREATE_GROUP: true (config_file)
	// RRH_AUTO_DELETE_GROUP: false (config_file)
	// RRH_CLONE_DESTINATION: . (default)
	// RRH_COLOR: repository:fg=red+group:fg=magenta+label:op=bold+configValue:fg=green (default)
	// RRH_CONFIG_PATH: ../../../../testdata/config.json (environment)
	// RRH_DATABASE_PATH: ../../../../testdata/test_db.json (environment)
	// RRH_DEFAULT_GROUP_NAME: no-group (default)
	// RRH_ENABLE_COLORIZED: false (default)
	// RRH_HOME: ../../../../testdata/ (environment)
	// RRH_SORT_ON_UPDATING: true (config_file)
	// RRH_TIME_FORMAT: relative (default)
}

func ExampleConfigCommand_Run() {
	os.Setenv(rrh.ConfigPath, "../../../../testdata/config.json")
	os.Setenv(rrh.Home, "../../../../testdata/")
	os.Setenv(rrh.DatabasePath, "${RRH_HOME}/database.json")
	cmd := New()
	cmd.SetArgs([]string{"list"})
	cmd.SetOut(os.Stdout)
	cmd.Execute()
	// Output:
	// RRH_ALIAS_PATH: ../../../../testdata/alias.json (default)
	// RRH_AUTO_CREATE_GROUP: true (config_file)
	// RRH_AUTO_DELETE_GROUP: false (config_file)
	// RRH_CLONE_DESTINATION: . (default)
	// RRH_COLOR: repository:fg=red+group:fg=magenta+label:op=bold+configValue:fg=green (default)
	// RRH_CONFIG_PATH: ../../../../testdata/config.json (environment)
	// RRH_DATABASE_PATH: ../../../../testdata/database.json (environment)
	// RRH_DEFAULT_GROUP_NAME: no-group (default)
	// RRH_ENABLE_COLORIZED: false (default)
	// RRH_HOME: ../../../../testdata/ (environment)
	// RRH_SORT_ON_UPDATING: true (config_file)
	// RRH_TIME_FORMAT: relative (default)
}
func Example_listCommand_Run() {
	os.Setenv(rrh.ConfigPath, "../../../../testdata/config.json")
	os.Setenv(rrh.Home, "../../../../testdata/")
	os.Unsetenv(rrh.DatabasePath)

	cmd := newListCommand()
	cmd.SetOut(os.Stdout)
	cmd.Execute()
	// Output:
	// RRH_ALIAS_PATH: ../../../../testdata/alias.json (default)
	// RRH_AUTO_CREATE_GROUP: true (config_file)
	// RRH_AUTO_DELETE_GROUP: false (config_file)
	// RRH_CLONE_DESTINATION: . (default)
	// RRH_COLOR: repository:fg=red+group:fg=magenta+label:op=bold+configValue:fg=green (default)
	// RRH_CONFIG_PATH: ../../../../testdata/config.json (environment)
	// RRH_DATABASE_PATH: ../../../../testdata/database.json (default)
	// RRH_DEFAULT_GROUP_NAME: no-group (default)
	// RRH_ENABLE_COLORIZED: false (default)
	// RRH_HOME: ../../../../testdata/ (environment)
	// RRH_SORT_ON_UPDATING: true (config_file)
	// RRH_TIME_FORMAT: relative (default)
}

func TestLoadConfigFile(t *testing.T) {
	os.Setenv(rrh.ConfigPath, "../../../../testdata/config.json")

	var testdata = []struct {
		key   string
		value string
		from  rrh.ReadFrom
	}{
		{rrh.AutoDeleteGroup, "false", rrh.ConfigFile},
		{rrh.AutoCreateGroup, "true", rrh.ConfigFile},
		{rrh.SortOnUpdating, "true", rrh.ConfigFile},
		{rrh.ConfigPath, "../../../../testdata/config.json", rrh.Env},
		{rrh.TimeFormat, rrh.Relative, rrh.Default},
		{rrh.EnableColorized, "false", rrh.Default},
		{"unknown", "", rrh.NotFound},
	}

	var config = rrh.OpenConfig()
	for _, data := range testdata {
		if val, from := config.GetString(data.key); val != data.value || from != data.from {
			t.Errorf("%s: want: (%s, %s), got: (%s, %s)", data.key, data.value, data.from, val, from)
		}
	}
}

func TestOpenConfig(t *testing.T) {
	os.Unsetenv(rrh.Home)
	os.Unsetenv(rrh.DatabasePath)
	os.Unsetenv(rrh.ConfigPath)
	var home, _ = os.UserHomeDir()
	var testdata = []struct {
		key  string
		want string
	}{
		{rrh.Home, fmt.Sprintf("%s/.config/rrh", home)},
		{rrh.ConfigPath, fmt.Sprintf("%s/.config/rrh/config.json", home)},
		{rrh.DatabasePath, fmt.Sprintf("%s/.config/rrh/database.json", home)},
		{rrh.DefaultGroupName, "no-group"},
		{rrh.CloneDestination, "."},
		{rrh.AutoCreateGroup, "false"},
		{rrh.AutoDeleteGroup, "false"},
		{rrh.SortOnUpdating, "false"},
		{rrh.TimeFormat, rrh.Relative},
		{"unknown", ""},
	}
	// os.Unsetenv(RrhConfigPath)
	// os.Unsetenv(RrhHome)

	var config = rrh.OpenConfig()
	for _, data := range testdata {
		if value := config.GetDefaultValue(data.key); value != data.want {
			t.Errorf("%s: want: %s, got: %s", data.key, data.want, value)
		}
	}
	assert(t, config.GetDefaultValue("unknown"), "")
}

func TestConfigSet(t *testing.T) {
	var testdata = []struct {
		args      []string
		wontError bool
		value     string
		location  rrh.ReadFrom
	}{
		{[]string{"RRH_DEFAULT_GROUP_NAME", "newgroup"}, false, "newgroup", rrh.ConfigFile},
		{[]string{"RRH_DEFAULT_GROUP_NAME"}, true, "", ""},
		{[]string{"RRH_AUTO_DELETE_GROUP", "yes"}, true, "", ""},
		{[]string{rrh.ConfigPath, "../testdata/broken.json"}, true, "", ""},
	}
	for _, td := range testdata {
		var dbfile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			cmd := newSetCommand()
			cmd.SetArgs(td.args)
			cmd.SetOut(os.Stdout)
			err := cmd.Execute()

			if err == nil && td.wontError || err != nil && !td.wontError {
				t.Errorf("%v: wont error %v, but got %v", td.args, td.wontError, err)
			}
			if err == nil {
				var config = rrh.OpenConfig()
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
