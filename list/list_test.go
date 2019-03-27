package list

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/tamada/rrh/common"
)

func open(jsonName string) *common.Database {
	os.Setenv(common.RrhDatabasePath, fmt.Sprintf("../testdata/%s", jsonName))
	var config = common.OpenConfig()
	var db, _ = common.Open(config)
	return db
}

func ExampleListCommand_Run() {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var list, _ = ListCommandFactory()
	list.Run([]string{"--desc", "--path"})
	// Output:
	// group1
	//     Description  desc1
	//     repo1        path1
	// group2
	//     Description  desc2
	// group3 (1 repository)
}

func TestRunByCsvOutput(t *testing.T) {
	os.Setenv(common.RrhDefaultGroupName, "group1")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var result, _ = common.CaptureStdout(func() {
		var list, _ = ListCommandFactory()
		list.Run([]string{"--all-entries", "--csv"})
	})
	result = strings.TrimSpace(result)
	var want = "group1,desc1,repo1,path1\ngroup3,desc3,repo2,path2"
	if result != want {
		t.Errorf("result did not match\ngot: %s\nwont: %s", result, want)
	}
}

func TestFailedByUnknownOption(t *testing.T) {
	common.CaptureStdout(func() {
		var list, _ = ListCommandFactory()
		if val := list.Run([]string{"--unknown"}); val != 1 {
			t.Error("unknown option parsed!?")
		}
	})
}

func TestListCommandHelpAndSynopsis(t *testing.T) {
	var list = ListCommand{&listOptions{}}
	var helpMessage = `rrh list [OPTIONS] [GROUPS...]
OPTIONS
    -a, --all           print all repositories, no omit repositories.
    -d, --desc          print description of group.
    -p, --path          print local paths (default).
    -r, --remote        print remote urls.
                        if any options of above are specified, '-a' are specified.
    -A, --all-entries   print all entries of each repository.

    -c, --csv           print result as csv format.
ARGUMENTS
    GROUPS    print managed repositories categorized in the groups.
              if no groups are specified, all groups are printed.`

	if list.Help() != helpMessage {
		t.Error("help message did not match")
	}
	if list.Synopsis() != "print managed repositories and their groups." {
		t.Error("Synopsis did not match")
	}
}

func TestFindResults(t *testing.T) {
	var db = open("tmp.json")
	var list = ListCommand{&listOptions{}}
	var testdata = []struct {
		targets []string
		want    []ListResult
	}{
		{[]string{"group1"}, []ListResult{{"group1", "desc1", false, []Repo{{"repo1", "path1", []common.Remote{}}}}}},
		{[]string{"group2"}, []ListResult{{"group2", "desc2", false, []Repo{}}}},
	}

	for _, data := range testdata {
		list.Options.args = data.targets
		var results, err = list.FindResults(db)
		if err != nil {
			t.Errorf("%v: group not found.", data.targets)
		}
		if results[0].GroupName != data.want[0].GroupName {
			t.Errorf("group name: want: %s, got:%s", data.want[0].GroupName, results[0].GroupName)
		}
		if results[0].Description != data.want[0].Description {
			t.Errorf("description: want: %s, got: %s", data.want[0].Description, results[0].Description)
		}
		if len(results[0].Repos) != len(data.want[0].Repos) {
			t.Errorf("# of repositories did not match: want: %d, got: %d\n", len(data.want[0].Repos), len(results[0].Repos))
		}
		if len(results[0].Repos) > 0 {
			if results[0].Repos[0].Name != data.want[0].Repos[0].Name {
				t.Errorf("repo name: want: %s, got:%s", data.want[0].Repos[0].Name, results[0].Repos[0].Name)
			}
			if results[0].Repos[0].Path != data.want[0].Repos[0].Path {
				t.Errorf("repo path: want: %s, got:%s", data.want[0].Repos[0].Path, results[0].Repos[0].Path)
			}
		}
	}
}
