package group

import (
	"os"
	"testing"

	"github.com/tamada/rrh/common"
)

func rollback(dbpath string, f func()) {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	os.Setenv(common.RrhDatabasePath, dbpath)
	var config = common.OpenConfig()
	var db, _ = common.Open(config)

	f()

	db.StoreAndClose()
}

func ExampleGroupCommand_Run() {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var gc, _ = GroupCommandFactory()
	gc.Run([]string{})
	// Output:
	// group1,(1 repository)
	// group2,(0 repositories)
}

func Example_groupListCommand_Run() {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var glc, _ = groupListCommandFactory()
	glc.Run([]string{"-d", "-r"})
	// Output:
	// group1,desc1,[repo1],(1 repository)
	// group2,desc2,[],(0 repositories)
}

func TestAddGroup(t *testing.T) {
	rollback("../testdata/tmp.json", func() {
		var gac, _ = groupAddCommandFactory()
		if val := gac.Run([]string{"-d", "desc3", "group3"}); val != 0 {
			t.Errorf("group add failed: %d", val)
		}
		var config = common.OpenConfig()
		var db2, _ = common.Open(config)
		if len(db2.Groups) != 3 {
			t.Fatal("group3 was not added.")
		}
		if db2.Groups[2].Name != "group3" || db2.Groups[2].Description != "desc3" {
			t.Errorf("want: group3 (desc3), got: %s (%s)", db2.Groups[2].Name, db2.Groups[2].Description)
		}
	})
}

func TestUpdateGroup(t *testing.T) {
	rollback("../testdata/tmp.json", func() {
		var guc, _ = groupUpdateCommandFactory()
		if val := guc.Run([]string{"-d", "newdesc2", "--name", "newgroup2", "group2"}); val != 0 {
			t.Errorf("group update failed: %d", val)
		}
		var config = common.OpenConfig()
		var db2, _ = common.Open(config)
		if len(db2.Groups) != 2 {
			t.Fatal("the length of group did not match")
		}
		if db2.Groups[1].Name != "newgroup2" || db2.Groups[1].Description != "newdesc2" {
			t.Errorf("want: newgroup2 (newdesc2), got: %s (%s)", db2.Groups[1].Name, db2.Groups[1].Description)
		}
	})
}

func TestRemoveGroup(t *testing.T) {
	rollback("../testdata/tmp.json", func() {
		var grc, _ = groupRemoveCommandFactory()
		if val := grc.Run([]string{"--force", "group1"}); val != 0 {
			t.Errorf("group remove failed: %d", val)
		}
		var config = common.OpenConfig()
		var db2, _ = common.Open(config)
		if len(db2.Groups) != 1 {
			t.Fatalf("the length of group did not match: %v", db2.Groups)
		}
		if db2.Groups[0].Name != "group2" || db2.Groups[0].Description != "desc2" {
			t.Errorf("want: group2 (desc2), got: %s (%s)", db2.Groups[0].Name, db2.Groups[0].Description)
		}
	})
}

func TestInvalidOptionInGroupList(t *testing.T) {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	common.CaptureStdout(func() {
		var glc, _ = groupListCommandFactory()
		if val := glc.Run([]string{"--unknown-option"}); val != 1 {
			t.Error("list subcommand accept unknown-option!")
		}
		var gac, _ = groupAddCommandFactory()
		if val := gac.Run([]string{"--unknown-option"}); val != 1 {
			t.Error("add subcommand accept unknown-option!")
		}
		var grc, _ = groupRemoveCommandFactory()
		if val := grc.Run([]string{"--unknown-option"}); val != 1 {
			t.Error("remove subcommand accept unknown-option!")
		}
		var guc, _ = groupUpdateCommandFactory()
		if val := guc.Run([]string{"--unknown-option"}); val != 1 {
			t.Error("update subcommand accept unknown-option!")
		}
	})
}

func TestHelp(t *testing.T) {
	var gac, _ = groupAddCommandFactory()
	var glc, _ = groupListCommandFactory()
	var grc, _ = groupRemoveCommandFactory()
	var guc, _ = groupUpdateCommandFactory()
	var gc, _ = GroupCommandFactory()

	var gacHelp = `rrh group add [OPTIONS] <GROUPS...>
OPTIONS
    -d, --desc <DESC>    give the description of the group
ARGUMENTS
    GROUPS               gives group names.`

	var glcHelp = `rrh group list [OPTIONS]
OPTIONS
    -d, --desc          show description.
    -r, --repository    show repositories in the group.`

	var grcHelp = `rrh group rm [OPTIONS] <GROUPS...>
OPTIONS
    -f, --force      force remove
    -i, --inquery    inquiry mode
    -v, --verbose    verbose mode
ARGUMENTS
    GROUPS           target group names.`

	var gucHelp = `rrh group update [OPTIONS] <GROUP>
OPTIONS
    -n, --name <NAME>   change group name to NAME.
    -d, --desc <DESC>   change description to DESC.
ARGUMENTS
    GROUP               update target group names.`

	var gcHelp = `rrh group <SUBCOMMAND>
SUBCOMMAND
    add       add new group.
    list      list groups (default).
    rm        remove group.
    update    update group`

	if gc.Help() != gcHelp {
		t.Error("help message did not match")
	}
	if glc.Help() != glcHelp {
		t.Error("help message did not match")
	}
	if guc.Help() != gucHelp {
		t.Error("help message did not match")
	}
	if gac.Help() != gacHelp {
		t.Error("help message did not match")
	}
	if grc.Help() != grcHelp {
		t.Error("help message did not match")
	}
}

func TestSynopsis(t *testing.T) {
	var gc, _ = GroupCommandFactory()
	if gc.Synopsis() != "add/list/update/remove groups." {
		t.Error("synopsis did not match")
	}

	var guc, _ = groupUpdateCommandFactory()
	if guc.Synopsis() != "update group." {
		t.Error("synopsis did not match")
	}
	var grc, _ = groupRemoveCommandFactory()
	if grc.Synopsis() != "remove given group." {
		t.Error("synopsis did not match")
	}
	var gac, _ = groupAddCommandFactory()
	if gac.Synopsis() != "add group." {
		t.Error("synopsis did not match")
	}
	var glc, _ = groupListCommandFactory()
	if glc.Synopsis() != "list groups." {
		t.Error("synopsis did not match")
	}
}
