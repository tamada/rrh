package lib

import (
	"os"
	"testing"

	c "github.com/gookit/color"
)

func TestEnableColorize(t *testing.T) {
	os.Setenv(RrhConfigPath, "../testdata/config.json")
	var config = OpenConfig()
	config.Update(RrhEnableColorized, "true")
	var cs = InitializeColor(config)

	var groupName1 = cs.ColorizedGroupName("group")
	var groupNameWont = c.FgMagenta.Render("group")
	if groupName1 != groupNameWont {
		t.Errorf("group name did not match: wont: %s, got: %s", groupName1, groupNameWont)
	}

	cs.SetColorize(false)
	var groupName2 = cs.ColorizedGroupName("group")
	if groupName2 != "group" {
		t.Errorf("group name did not match: wont: %s, got: %s", groupName2, "group")
	}

	cs.SetColorize(true)
	var groupName3 = cs.ColorizedGroupName("group")
	if groupName3 != groupNameWont {
		t.Errorf("groupName did not match: wont: %s, got: %s", groupName3, groupNameWont)
	}
	cs.ClearColorize()
}

func TestParse(t *testing.T) {
	var testcases = []struct {
		givenString    string
		repoColor      string
		groupColor     string
		labelColor     string
		boolTrueColor  string
		boolFalseColor string
		wontRepo       string
		wontGroup      string
		wontLabel      string
		wontBoolTrue   string
		wontBoolFalse  string
	}{
		{"repository:fg=white;op=bold,underscore", "37;1;4", "", "", "", "", c.Style{c.FgWhite, c.Bold, c.OpUnderscore}.Render("repository"), "groupName", "label", "true", "false"},
		{"group: fg=red+repository:fg=white;op=bold,underscore", "37;1;4", "31", "", "", "", c.Style{c.FgWhite, c.Bold, c.OpUnderscore}.Render("repository"), c.FgRed.Render("groupName"), "label", "true", "false"},
		{"group: fg=red+group: fg=blue+label:op=bold", "", "34", "1", "", "", "repository", c.FgBlue.Render("groupName"), c.Bold.Render("label"), "true", "false"},
		{"boolTrue: fg=green+boolFalse: fg=blue", "", "", "", "32", "34", "repository", "groupName", "label", c.FgGreen.Render("true"), c.FgBlue.Render("false")},
	}

	for _, tc := range testcases {
		var cs = Color{settings: colorSettings{}, funcs: colorFuncs{}}
		cs.parse(tc.givenString)
		if v, ok := cs.settings["repository"]; !ok || v != tc.repoColor {
			t.Errorf("%v: repo color did not match, wont: %s, got: %s", tc.givenString, tc.repoColor, v)
		}
		if v, ok := cs.settings["group"]; !ok || v != tc.groupColor {
			t.Errorf("%v: group color did not match, wont: %s, got: %s", tc.givenString, tc.groupColor, v)
		}
		if v, ok := cs.settings["label"]; !ok || v != tc.labelColor {
			t.Errorf("%v: label color did not match, wont: %s, got: %s", tc.givenString, tc.labelColor, v)
		}
		if name := cs.ColorizedRepositoryID("repository"); name != tc.wontRepo {
			t.Errorf("repository id did not match: wont: %s, got: %s", tc.wontRepo, name)
		}
		if name := cs.ColorizedGroupName("groupName"); name != tc.wontGroup {
			t.Errorf("group name did not match: wont: %s, got: %s", tc.wontGroup, name)
		}
		if name := cs.ColorizedLabel("label"); name != tc.wontLabel {
			t.Errorf("label did not match: wont: %s, got: %s", tc.wontLabel, name)
		}
	}
}
