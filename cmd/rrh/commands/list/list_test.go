package list

import (
	"bytes"
	"os"
	"testing"

	"github.com/tamada/rrh"
)

func ExampleListCommand() {
	var dbFile = rrh.Rollback("../../../../testdata/database.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
		cmd := New()
		cmd.SetOut(os.Stdout)
		cmd.Execute()
	})
	defer os.Remove(dbFile)
	// Output:
	// no-group (2 repositories)
	//     rrh           ~/go/src/github.com/tamada/rrh
	//     helloworld    ../testdata/helloworld
	// 1 group, and 2 repositories
}

func ExampleListCommand_Run() {
	var dbFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
		cmd := New()
		cmd.SetArgs([]string{"-e", "group,count,id", "--entry", "note", "-e", "path"})
		cmd.SetOut(os.Stdout)
		cmd.Execute()
	})
	defer os.Remove(dbFile)
	// Output:
	// group1 (1 repository)
	//     Note: desc1
	//     repo1    path1
	// group2 (0 repositories)
	//     Note: desc2
	// group3 (1 repository) (abbreviate repositories)
	// 3 groups, and 2 repositories
}

func TestRunByCsvOutput(t *testing.T) {
	os.Setenv(rrh.DefaultGroupName, "group1")
	defer os.Unsetenv(rrh.DefaultGroupName)
	var dbFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
		buffer := bytes.NewBuffer([]byte{})
		cmd := New()
		cmd.SetArgs([]string{"--entry", "all", "--format", "csv"})
		cmd.SetOut(buffer)
		cmd.Execute()

		result := rrh.ReplaceNewline(buffer.String(), "&")
		var want = "group1,desc1,repo1,,path1,,&group3,desc3,repo2,,path2,origin,git@github.com:example/repo2.git"
		if result != want {
			t.Errorf("result did not match, wont: %s, got: %s", want, result)
		}
	})
	defer os.Remove(dbFile)
}

func TestFindResults(t *testing.T) {
	var testdata = []struct {
		targets []string
		want    []Result
	}{
		{[]string{"group1"}, []Result{{"group1", "desc1", false, []*Repo{{"repo1", "path1", "", []*rrh.Remote{}}}}}},
		{[]string{"group2"}, []Result{{"group2", "desc2", false, []*Repo{}}}},
	}

	for _, data := range testdata {
		var dbFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
			var results, err = FindResults(db, data.targets)
			if err != nil {
				t.Errorf("%v: group not found.", data.targets)
			}
			if results[0].GroupName != data.want[0].GroupName {
				t.Errorf("group name: want: %s, got:%s", data.want[0].GroupName, results[0].GroupName)
			}
			if results[0].Note != data.want[0].Note {
				t.Errorf("description: want: %s, got: %s", data.want[0].Note, results[0].Note)
			}
			if len(results[0].Repos) != len(data.want[0].Repos) {
				t.Errorf("# of repositories did not match: want: %d, got: %d", len(data.want[0].Repos), len(results[0].Repos))
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
