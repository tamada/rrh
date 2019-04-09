package prune

import (
	"os"
	"testing"

	"github.com/tamada/rrh/common"
)

func open() *common.Database {
	var config = common.OpenConfig()
	var db, _ = common.Open(config)
	return db
}

func TestSynopsis(t *testing.T) {
	var prune, _ = CommandFactory()
	if prune.Synopsis() != "prune unnecessary repositories and groups." {
		t.Error("Synopsis message is not matched.")
	}
}
func TestHelp(t *testing.T) {
	var prune = Command{}
	if prune.Help() != "rrh prune" {
		t.Error("Help message is not matched.")
	}
}

type groupExistChecker struct {
	groupName string
	existFlag bool
}

type repositoryExistChecker struct {
	repoName  string
	existFlag bool
}

func TestPrune(t *testing.T) {
	var tc = struct {
		gchecker []groupExistChecker
		rchecker []repositoryExistChecker
	}{
		[]groupExistChecker{{"group1", true}, {"group2", false}, {"group3", true}},
		[]repositoryExistChecker{{"repo1", true}, {"repo2", true}, {"repo3", false}},
	}
	common.WithDatabase("../testdata/tmp.json", "../testdata/config.json", func() {
		var db = open()
		db.Prune()

		for _, gc := range tc.gchecker {
			if db.HasGroup(gc.groupName) != gc.existFlag {
				t.Errorf("group %s exist flag did not match: wont: %v, got: %v", gc.groupName, gc.existFlag, !gc.existFlag)
			}
		}
		for _, rc := range tc.rchecker {
			if db.HasRepository(rc.repoName) != rc.existFlag {
				t.Errorf("repository %s exist flag did not match: wont: %v, got: %v", rc.repoName, rc.existFlag, !rc.existFlag)
			}
		}
	})
}

func TestCommandRunFailedByBrokenDBFile(t *testing.T) {
	os.Setenv(common.RrhDatabasePath, "../testdata/broken.json")
	var prune, _ = CommandFactory()
	if prune.Run([]string{}) != 1 {
		t.Error("broken database read successfully.")
	}
}

func ExampleCommand_Run() {
	common.Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
		var prune, _ = CommandFactory()
		prune.Run([]string{})
	})
	// Output: Pruned 3 groups, 2 repositories
}
