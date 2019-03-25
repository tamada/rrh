package move

import (
	"fmt"
	"os"
	"testing"

	"github.com/tamada/rrh/common"
)

func rollback(dbpath string, f func()) {
	var db = open(dbpath)

	f()

	db.StoreAndClose()
}
func open(jsonName string) *common.Database {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	os.Setenv(common.RrhDatabasePath, fmt.Sprintf("../testdata/%s", jsonName))
	var config = common.OpenConfig()
	var db, _ = common.Open(config)
	return db
}

func TestMoveRepositoryToRepository(t *testing.T) {
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
		{"unrelate, then relate", []string{"-v", "group1/repo1", "group3/repo1"},
			[]relation{{"group3", "repo1", true}, {"group1", "repo1", false}}},
		{"relate only", []string{"repo1", "group3"}, []relation{
			{"group3", "repo1", true},
			{"group1", "repo1", true}}},
		{"relate to new group", []string{"repo1", "group4"}, []relation{
			{"group1", "repo1", true},
			{"group4", "repo1", true}}},
		{"group to group", []string{"group1", "group4"}, []relation{
			{"group4", "repo1", true},
			{"group1", "repo1", false}}},
		{"groups to group", []string{"-verbose", "group1", "group3", "group4"}, []relation{
			{"group4", "repo1", true},
			{"group4", "repo2", true},
			{"group3", "repo2", false},
			{"group1", "repo1", false}}},
	}
	for _, item := range cases {
		rollback("tmp.json", func() {
			var mv, _ = MoveCommandFactory()
			mv.Run(item.args)
			var db = open("tmp.json")
			for _, rel := range item.relations {
				if db.HasRelation(rel.group, rel.repo) != rel.hasRelation {
					t.Errorf("rrh mv %v failed: relation: group %s and repo %s: %v", item.args, rel.group, rel.repo, !rel.hasRelation)
				}
			}
		})
	}
}

func TestParseType(t *testing.T) {
	var cases = []struct {
		gives     string
		wont      int
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

	var db = open("tmp.json")
	for _, item := range cases {
		var got, err = parseType(db, item.gives)
		if got.targetType != item.wont && (item.errorFlag && err == nil) {
			t.Errorf("%s: gives: %v, wont: %d, got: %d", item.message, item.gives, item.wont, got.targetType)
		}
	}
}

func TestVerifyArguments(t *testing.T) {
	var cases = []struct {
		givesFrom []string
		givesTo   string
		wont      int
		errorFlag bool
		message   string
	}{
		{[]string{"group1/repo1"}, "repo5", RepositoriesToGroup, false, "repo5 treats as a group"},
		{[]string{"group3"}, "repo5", GroupToGroup, false, "repo5 treats as a group"},
		{[]string{"group1"}, "group3", GroupToGroup, false, ""},
		{[]string{"repo1"}, "group3", RepositoriesToGroup, false, ""},
		{[]string{"repo1"}, "group5", RepositoriesToGroup, false, ""},
		{[]string{"group1/repo1"}, "repo5/repo1", RepositoryToRepository, false, ""},
		{[]string{"group1/repo1"}, "group2", RepositoriesToGroup, false, ""},
		{[]string{"group1/repo1", "group3/repo2"}, "group1", RepositoriesToGroup, false, ""},
		{[]string{"group1"}, "group3", GroupToGroup, false, ""},
		{[]string{"group1", "group2"}, "group3", GroupsToGroup, false, ""},
		{[]string{"repo1", "repo2"}, "group3/repo1", Invalid, true, "Multiple from moves to only group"},
		{[]string{"repo1"}, "group5/repo1", RepositoryToRepository, false, ""},
	}

	var db = open("tmp.json")
	for _, item := range cases {
		var froms, to = convertToTarget(db, item.givesFrom, item.givesTo)
		var got, _ = verifyArguments(db, froms, to)
		if got != item.wont {
			t.Errorf("%s: gives: %v, %s, wont: %d, got: %d", item.message, item.givesFrom, item.givesTo, item.wont, got)
		}
	}
}

func TestMergeType(t *testing.T) {
	var cases = []struct {
		gives     []int
		wont      int
		errorFlag bool
	}{
		{[]int{GroupType, GroupType, GroupType}, GroupType, false},
		{[]int{GroupType, RepositoryType, GroupType}, Unknown, true},
	}

	for _, item := range cases {
		var got, err = mergeType(item.gives)
		if got != item.wont || (item.errorFlag && err == nil) {
			t.Errorf("gives: %v, wont: %v, got: %v", item.gives, item.wont, got)
		}
	}
}

func TestSynopsis(t *testing.T) {
	var mv, _ = MoveCommandFactory()
	if mv.Synopsis() != "move the repositories from groups to another group" {
		t.Error("Synopsis message is not matched.")
	}
}
func TestHelp(t *testing.T) {
	var mv = MoveCommand{}
	const helpMessage = `rrh mv [OPTIONS] <FROMS...> <TO>
OPTIONS
    -v, --verbose   verbose mode
    -i, --inquiry   inquiry mode

ARGUMENTS
    FROMS...        specifies move from, formatted in <GROUP_NAME/REPO_ID>, or <GROUP_NAME>
    TO              specifies move to, formatted in <GROUP_NAME>`
	if mv.Help() != helpMessage {
		t.Error("Help message is not matched.")
	}
}
