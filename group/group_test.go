package group

import (
	"os"
	"testing"

	"github.com/tamada/rrh/common"
)

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

func TestInvalidOptionInGroupList(t *testing.T) {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var glc, _ = groupListCommandFactory()
	if val := glc.Run([]string{"--unknown-option"}); val != 1 {
		t.Error("list subcommand accept unknown-option!")
	}
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
    GROUP                gives group names.`

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
