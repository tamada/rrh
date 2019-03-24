package remove

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

func ExampleRemoveCommand_Run() {
	var db = open("tmp.json")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var rm, _ = RemoveCommandFactory()
	rm.Run([]string{"-v", "group2", "repo1"})
	// Output: group2: group removed
	// repo1: repository removed

	db.StoreAndClose()
}

func TestRemoveCommandUnknownGroupAndRepository(t *testing.T) {
	var db = open("tmp.json")

	var rm = RemoveCommand{}
	var err = rm.executeRemove(db, "not_exist_group_and_repository")
	if err.Error() != "not_exist_group_and_repository: not found in repositories and groups" {
		t.Error("not exist group and repository found!?")
	}

	db.StoreAndClose()
}

func TestRemoveRepository(t *testing.T) {
	var db = open("tmp.json")
	var rm = RemoveCommand{&removeOptions{}}
	if err := rm.executeRemoveRepository(db, "unknown-repo"); err == nil {
		t.Error("unknown-repo: found")
	}
	var err = rm.executeRemoveRepository(db, "repo1")
	if err != nil || len(db.Repositories) != 1 {
		t.Errorf("repo1 did not remove?: %s", err.Error())
	}
	if len(db.Groups) != 2 {
		t.Error("the number of groups changed")
	}
	if db.ContainsCount("group1") != 0 || db.ContainsCount("group2") != 0 {
		t.Error("Unrelate repo was failed?")
	}
}

func TestRemoveGroup(t *testing.T) {
	var db = open("tmp.json")
	var rm = RemoveCommand{&removeOptions{}}
	if err := rm.executeRemoveGroup(db, "unknown-group"); err == nil {
		t.Error("unknown-group: found")
	}
	if err := rm.executeRemoveGroup(db, "group1"); err == nil {
		t.Error("group1 has a entry!")
	}
	rm.Options.recursive = true
	if err := rm.executeRemoveGroup(db, "group1"); err != nil {
		t.Error("group1 cannot remove recursively.")
	}
}

func TestRemoveCommandRemoveTargetIsBothInGroupAndRepository(t *testing.T) {
	var db = open("nulldb.json")

	db.CreateGroup("groupOrRepo", "same name as Repository")
	db.CreateRepository("groupOrRepo", "unknownpath", []common.Remote{})
	var rm = RemoveCommand{&removeOptions{}}
	var err = rm.executeRemove(db, "groupOrRepo")
	if err.Error() != "groupOrRepo: exists in repositories and groups" {
		t.Error("not failed!?")
	}
}

func TestRemoveEntryFailed(t *testing.T) {
	var db = open("tmp.json")
	var rm = RemoveCommand{&removeOptions{}}
	var err = rm.executeRemoveFromGroup(db, "group2", "repo2")
	if err != nil {
		t.Error("Successfully remove unrelated group and repository.")
	}
}

func TestRemoveRelation(t *testing.T) {
	var db = open("tmp.json")

	var rm, _ = RemoveCommandFactory()
	rm.Run([]string{"-v", "group1/repo1"})
	var db2 = open("tmp.json")
	if len(db2.Repositories) != 2 && len(db2.Groups) != 2 {
		t.Error("repositories and groups are removed!")
	}
	if db2.ContainsCount("group1") != 0 || db2.ContainsCount("group2") != 0 {
		t.Error("relation was not removed")
	}

	db.StoreAndClose()
}

func TestRunRemoveRepository(t *testing.T) {
	var db = open("tmp.json")

	var rm, _ = RemoveCommandFactory()
	rm.Run([]string{"-v", "group2", "repo1"})
	var db2 = open("tmp.json")
	if len(db2.Repositories) != 1 && len(db2.Groups) != 1 {
		t.Errorf("repositories: %d, groups: %d\n", len(db2.Repositories), len(db2.Groups))
	}
	if db2.ContainsCount("group1") != 0 {
		t.Errorf("database was broken")
	}

	db.StoreAndClose()
}

func TestRemoveRepository2(t *testing.T) {
	var db = open("tmp.json")

	var rm, _ = RemoveCommandFactory()
	os.Setenv(common.RrhAutoDeleteGroup, "true")
	rm.Run([]string{"-v", "group2", "repo1"})
	var db2 = open("tmp.json")
	if len(db2.Repositories) != 1 && len(db2.Groups) != 0 {
		t.Errorf("repositories: %d, groups: %d\n", len(db2.Repositories), len(db2.Groups))
	}

	db.StoreAndClose()
}

func TestBrokenDatabase(t *testing.T) {
	os.Setenv(common.RrhDatabasePath, "../testdata/broken.json")
	var rm, _ = RemoveCommandFactory()
	if result := rm.Run([]string{}); result != 2 {
		t.Errorf("broken database are successfully read!?")
	}
}

func TestUnknownOptions(t *testing.T) {
	var db = open("tmp.json")

	var rm, _ = RemoveCommandFactory()
	if result := rm.Run([]string{"--unknown"}); result != 1 {
		t.Errorf("unknown option was not failed: %d", result)
	}

	db.StoreAndClose()
}

func TestHelp(t *testing.T) {
	var helpMessage = `rrh rm [OPTIONS] <REPO_ID|GROUP_ID|REPO_ID/GROUP_ID...>
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

func TestSynopsis(t *testing.T) {
	var rm, _ = RemoveCommandFactory()
	if rm.Synopsis() != "remove given repository from database." {
		t.Error("synopsis did not match")
	}
}

func TestRemove(t *testing.T) {
	// var db = open("tmp.json")
	// var rm, _ = RemoveCommandFactory()

}
