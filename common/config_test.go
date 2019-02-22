package common

import (
	"fmt"
	"os"
	"testing"
)

func assert(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Errorf("expected: %s, however actually %s", expected, actual)
	}
}

func TestOpenConfigBrokenJson(t *testing.T) {
	os.Setenv(RrhConfigPath, "./testdata/broken.json")
	var config = OpenConfig()
	if config != nil {
		t.Error("broken json returns nil")
	}
}

func TestLoadConfigFile(t *testing.T) {
	os.Setenv(RrhConfigPath, "./testdata/config.json")
	var config = OpenConfig()
	if val, from := config.GetString(RrhAutoDeleteGroup); val != "false" || from != ConfigFile {
		t.Error("The RRH_AUTO_DELETE_GROUP was false in the config file!")
	}
	if val, from := config.GetString(RrhAutoCreateGroup); val != "true" || from != ConfigFile {
		t.Error("The RRH_AUTO_CREATE_GROUP was true in the config file!")
	}
	if val, from := config.GetString(RrhConfigPath); val != "./testdata/config.json" || from != Env {
		t.Error("The RRH_CONFIG_FILE was \"./testdata/config.json\" in environment value!")
	}
	if val, from := config.GetString(RrhTimeFormat); val != Relative || from != Default {
		t.Error("The RRH_TIME_FORMAT was \"Relative\" in environment value!")
	}
	if val, from := config.GetString(RrhOnError); val != Warn || from != Default {
		t.Error("The RRH_ON_ERRORwas WARN in default!")
	}
	if val, from := config.GetString("UnknownParameter"); val != "" || from != NotFound {
		t.Error("The UnknownParameter was \"\"!")
	}
}

func TestUpdateTrueFalseValue(t *testing.T) {
	os.Setenv(RrhConfigPath, "./testdata/testconfig.json")
	var config = OpenConfig()
	if err := config.Update(RrhAutoDeleteGroup, "True"); err != nil {
		t.Error(err.Error())
	}
	assert(t, config.GetValue(RrhAutoDeleteGroup), "true")
	if err := config.Update(RrhAutoDeleteGroup, "FALSE"); err != nil {
		t.Error(err.Error())
	}
	assert(t, config.GetValue(RrhAutoDeleteGroup), "false")
	if err := config.Update(RrhAutoDeleteGroup, "YES"); err == nil {
		t.Error("only accept true or false")
	}

	if err := config.Update(RrhAutoCreateGroup, "FALSE"); err != nil {
		t.Error(err.Error())
	}
	assert(t, config.GetValue(RrhAutoCreateGroup), "false")
	if err := config.Update(RrhAutoCreateGroup, "YES"); err == nil {
		t.Error("only accept true or false")
	}
}

func TestUpdateOnError(t *testing.T) {
	os.Setenv(RrhConfigPath, "./testdata/testconfig.json")
	var config = OpenConfig()
	if err := config.Update(RrhOnError, Ignore); err != nil {
		t.Error(err.Error())
	}
	if err := config.Update(RrhOnError, Fail); err != nil {
		t.Error(err.Error())
	}
	if err := config.Update(RrhOnError, FailImmediately); err != nil {
		t.Error(err.Error())
	}
	if err := config.Update(RrhOnError, Warn); err != nil {
		t.Error(err.Error())
	}
	if err := config.Update(RrhOnError, "unknown"); err == nil {
		t.Error("cannot set unknown to RrhOnError")
	}
}

func TestUpdateValue(t *testing.T) {
	os.Setenv(RrhConfigPath, "./testdata/testconfig.json")
	var config = OpenConfig()
	if err := config.Update(RrhConfigPath, "hogehoge"); err == nil {
		t.Error("RrhConfigPath cannot update")
	}
	if err := config.Update(RrhHome, "hoge1"); err != nil {
		t.Error(err.Error())
	}
	if err := config.Update(RrhDatabasePath, "hoge2"); err != nil {
		t.Error(err.Error())
	}
	if err := config.Update(RrhDefaultGroupName, "hoge3"); err != nil {
		t.Error(err.Error())
	}
	assert(t, config.GetValue(RrhHome), "hoge1")
	assert(t, config.GetValue(RrhDatabasePath), "hoge2")
	assert(t, config.GetValue(RrhDefaultGroupName), "hoge3")

	if err := config.Update("unknown", "hoge4"); err == nil {
		t.Error("unknown property was unknown")
	}
}

func TestOpenConfig(t *testing.T) {
	var config = OpenConfig()
	assert(t, config.GetDefaultValue(RrhHome), fmt.Sprintf("%s/.rrh", os.Getenv("HOME")))
	assert(t, config.GetDefaultValue(RrhConfigPath), fmt.Sprintf("%s/.rrh/config.json", os.Getenv("HOME")))
	assert(t, config.GetDefaultValue(RrhDatabasePath), fmt.Sprintf("%s/.rrh/database.json", os.Getenv("HOME")))
	assert(t, config.GetDefaultValue(RrhDefaultGroupName), "no-group")
	assert(t, config.GetDefaultValue(RrhOnError), Warn)
	assert(t, config.GetDefaultValue(RrhAutoDeleteGroup), "false")
	assert(t, config.GetDefaultValue(RrhAutoCreateGroup), "false")
	assert(t, config.GetDefaultValue("unknown"), "")
}

func TestFromatVariableAndValue(t *testing.T) {
	var config = OpenConfig()
	assert(t, config.formatVariableAndValue(RrhDefaultGroupName), "RRH_DEFAULT_GROUP_NAME: no-group (default)")
}
