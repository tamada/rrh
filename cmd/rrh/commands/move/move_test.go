package move

import (
	"bytes"
	"os"
	"testing"

	"github.com/tamada/rrh"
	"github.com/tamada/rrh/common"
)

func TestParseError(t *testing.T) {
	testcases := []struct {
		args      []string
		wontError bool
	}{
		{[]string{"--unknown-option"}, true},
		{[]string{"too_few_arguments"}, true},
		{[]string{"group1/repo1", "group3/repo5"}, true},
		{[]string{"group1/repo1", "repo2", "group5"}, true},
	}
	for _, testcase := range testcases {
		buffer := bytes.NewBuffer([]byte{})
		var dbFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			cmd := New()
			cmd.SetOut(buffer)
			cmd.SetArgs(testcase.args)
			err := cmd.Execute()
			if err != nil && !testcase.wontError || err == nil && testcase.wontError {
				t.Errorf("%v: wont error %v, but got %v", testcase.args, testcase.wontError, err)
			}
		})
		os.Remove(dbFile)
	}
}

func TestMoveCommand(t *testing.T) {
	type relation struct {
		group       string
		repo        string
		hasRelation bool
	}
	var cases = []struct {
		message   string
		args      []string
		relations []relation
	}{
		{"unrelate, then relate", []string{"group1/repo1", "group3/repo1"},
			[]relation{{"group3", "repo1", true}, {"group1", "repo1", false}}},
		{"unrelate, then relate", []string{"group1/repo1", "group5"},
			[]relation{{"group5", "repo1", true}, {"group1", "repo1", false}}},
		{"different repository name", []string{"group1/repo1", "group5/repo1"},
			[]relation{{"group5", "repo1", true}, {"group1", "repo1", false}}},
		{"relate only", []string{"repo1", "group3"}, []relation{
			{"group3", "repo1", true},
			{"group1", "repo1", true}}},
		{"relate to new group", []string{"repo1", "group4"}, []relation{
			{"group1", "repo1", true},
			{"group4", "repo1", true}}},
		{"group to group", []string{"group1", "group4"}, []relation{
			{"group4", "repo1", true},
			{"group1", "repo1", false}}},
		{"groups to group", []string{"group1", "group3", "group4"}, []relation{
			{"group4", "repo1", true},
			{"group4", "repo2", true},
			{"group3", "repo2", false},
			{"group1", "repo1", false}}},
	}
	for _, item := range cases {
		dbFile := rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			cmd := New()
			cmd.SetArgs(item.args)
			cmd.Execute()

			var db, _ = rrh.Open(config)
			for _, rel := range item.relations {
				if db.HasRelation(rel.group, rel.repo) != rel.hasRelation {
					t.Errorf("rrh mv %v failed: relation: group %s and repo %s: %v", item.args, rel.group, rel.repo, !rel.hasRelation)
				}
			}
		})
		defer os.Remove(dbFile)
	}
}

func TestParseType(t *testing.T) {
	var cases = []struct {
		gives     string
		wont      targetKind
		errorFlag bool
		message   string
	}{
		{"group1/repo1", GroupAndRepoType, false, ""},
		{"not-exist-group/not-exist-repo", Unknown, true, "group and repository not found"},
		{"not-exist-group/repo1", Unknown, true, "group not found"},
		{"group1/not-exist-repo", Unknown, true, "repository not found"},
		{"group3/repo1", GroupAndRepoType, true, "no relation between group3 and repo1"},
		{"group1", GroupType, false, ""},
		{"repo1", RepositoryType, false, "should be <GROUP/REPO> or <GROUP>"},
		{"not-exist", GroupOrRepoType, false, "group not found"},
	}

	var dbFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		for _, item := range cases {
			var got, err = parseType(db, item.gives)
			if got.kind != item.wont && (item.errorFlag && err == nil) {
				t.Errorf("%s: gives: %v, wont: %d, got: %d", item.message, item.gives, item.wont, got.kind)
			}
		}
	})
	defer os.Remove(dbFile)
}

func TestVerifyArguments(t *testing.T) {
	var cases = []struct {
		givesFrom []string
		givesTo   string
		wont      executionType
		errorFlag bool
		message   string
	}{
		{[]string{"group1/repo1"}, "repo5", RepositoriesToGroup, false, "repo5 treats as a group"},
		{[]string{"group3"}, "repo5", GroupToGroup, false, "repo5 treats as a group"},
		{[]string{"group1"}, "group3", GroupToGroup, false, ""},
		{[]string{"repo1"}, "group3", RepositoriesToGroup, false, ""},
		{[]string{"repo1"}, "group5", RepositoriesToGroup, false, ""},
		{[]string{"group1/repo1"}, "repo5/repo1", RepositoryToRepository, false, ""},
		{[]string{"group1/repo1"}, "group3/repo5", RepositoryToRepository, false, ""},
		{[]string{"group1/repo1"}, "group2", RepositoriesToGroup, false, ""},
		{[]string{"group1/repo1", "group3/repo2"}, "group1", RepositoriesToGroup, false, ""},
		{[]string{"group1"}, "group3", GroupToGroup, false, ""},
		{[]string{"group1", "group2"}, "group3", GroupsToGroup, false, ""},
		{[]string{"repo1", "repo2"}, "group3/repo1", Invalid, true, "Multiple from moves to only group"},
		{[]string{"repo1"}, "group5/repo1", RepositoryToRepository, false, ""},
	}

	var dbFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		for _, item := range cases {
			var froms, to = convertToTarget(db, item.givesFrom, item.givesTo)
			var got, _ = verifyArguments(db, froms, to)
			if got != item.wont {
				t.Errorf("%s: gives: %v, %s, wont: %d, got: %d", item.message, item.givesFrom, item.givesTo, item.wont, got)
			}
		}
	})
	defer os.Remove(dbFile)
}

func TestMergeType(t *testing.T) {
	var cases = []struct {
		gives     []targetKind
		wont      targetKind
		errorFlag bool
	}{
		{[]targetKind{GroupType, GroupType, GroupType}, GroupType, false},
		{[]targetKind{GroupType, RepositoryType, GroupType}, Unknown, true},
	}

	for _, item := range cases {
		var got, err = mergeType(item.gives)
		if got != item.wont || (item.errorFlag && err == nil) {
			t.Errorf("gives: %v, wont: %v, got: %v", item.gives, item.wont, got)
		}
	}
}

func TestErrorOnPerformImpl(t *testing.T) {
	err := performImpl(nil, targets{}, Invalid)
	el, ok := err.(common.ErrorList)
	if ok && len(el) != 1 {
		t.Errorf("return code of performImpl did not match, wont: 1, got: %d", len(el))
	}
}

func TestVerifyArgumentsOneToOne(t *testing.T) {
	var testcases = []struct {
		fromType    targetKind
		toType      targetKind
		resultType  executionType
		shouldError bool
	}{
		{Unknown, Unknown, Invalid, true},
		{GroupType, RepositoryType, Invalid, true},
	}
	for _, tc := range testcases {
		var dbFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			var db, _ = rrh.Open(config)
			var resultType, err = verifyArgumentsOneToOne(db, target{kind: tc.fromType}, target{kind: tc.toType})
			if resultType != tc.resultType {
				t.Errorf("%v: result type did not match, wont: %d, got: %d", tc, tc.resultType, resultType)
			}
			if (err == nil) == tc.shouldError {
				t.Errorf("verifyArgumentsOneToOne(%d, %d): should error, wont: %v, got: %v", tc.fromType, tc.toType, tc.shouldError, !tc.shouldError)
			}
		})
		defer os.Remove(dbFile)
	}
}
