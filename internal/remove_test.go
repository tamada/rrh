package internal

import (
	"os"
	"testing"

	"github.com/tamada/rrh"
)

func ExampleRemoveCommand_Run() {
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
		var rm, _ = RemoveCommandFactory()
		rm.Run([]string{"-v", "group2", "repo1"})
	})
	defer os.Remove(dbFile)
	// Output: group2: group removed
	// repo1: repository removed
}

func TestCommandUnknownGroupAndRepository(t *testing.T) {
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		var rm = RemoveCommand{}
		var err = rm.executeRemove(db, "not_exist_group_and_repository")
		if err.Error() != "not_exist_group_and_repository: not found in repositories and groups" {
			t.Error("not exist group and repository found!?")
		}
	})
	defer os.Remove(dbFile)
}

func TestRemoveRepository(t *testing.T) {
	var testcases = []struct {
		repositoryName   string
		removeSuccess    bool
		belongedGroup    string
		repoCountInGroup int
	}{
		{"unknown-repo", false, "", 0},
		{"repo1", true, "group1", 0},
	}
	for _, tc := range testcases {
		var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
			var rm = RemoveCommand{&removeOptions{}}
			var err = rm.executeRemoveRepository(db, tc.repositoryName)
			if (err == nil) != tc.removeSuccess {
				t.Errorf("%v: remove result not match: wont: %v, got: %v", tc.repositoryName, tc.removeSuccess, !tc.removeSuccess)
			}
			if tc.belongedGroup != "" {
				var count = db.ContainsCount(tc.belongedGroup)
				if count != tc.repoCountInGroup {
					t.Errorf("%v: repo count in group %s did not match: wont: %d, got: %d", tc.repositoryName, tc.belongedGroup, tc.repoCountInGroup, count)
				}
			}
		})
		defer os.Remove(dbFile)
	}
}

func TestRemoveCommandForGroup(t *testing.T) {
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		var rm = RemoveCommand{&removeOptions{}}
		if err := rm.executeRemoveGroup(db, "unknown-group"); err == nil {
			t.Error("unknown-group: found")
		}
		if err := rm.executeRemoveGroup(db, "group1"); err == nil {
			t.Error("group1 has a entry!")
		}
		rm.options.recursive = true
		if err := rm.executeRemoveGroup(db, "group1"); err != nil {
			t.Error("group1 cannot remove recursively.")
		}
	})
	defer os.Remove(dbFile)
}

func TestRemoveCommandRemoveTargetIsBothInGroupAndRepository(t *testing.T) {
	var dbFile = rrh.Rollback("../testdata/nulldb.json", "../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		db.CreateGroup("groupOrRepo", "same name as Repository", false)
		db.CreateRepository("groupOrRepo", "unknownpath", "desc", []*rrh.Remote{})
		var rm = RemoveCommand{&removeOptions{}}
		var err = rm.executeRemove(db, "groupOrRepo")
		if err.Error() != "groupOrRepo: exists in repositories and groups" {
			t.Error("not failed!?")
		}
	})
	defer os.Remove(dbFile)
}

func TestRemoveEntryFailed(t *testing.T) {
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		var rm = RemoveCommand{&removeOptions{}}
		var err = rm.executeRemoveFromGroup(db, "group2", "repo2")
		if err != nil {
			t.Error("Successfully remove unrelated group and repository.")
		}
	})
	defer os.Remove(dbFile)
}

func TestRemoveRelation(t *testing.T) {
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
		var rm, _ = RemoveCommandFactory()
		rm.Run([]string{"-v", "group1/repo1"})
		var db2, _ = rrh.Open(config)
		if len(db2.Repositories) != 2 && len(db2.Groups) != 3 {
			t.Error("repositories and groups are removed!")
		}
		if db2.ContainsCount("group1") != 0 || db2.ContainsCount("group2") != 0 {
			t.Error("relation was not removed")
		}
	})
	defer os.Remove(dbFile)
}

func TestRunRemoveRepository(t *testing.T) {
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
		var rm, _ = RemoveCommandFactory()
		rm.Run([]string{"-v", "group2", "repo1"})
		var db2, _ = rrh.Open(config)
		if len(db2.Repositories) != 1 && len(db2.Groups) != 2 {
			t.Errorf("repositories: %d, groups: %d\n", len(db2.Repositories), len(db2.Groups))
		}
		if db2.ContainsCount("group1") != 0 {
			t.Errorf("database was broken")
		}
	})
	defer os.Remove(dbFile)
}

func TestRemoveRepository2(t *testing.T) {
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
		var rm, _ = RemoveCommandFactory()
		os.Setenv(rrh.AutoDeleteGroup, "true")
		rm.Run([]string{"-v", "group2", "repo1"})
		var db2, _ = rrh.Open(config)
		if len(db2.Repositories) != 1 && len(db2.Groups) != 0 {
			t.Errorf("repositories: %d, groups: %d\n", len(db2.Repositories), len(db2.Groups))
		}
	})
	defer os.Remove(dbFile)
}

func TestBrokenDatabaseOnRemoveCommand(t *testing.T) {
	os.Setenv(rrh.DatabasePath, "../testdata/broken.json")
	var rm, _ = RemoveCommandFactory()
	if result := rm.Run([]string{}); result != 2 {
		t.Errorf("broken database are successfully read!?")
	}
}

func TestUnknownOptionsOnRemoveCommand(t *testing.T) {
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
		var rm, _ = RemoveCommandFactory()
		if result := rm.Run([]string{"--unknown"}); result != 1 {
			t.Errorf("unknown option was not failed: %d", result)
		}
	})
	defer os.Remove(dbFile)
}

func TestHelpOfRemove(t *testing.T) {
	var helpMessage = `rrh rm [OPTIONS] <REPO_ID|GROUP_ID|GROUP_ID/REPO_ID...>
OPTIONS
    -i, --inquiry       inquiry mode.
    -r, --recursive     recursive mode.
    -v, --verbose       verbose mode.

ARGUMENTS
    REPOY_ID            repository name for removing.
    GROUP_ID            group name. if the group contains repositories,
                        remove will fail without '-r' option.
    GROUP_ID/REPO_ID    remove the relation between the given REPO_ID and GROUP_ID.`
	var rm, _ = RemoveCommandFactory()
	if rm.Help() != helpMessage {
		t.Error("help message did not match.")
	}
}

func TestSynopsisOfRemove(t *testing.T) {
	var rm, _ = RemoveCommandFactory()
	if rm.Synopsis() != "remove given repository from database." {
		t.Error("synopsis did not match")
	}
}

func TestRemove(t *testing.T) {
	// var db = open("test_db.json")
	// var rm, _ = RemoveCommandFactory()

}
