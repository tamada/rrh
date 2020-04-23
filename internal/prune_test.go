package internal

import (
	"os"
	"testing"

	"github.com/tamada/rrh"
)

func TestSynopsis(t *testing.T) {
	var prune, _ = PruneCommandFactory()
	if prune.Synopsis() != "prune unnecessary repositories and groups." {
		t.Error("Synopsis message is not matched.")
	}
}
func TestHelp(t *testing.T) {
	var prune = PruneCommand{}
	if prune.Help() != `rrh prune [OPTIONS]
OPTIONS
    -d, --dry-run    dry-run mode.
    -v, --verbose    verbose mode.` {
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
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
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
	defer os.Remove(dbFile)
}

func TestPruneCommandRunFailedByBrokenDBFile(t *testing.T) {
	os.Setenv(rrh.RrhDatabasePath, "../testdata/broken.json")
	var prune, _ = PruneCommandFactory()
	if prune.Run([]string{}) != 1 {
		t.Error("broken database read successfully.")
	}
}

func TestPruneCommandRunFailedByInvalidArgs(t *testing.T) {
	os.Setenv(rrh.RrhDatabasePath, "../testdata/test_db.json")
	var prune, _ = PruneCommandFactory()
	if prune.Run([]string{"--help"}) != 1 {
		t.Error("successing invalid option parsing.")
	}
}

func ExamplePruneCommand_Run() {
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		var prune, _ = PruneCommandFactory()
		prune.Run([]string{})
	})
	defer os.Remove(dbFile)
	// Output: Pruned 3 groups, 2 repositories
}

func ExamplePruneCommand_Run_DryrunMode() {
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		var prune, _ = PruneCommandFactory()
		prune.Run([]string{"--dry-run"})
	})
	defer os.Remove(dbFile)
	// Output: Pruned 3 groups, 2 repositories (dry-run mode)
	// repo1: repository pruned (not exists)
	// repo2: repository pruned (not exists)
	// group1: group pruned (no relations)
	// group2: group pruned (no relations)
	// group3: group pruned (no relations)
}

func ExamplePruneCommand_Run_VerboseMode() {
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		var prune, _ = PruneCommandFactory()
		prune.Run([]string{"--verbose"})
	})
	defer os.Remove(dbFile)
	// Output: Pruned 3 groups, 2 repositories
	// repo1: repository pruned (not exists)
	// repo2: repository pruned (not exists)
	// group1: group pruned (no relations)
	// group2: group pruned (no relations)
	// group3: group pruned (no relations)
}
