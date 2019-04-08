package common

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	var testcases = []struct {
		givenString string
		repo        string
		group       string
	}{
		{"repository:fg=white;op=bold,underscore", "37;1;4", ""},
		{"group: fg=red+repository:fg=white;op=bold,underscore", "37;1;4", "31"},
		{"group: fg=red+group: fg=blue", "", "34"},
	}

	for _, tc := range testcases {
		parse(tc.givenString)
		if repoColor != tc.repo {
			t.Errorf("%v: repo color did not match, wont: %s, got: %s", tc.givenString, tc.repo, repoColor)
		}
		if groupColor != tc.group {
			t.Errorf("%v: group color did not match, wont: %s, got: %s", tc.givenString, tc.group, groupColor)
		}
		fmt.Printf("repo: %s, group: %s\n", ColorrizedRepositoryID("repository"), ColorrizedGroupName("groupName"))
		ClearColorize()
	}
}
