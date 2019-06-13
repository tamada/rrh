package internal

import (
	"os"
	"testing"

	"github.com/tamada/rrh/lib"
)

func ExampleListCommand() {
	var dbFile = lib.Rollback("../testdata/database.json", "../testdata/config.json", func(config *lib.Config, oldDB *lib.Database) {
		var list, _ = ListCommandFactory()
		list.Run([]string{})
	})
	defer os.Remove(dbFile)
	// Output:
	// no-group (1 repository)
	//     rrh          ~/go/src/github.com/tamada/rrh
	// 1 group, 1 repository
}

func ExampleListCommand_Run() {
	var dbFile = lib.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *lib.Config, oldDB *lib.Database) {
		var list, _ = ListCommandFactory()
		list.Run([]string{"--desc", "--path"})
	})
	defer os.Remove(dbFile)
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
	os.Setenv(lib.RrhDefaultGroupName, "group1")
	defer os.Unsetenv(lib.RrhDefaultGroupName)
	var dbFile = lib.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *lib.Config, oldDB *lib.Database) {
		var result = lib.CaptureStdout(func() {
			var list, _ = ListCommandFactory()
			list.Run([]string{"--all-entries", "--csv"})
		})
		result = lib.ReplaceNewline(result, "&")
		var want = "group1,desc1,repo1,path1&group3,desc3,repo2,path2,origin,git@github.com:example/repo2.git"
		if result != want {
			t.Errorf("result did not match, wont: %s, got: %s", want, result)
		}
	})
	defer os.Remove(dbFile)
}

func TestSimpleResults(t *testing.T) {
	var testcases = []struct {
		args   []string
		status int
		result string
	}{
		{[]string{"--only-repositoryname"}, 0, "repo1,repo2"},
		{[]string{"--group-repository-form"}, 0, "group1/repo1,group3/repo2"},
		{[]string{"not-included-group"}, 3, ""},
	}
	for _, tc := range testcases {
		var dbFile = lib.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *lib.Config, oldDB *lib.Database) {
			var result = lib.CaptureStdout(func() {
				var list, _ = ListCommandFactory()
				var status = list.Run(tc.args)
				if status != tc.status {
					t.Errorf("%v: status code did not match: wont: %d, got: %d", tc.args, tc.status, status)
				}
			})
			result = lib.ReplaceNewline(result, ",")
			if tc.status == 0 && result != tc.result {
				t.Errorf("%v: result did not match: wont: %s, got: %s", tc.args, tc.result, result)
			}
		})
		defer os.Remove(dbFile)
	}
}

func TestFailedByUnknownOption(t *testing.T) {
	lib.CaptureStdout(func() {
		var list, _ = ListCommandFactory()
		if val := list.Run([]string{"--unknown"}); val != 1 {
			t.Error("unknown option parsed!?")
		}
	})
}

func TestCommandHelpAndSynopsisOfListCommand(t *testing.T) {
	var list = ListCommand{&listOptions{}}
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
	var list = ListCommand{&listOptions{}}
	var testdata = []struct {
		targets []string
		want    []Result
	}{
		{[]string{"group1"}, []Result{{"group1", "desc1", false, []Repo{{"repo1", "path1", []lib.Remote{}}}}}},
		{[]string{"group2"}, []Result{{"group2", "desc2", false, []Repo{}}}},
	}

	for _, data := range testdata {
		var dbFile = lib.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *lib.Config, db *lib.Database) {
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
		})
		defer os.Remove(dbFile)
	}
}
