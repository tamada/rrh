package common

import (
	"testing"

	c "github.com/gookit/color"
)

func TestParse(t *testing.T) {
	var testcases = []struct {
		givenString string
		repoColor   string
		groupColor  string
		labelColor  string
		wontRepo    string
		wontGroup   string
		wontLabel   string
	}{
		{"repository:fg=white;op=bold,underscore", "37;1;4", "", "", c.Style{c.FgWhite, c.Bold, c.OpUnderscore}.Render("repository"), "groupName", "label"},
		{"group: fg=red+repository:fg=white;op=bold,underscore", "37;1;4", "31", "", c.Style{c.FgWhite, c.Bold, c.OpUnderscore}.Render("repository"), c.FgRed.Render("groupName"), "label"},
		{"group: fg=red+group: fg=blue+label:op=bold", "", "34", "1", "repository", c.FgBlue.Render("groupName"), c.Bold.Render("label")},
	}

	for _, tc := range testcases {
		parse(tc.givenString)
		if repoColor != tc.repoColor {
			t.Errorf("%v: repo color did not match, wont: %s, got: %s", tc.givenString, tc.repoColor, repoColor)
		}
		if groupColor != tc.groupColor {
			t.Errorf("%v: group color did not match, wont: %s, got: %s", tc.givenString, tc.groupColor, groupColor)
		}
		if groupColor != tc.groupColor {
			t.Errorf("%v: group color did not match, wont: %s, got: %s", tc.givenString, tc.groupColor, groupColor)
		}
		if labelColor != tc.labelColor {
			t.Errorf("%v: label color did not match, wont: %s, got: %s", tc.givenString, tc.labelColor, labelColor)
		}
		if name := ColorizedRepositoryID("repository"); name != tc.wontRepo {
			t.Errorf("repository id did not match: wont: %s, got: %s", tc.wontRepo, name)
		}
		if name := ColorizedGroupName("groupName"); name != tc.wontGroup {
			t.Errorf("group name did not match: wont: %s, got: %s", tc.wontGroup, name)
		}
		if name := ColorizedLabel("label"); name != tc.wontLabel {
			t.Errorf("label did not match: wont: %s, got: %s", tc.wontLabel, name)
		}
		ClearColorize()
	}
}
