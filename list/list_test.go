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

func TestFindResults(t *testing.T) {
	var db = open("tmp.json")
	var list = ListCommand{}
	var testdata = []struct {
		target          string
		wantGroupName   string
		wantDescription string
		wantRepos       []struct {
			wantRepoName     string
			wantPath         string
			wantRemoteLength int
		}
	}{
		{"group1", "group1", "", {"repo1", "path1", 0}},
		{"group2", "group2", "", {}},
	}

	for _, data := range testdata {
		list.Options.args = []string{data.target}
		var results, err = list.FindResults(db)
		if err != nil {
			t.Errorf("%s: group not found.", data.target)
		}
		if results[0].GroupName != data.wantGroupName {
			t.Errorf("group name: want: %s, got:%s", data.wantGroupName, results[0].GroupName)
		}
		if results[0].Description != data.wantDescription {
			t.Errorf("description: want: %s, got: %s", data.wantDescription, results[0].Description)
		}
		if results[0].Repos[0].Name != data.wantRepos[0].wantRepoName {
			t.Errorf("repo name: want: %s, got:%s", data.wantRepos[0].wantRepoName, results[0].Repos[0].Name)
		}
		if results[0].Repos[0].Path != data.wantRepos[0].wantPath {
			t.Errorf("repo path: want: %s, got:%s", data.wantRepos[0].wantPath, results[0].Repos[0].Path)
		}
	}
}
