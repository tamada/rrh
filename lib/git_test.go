package lib

import (
	"testing"
)

func TestNewStatusOptions(t *testing.T) {
	var opts = NewStatusOption()
	if opts.BranchStatus != false && opts.RemoteStatus != false {
		t.Errorf("default values of StatusOption did not match.")
	}
}

func TestGenerateMessage(t *testing.T) {
	var testcases = []struct {
		stagingFlag bool
		changesFlag bool
		wontMessage string
	}{
		{true, true, "Changes in staging"},
		{false, true, "Changes in workspace"},
		{true, false, "No changes"},
		{false, false, "No changes"},
	}
	for _, tc := range testcases {
		var message = generateMessage(tc.stagingFlag, tc.changesFlag)
		if message != tc.wontMessage {
			t.Errorf("generateMessage(%v, %v) did not match, wont: %s, got: %s", tc.stagingFlag, tc.changesFlag, tc.wontMessage, message)
		}
	}
}

func TestFindRemotes(t *testing.T) {
	var testdata = []struct {
		path      string
		errorFlag bool
		count     int
	}{
		{"../testdata/dummygit", true, 0},
		{"../testdata/helloworld", false, 1},
		{"../testdata/helloworld_noremote", false, 0},
	}
	for _, td := range testdata {
		var remotes, err = FindRemotes(td.path)
		if (err == nil) == td.errorFlag {
			t.Errorf("%s: error flag did not match, wont: %v, got: %v, %v", td.path, td.errorFlag, !td.errorFlag, err)
		}
		if err != nil && td.count != len(remotes) {
			t.Errorf("%s: remote count did not match, wont: %d, got: %d", td.path, td.count, len(remotes))
		}
	}
}
