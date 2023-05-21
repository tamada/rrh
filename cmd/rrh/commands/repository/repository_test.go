package repository

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tamada/rrh"
)

func TestRepository(t *testing.T) {
	var testcases = []struct {
		args         []string
		hasError     bool
		output       string
		ignoreOutput bool
	}{
		{[]string{"unknown-command"}, true, "", true},
		{[]string{"list"}, false, "group1/repo1    path1    +group3/repo2    path2", false},
		{[]string{"list", "--entry", "id"}, false, "repo1    +repo2", false},
		{[]string{"list", "repo2"}, true, "", true},
		{[]string{"list", "-e", "id,path", "group3"}, false, "repo2    path2", false},
		{[]string{"list", "--entry", "group,id", "group1"}, false, "group1/repo1", false},
	}
	for _, tc := range testcases {
		var dbFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			var output = rrh.CaptureStdout(func() {
				cmd := New()
				cmd.SetArgs(tc.args)
				cmd.SetOutput(os.Stdout)
				err := cmd.Execute()
				if tc.hasError && err == nil || !tc.hasError && err != nil {
					t.Errorf("%v: status code did not match, wont: %v, got: %v", tc.args, tc.hasError, err)
				}
			})
			if !tc.ignoreOutput {
				output = strings.TrimSpace(output)
				output = rrh.ReplaceNewline(output, "+")
				if output != tc.output {
					t.Errorf("%v: output did not match, wont: %s, got: %s", tc.args, tc.output, output)
				}
			}
		})
		defer os.Remove(dbFile)
	}
}

func TestListRepository(t *testing.T) {
	var testcases = []struct {
		args     []string
		hasError bool
		output   string
	}{
		{[]string{"--entry", "id"}, false, "repo1    +repo2"},
		{[]string{"--entry", "path"}, false, "path1    +path2"},
		{[]string{"--entry", "id,group"}, false, "group1/repo1    +group3/repo2"},
		{[]string{"--entry", "id", "group3"}, false, "repo2"},
		{[]string{"--entry", "path", "group1"}, false, "path1"},
		{[]string{"--entry", "group,id", "group3"}, false, "group3/repo2"},
		{[]string{}, false, "group1/repo1    path1    +group3/repo2    path2"},
	}
	for _, tc := range testcases {
		var dbFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			var output = rrh.CaptureStdout(func() {
				cmd := newListCommand()
				cmd.SetArgs(tc.args)
				cmd.SetOutput(os.Stdout)
				err := cmd.Execute()
				if err == nil && tc.hasError || err != nil && !tc.hasError {
					t.Errorf("%v: status code did not match, wont: %v, got: %v", tc.args, tc.hasError, err)
				}
			})
			if !tc.hasError {
				output = strings.TrimSpace(output)
				output = rrh.ReplaceNewline(output, "+")
				if output != tc.output {
					t.Errorf("%v: output did not match, wont: %s, got: %s", tc.args, tc.output, output)
				}
			}
		})
		defer os.Remove(dbFile)
	}
}

func TestInfoRepository(t *testing.T) {
	var testcases = []struct {
		args     []string
		hasError bool

		output       string
		ignoreOutput bool
	}{
		{[]string{"repo1"}, false, `Repository Id: repo1+Group: group1+Path: path1`, false},
		{[]string{"repo4"}, true, "", true},
	}

	for _, tc := range testcases {
		var dbFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			buffer := bytes.NewBuffer([]byte{})
			cmd := newInfoCommand()
			cmd.SetArgs(tc.args)
			cmd.SetOut(buffer)
			err := cmd.Execute()
			if err == nil && tc.hasError || err != nil && !tc.hasError {
				t.Errorf("%v: status code did not match, wont: %v, got: %v", tc.args, tc.hasError, err)
			}
			output := buffer.String()
			if !tc.ignoreOutput {
				output = strings.TrimSpace(output)
				output = rrh.ReplaceNewline(output, "+")
				if output != tc.output {
					t.Errorf("%v: result did not match, wont: \"%s\", got: \"%s\"", tc.args, tc.output, output)
				}
			}
		})
		defer os.Remove(dbFile)
	}
}

func TestUpdateRepository(t *testing.T) {
	pwd, _ := filepath.Abs(".")
	var testcases = []struct {
		args      []string
		hasError  bool
		newRepoID string
		wontRepo  *rrh.Repository
	}{
		{[]string{"--id", "newRepo1", "--path", "../../../../testdata/helloworld", "--description", "desc1", "repo1"}, false, "newRepo1", &rrh.Repository{ID: "newRepo1", Description: "desc1", Path: filepath.Clean(filepath.Join(pwd, "../../../../testdata/helloworld"))}},
		{[]string{"--description", "desc2", "repo2"}, false, "repo2", &rrh.Repository{ID: "repo2", Description: "desc2", Path: "path2"}},
		{[]string{"repo4"}, true, "repo4", nil},                                        // unknown repository
		{[]string{"--description", "desc", "repo1", "repo3"}, true, "never used", nil}, // too many arguments.
	}
	for _, tc := range testcases {
		var dbFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			cmd := newUpdateCommand()
			cmd.SetOut(os.Stdout)
			cmd.SetArgs(tc.args)
			err := cmd.Execute()
			if err == nil && tc.hasError || err != nil && !tc.hasError {
				t.Errorf("%v: status code did not match, wont: %v, got: %v", tc.args, tc.hasError, err)
			}
			if err != nil {
				return
			}
			var db, _ = rrh.Open(config)
			var repo = db.FindRepository(tc.newRepoID)
			if repo == nil {
				t.Errorf("%s: new repository do not found", tc.newRepoID)
				return
			}
			if repo.ID != tc.wontRepo.ID {
				t.Errorf("%v: id did not match: wont: %s, got: %s", tc.args, tc.wontRepo.ID, repo.ID)
			}
			if repo.Path != tc.wontRepo.Path {
				t.Errorf("%v: path did not match: wont: %s, got: %s", tc.args, tc.wontRepo.Path, repo.Path)
			}
			if repo.Description != tc.wontRepo.Description {
				t.Errorf("%v: description did not match: wont: %s, got: %s", tc.args, tc.wontRepo.Description, repo.Description)
			}
		})
		defer os.Remove(dbFile)
	}
}
