package group

import (
	"os"
	"testing"

	"github.com/tamada/rrh"
)

func TestGroupRemove(t *testing.T) {
	testcases := []struct {
		args       []string
		wontError  bool
		existGroup bool
	}{
		{[]string{"rm", "unknown_group"}, true, false},
		{[]string{"rm", "group2"}, false, false},
		{[]string{"rm", "group1"}, true, true},
		{[]string{"rm", "group1", "--force"}, false, false},
	}
	for _, tc := range testcases {
		dbFile := rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
			cmd := New()
			cmd.SetArgs(tc.args)
			cmd.SetOut(os.Stdout)
			err := cmd.Execute()

			if err == nil && tc.wontError || err != nil && !tc.wontError {
				t.Errorf("%v: group %v, wont error: %v, but got error %v", tc.args, tc.args, tc.wontError, err != nil)
			}

			db2, _ := rrh.Open(config)
			group := db2.FindGroup(tc.args[1])
			if group == nil && tc.existGroup || group != nil && !tc.existGroup {
				t.Errorf("%v: %s wont exist flag %v, but existence is %v", tc.args, tc.args[1], tc.existGroup, group != nil)
			}
		})
		defer os.Remove(dbFile)
	}
}

func TestGroupOf(t *testing.T) {
	testcases := []struct {
		args      []string
		wontError bool
	}{
		{[]string{"of", "not_exist_repo"}, true},
	}

	for _, tc := range testcases {
		dbFile := rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
			cmd := New()
			cmd.SetArgs(tc.args)
			cmd.SetOut(os.Stdout)
			err := cmd.Execute()
			if err == nil && tc.wontError || err != nil && !tc.wontError {
				t.Errorf("group %v wont error: %v, but got error: %v", tc.args, tc.wontError, err != nil)
			}
		})
		defer os.Remove(dbFile)
	}
}

type groupChecker struct {
	groupName   string
	existFlag   bool
	description string
	abbrevFlag  bool
}

func TestGroupAdd(t *testing.T) {
	var testcases = []struct {
		args       []string
		statusCode int
		checkers   []groupChecker
	}{
		{[]string{"add", "--note", "desc4", "group4"}, 0, []groupChecker{{"group4", true, "desc4", false}}},
		{[]string{"add", "-n", "desc4", "-a", "true", "group4"}, 0, []groupChecker{{"group4", true, "desc4", true}}},
		{[]string{"add", "-n", "desc4", "--abbrev", "false", "group4"}, 0, []groupChecker{{"group4", true, "desc4", false}}},
		{[]string{"add", "-n", "desc4", "-a", "true", "group1"}, 4, []groupChecker{}},
		{[]string{"add"}, 3, []groupChecker{}},
	}
	for _, testcase := range testcases {
		dbFile := rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
			cmd := New()
			cmd.SetArgs(testcase.args)
			cmd.SetOut(os.Stdout)
			cmd.Execute()

			db2, _ := rrh.Open(config)
			for _, checker := range testcase.checkers {
				if db2.HasGroup(checker.groupName) != checker.existFlag {
					t.Errorf("%v: group check failed: %s, wont: %v, got: %v", testcase.args, checker.groupName, checker.existFlag, !checker.existFlag)
				}
				if checker.existFlag {
					var group = db2.FindGroup(checker.groupName)
					if group != nil && group.Description != checker.description {
						t.Errorf("%v: group description did not match: wont: %s, got: %s", testcase.args, checker.description, group.Description)
					}
					if group != nil && group.OmitList != checker.abbrevFlag {
						t.Errorf("%v: group OmitList did not match: wont: %v, got: %v", testcase.args, checker.abbrevFlag, group.OmitList)
					}
				}
			}

		})
		defer os.Remove(dbFile)
	}

}

func ExampleGroupListCommand_Run() {
	dbFile := rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		cmd := New()
		cmd.SetArgs([]string{"list", "-f", "csv", "--no-header", "-e", "name,count,repo,note"})
		cmd.SetOut(os.Stdout)
		cmd.Execute()
	})
	defer os.Remove(dbFile)
	// Output:
	// group1,desc1,[repo1],1 repository
	// group2,desc2,[],0 repositories
	// group3,desc3,[repo2],1 repository
}

func ExampleGroupOfCommand_Run() {
	dbFile := rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		cmd := New()
		cmd.SetArgs([]string{"of", "repo1"})
		cmd.SetOut(os.Stdout)
		cmd.Execute()
	})
	defer os.Remove(dbFile)
	// Output:
	// repo1: group1
}

func ExampleGroupInfoCommand_Run() {
	dbFile := rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		cmd := New()
		cmd.SetArgs([]string{"info", "group1", "group2", "groupN"})
		cmd.SetOut(os.Stdout)
		cmd.SetErr(os.Stdout)
		cmd.Execute()
	})
	defer os.Remove(dbFile)
	// Output:
	// group1: desc1 (1 repository, abbrev: false)
	// group2: desc2 (0 repositories, abbrev: false)
	// Error: groupN: group not found
}
