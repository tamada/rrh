package list

import (
	"fmt"
	"os"
	"testing"

	"github.com/tamada/rrh/common"
)

func open(jsonName string) *common.Database {
	os.Setenv(common.RrhDatabasePath, fmt.Sprintf("../testdata/%s", jsonName))
	var config = common.OpenConfig()
	var db, _ = common.Open(config)
	return db
}

func ExampleCommand() {
	os.Setenv(common.RrhDatabasePath, "../testdata/database.json")
	var list, _ = CommandFactory()
	list.Run([]string{})
	// Output:
	// no-group (1 repository)
	//     rrh          ~/go/src/github.com/tamada/rrh
	// 1 group, 1 repository
}

func ExampleCommand_Run() {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var list, _ = CommandFactory()
	list.Run([]string{"--desc", "--path"})
	// Output:
	// group1 (1 repository)
	//     Description  desc1
	//     repo1        path1
	// group2 (0 repositories)
	//     Description  desc2
	// group3 (1 repository)
	// 3 groups, 2 repositories
}

func TestRunByCsvOutput(t *testing.T) {
	os.Setenv(common.RrhDefaultGroupName, "group1")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var result = common.CaptureStdout(func() {
		var list, _ = CommandFactory()
		list.Run([]string{"--all-entries", "--csv"})
	})
	result = common.ReplaceNewline(result, "&")
	var want = "group1,desc1,repo1,path1&group3,desc3,repo2,path2"
	if result != want {
		t.Errorf("result did not match, wont: %s, got: %s", want, result)
	}
}

func TestSimpleResults(t *testing.T) {
	var testcases = []struct {
		args   []string
		status int
		result string
	}{
		{[]string{"--only-repositoryname"}, 0, "repo1,repo2"},
		{[]string{"--group-repository-form"}, 0, "group1/repo1,group3/repo2"},
	}
	for _, tc := range testcases {
		var result = common.CaptureStdout(func() {
			var list, _ = CommandFactory()
			var status = list.Run(tc.args)
			if status != tc.status {
				t.Errorf("%v: status code did not match: wont: %d, got: %d", tc.args, tc.status, status)
			}
		})
		result = common.ReplaceNewline(result, ",")
		if result != tc.result {
			t.Errorf("%v: result did not match: wont: %s, got: %s", tc.args, tc.result, result)
		}
	}
}

func TestFailedByUnknownOption(t *testing.T) {
	common.CaptureStdout(func() {
		var list, _ = CommandFactory()
		if val := list.Run([]string{"--unknown"}); val != 1 {
			t.Error("unknown option parsed!?")
		}
	})
}

func TestCommandHelpAndSynopsis(t *testing.T) {
	var list = Command{&options{}}
	var helpMessage = `rrh list [OPTIONS] [GROUPS...]
OPTIONS
    -d, --desc          print description of group.
    -p, --path          print local paths (default).
    -r, --remote        print remote urls.
    -A, --all-entries   print all entries of each repository.

    -a, --all           print all repositories, no omit repositories.
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
	var list = Command{&options{}}
	var testdata = []struct {
		targets []string
		want    []Result
	}{
		{[]string{"group1"}, []Result{{"group1", "desc1", false, []Repo{{"repo1", "path1", []common.Remote{}}}}}},
		{[]string{"group2"}, []Result{{"group2", "desc2", false, []Repo{}}}},
	}

	for _, data := range testdata {
		list.options.args = data.targets
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
