package path

import (
	"os"
	"strings"
	"testing"

	"github.com/tamada/rrh/common"
)

func TestSynopsis(t *testing.T) {
	var path, _ = PathCommandFactory()
	if path.Synopsis() != "print paths of specified repositories." {
		t.Error("Synopsis message is not matched.")
	}
}
func TestHelp(t *testing.T) {
	var path = PathCommand{}
	var message = `rrh path [OPTIONS] <REPOSITORIES...>
OPTIONS
    -m, --partial-match    treats the arguments as the patterns.
    -p, --show-only-path   show path only.
ARGUMENTS
    REPOSITORIES           repository ids.`
	if path.Help() != message {
		t.Error("Help message is not matched.")
	}
}

func TestPathCommand(t *testing.T) {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var testcases = []struct {
		args    []string
		status  int
		results string
	}{
		{[]string{}, 0, "repo1 path1,repo2 path2"},
		{[]string{"repo1"}, 0, "repo1 path1"},
		{[]string{"repo3"}, 0, ""},
		{[]string{"-p"}, 0, "path1,path2"},
		{[]string{"--partial-match", "2"}, 0, "repo2 path2"},
		{[]string{"--partial-match"}, 0, "repo1 path1,repo2 path2"},
		{[]string{"-p", "-m", "1"}, 0, "path1"},
		{[]string{"--unknown-option"}, 1, ""},
	}

	for _, tc := range testcases {
		var path, _ = PathCommandFactory()
		var output, _ = common.CaptureStdout(func() {
			var status = path.Run(tc.args)
			if status != tc.status {
				t.Errorf("%v: status code did not match: wont: %d, got: %d", tc.args, tc.status, status)
			}
		})
		if tc.status == 0 {
			output = strings.TrimSpace(output)
			output = strings.Replace(output, "\n", ",", -1)
			if output != tc.results {
				t.Errorf("%v: output did not match: wont: %s, got: %s", tc.args, tc.results, output)
			}
		}
	}
}