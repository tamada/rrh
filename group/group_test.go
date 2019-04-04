package group

import (
	"os"
	"strings"
	"testing"

	"github.com/tamada/rrh/common"
)

func Example() {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var gc, _ = GroupCommandFactory()
	gc.Run([]string{})
	// Output:
	// group1,1 repository
	// group2,0 repositories
	// group3,1 repository
}

func ExampleGroupCommand_Run() {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var gc, _ = GroupCommandFactory()
	gc.Run([]string{"list"})
	// Output:
	// group1,1 repository
	// group2,0 repositories
	// group3,1 repository
}

func Example_groupListCommand_Run() {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var glc, _ = groupListCommandFactory()
	glc.Run([]string{"-d", "-r"})
	// Output:
	// group1,desc1,[repo1],1 repository
	// group2,desc2,[],0 repositories
	// group3,desc3,[repo2],1 repository
}

func Example_groupOfCommand_Run() {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var goc, _ = groupOfCommandFactory()
	goc.Run([]string{"repo1"})
	// Output:
	// repo1, [group1]
}

func TestGroupListOnlyName(t *testing.T) {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var output = common.CaptureStdout(func() {
		var glc, _ = GroupCommandFactory()
		glc.Run([]string{"list", "--only-groupname"})
	})
	var wontOutput = `group1
group2
group3`
	if strings.TrimSpace(output) != wontOutput {
		t.Errorf("the result with option only-groupname did not match\nwont: %s, got: %s", wontOutput, output)
	}
}

func TestGroupOfCommand(t *testing.T) {
	var testcases = []struct {
		args   []string
		output string
	}{
		{[]string{"unknown-repo"}, "unknown-repo: repository not found"},
		{[]string{"repo2"}, "repo2, [group3]"},
		{[]string{}, `rrh group of <REPOSITORY_ID>
ARGUMENTS
    REPOSITORY_ID     show the groups of the repository.`},
	}
	for _, tc := range testcases {
		var output = common.CaptureStdout(func() {
			var command, _ = groupOfCommandFactory()
			command.Run(tc.args)
		})
		output = strings.TrimSpace(output)
		if output != tc.output {
			t.Errorf("%v: output did not match, wont: %s, got: %s", tc.args, tc.output, output)
		}
	}
}

type groupChecker struct {
	groupName   string
	existFlag   bool
	description string
	omitList    bool
}

func TestAddGroup(t *testing.T) {
	var testcases = []struct {
		args       []string
		statusCode int
		checkers   []groupChecker
	}{
		{[]string{"add", "-d", "desc4", "group4"}, 0, []groupChecker{{"group4", true, "desc4", false}}},
		{[]string{"add", "-d", "desc4", "-o", "true", "group4"}, 0, []groupChecker{{"group4", true, "desc4", true}}},
		{[]string{"add", "-d", "desc4", "-o", "hoge", "group4"}, 0, []groupChecker{{"group4", true, "desc4", false}}},
		{[]string{"add", "-d", "desc4", "-o", "hoge", "group1"}, 4, []groupChecker{}},
		{[]string{"add"}, 3, []groupChecker{}},
	}
	for _, testcase := range testcases {
		common.Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
			var gac, _ = GroupCommandFactory()
			if val := gac.Run(testcase.args); val != testcase.statusCode {
				t.Errorf("%v: test failed, wont: %d, got: %d", testcase.args, testcase.statusCode, val)
			}
			var config = common.OpenConfig()
			var db2, _ = common.Open(config)
			for _, checker := range testcase.checkers {
				if db2.HasGroup(checker.groupName) != checker.existFlag {
					t.Errorf("%v: group check failed: %s, wont: %v, got: %v", testcase.args, checker.groupName, checker.existFlag, !checker.existFlag)
				}
				if checker.existFlag {
					var group = db2.FindGroup(checker.groupName)
					if group != nil && group.Description != checker.description {
						t.Errorf("%v: group description did not match: wont: %s, got: %s", testcase.args, checker.description, group.Description)
					}
					if group != nil && group.OmitList != checker.omitList {
						t.Errorf("%v: group OmitList did not match: wont: %v, got: %v", testcase.args, checker.omitList, group.OmitList)
					}
				}
			}
		})
	}
}

func TestUpdateGroupFailed(t *testing.T) {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")

	var testcases = []struct {
		opt     updateOptions
		errFlag bool
	}{
		{updateOptions{"newName", "desc", "omitList", "target"}, true},
	}
	for _, testcase := range testcases {
		common.Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
			var guc = groupUpdateCommand{}
			var config = common.OpenConfig()
			var db, _ = common.Open(config)
			var err = guc.updateGroup(db, &testcase.opt)
			if (err != nil) != testcase.errFlag {
				t.Errorf("%v: test failed: err wont: %v, got: %v: err (%v)", testcase.opt, testcase.errFlag, !testcase.errFlag, err)
			}
		})
	}
}

func TestUpdateGroup(t *testing.T) {
	type relation struct {
		groupID      string
		repositoryID string
		exist        bool
	}
	var testcases = []struct {
		args       []string
		statusCode int
		gexists    []groupChecker
		relations  []relation
	}{
		{[]string{"update", "-d", "newdesc2", "--name", "newgroup2", "-o", "true", "group2"}, 0,
			[]groupChecker{{"newgroup2", true, "newdesc2", true}, {"group2", false, "", false}},
			[]relation{}},
		{[]string{"update", "-n", "newgroup3", "group3"}, 0,
			[]groupChecker{{"newgroup3", true, "desc3", true}, {"group3", false, "desc3", false}},
			[]relation{{"newgroup3", "repo2", true}, {"group3", "repo2", false}},
		},
		{[]string{"update", "-o", "true", "group1"}, 0,
			[]groupChecker{{"group1", true, "desc1", true}},
			[]relation{{"group1", "repo1", true}},
		},
		{[]string{"update", "group4"}, 3, []groupChecker{}, []relation{}},
		{[]string{"update"}, 1, []groupChecker{}, []relation{}},
		{[]string{"update", "group1", "group4"}, 1, []groupChecker{}, []relation{}},
	}
	for _, testcase := range testcases {
		common.Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
			var guc, _ = GroupCommandFactory()
			if val := guc.Run(testcase.args); val != testcase.statusCode {
				t.Errorf("%v: group update failed status code wont: %d, got: %d", testcase.args, testcase.statusCode, val)
			}
			var config = common.OpenConfig()
			var db2, _ = common.Open(config)
			for _, gec := range testcase.gexists {
				if db2.HasGroup(gec.groupName) != gec.existFlag {
					t.Errorf("%s: exist check failed wont: %v, got: %v", gec.groupName, gec.existFlag, !gec.existFlag)
				}
				if gec.existFlag {
					var group = db2.FindGroup(gec.groupName)
					if group.Description != gec.description {
						t.Errorf("%s: description did not match: wont: %s, got: %s", gec.groupName, gec.description, group.Description)
					}
					if group.OmitList != gec.omitList {
						t.Errorf("%s: omitList did not match: wont: %v, got: %v", gec.groupName, gec.omitList, !gec.omitList)
					}
				}
			}
			for _, rel := range testcase.relations {
				if db2.HasRelation(rel.groupID, rel.repositoryID) != rel.exist {
					t.Errorf("relation between %s and %s: wont: %v, got: %v", rel.groupID, rel.repositoryID, rel.exist, !rel.exist)
				}
			}
		})
	}
}

func TestRemoveGroup(t *testing.T) {
	var testcases = []struct {
		args       []string
		statusCode int
		checkers   []groupChecker
	}{
		{[]string{"rm", "--force", "-v", "group1"}, 0, []groupChecker{{"group1", false, "", false}}},
		{[]string{"rm", "group2"}, 0, []groupChecker{{"group2", false, "desc2", false}}},
		{[]string{"rm", "group1"}, 3, []groupChecker{{"group1", true, "desc1", false}}},
		{[]string{"rm", "group4"}, 0, []groupChecker{{"group4", false, "not exist group", false}}},
		{[]string{"rm"}, 1, []groupChecker{}},
	}
	for _, testcase := range testcases {
		common.Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
			var grc, _ = GroupCommandFactory()
			if val := grc.Run(testcase.args); val != testcase.statusCode {
				t.Errorf("%v: group remove failed: wont: %d, got: %d", testcase.args, testcase.statusCode, val)
			}
			var config = common.OpenConfig()
			var db2, _ = common.Open(config)
			for _, checker := range testcase.checkers {
				if db2.HasGroup(checker.groupName) != checker.existFlag {
					t.Errorf("%v: exist check failed: wont: %v, got: %v", testcase.args, checker.existFlag, !checker.existFlag)
				}
				if checker.existFlag {
					var group = db2.FindGroup(checker.groupName)
					if group != nil && group.Description != checker.description {
						t.Errorf("%s: description did not match: wont: %s, got: %s", checker.groupName, checker.description, group.Description)
					}
					if group != nil && group.OmitList != checker.omitList {
						t.Errorf("%s: omitList did not match: wont: %v, got: %v", checker.groupName, checker.omitList, !checker.omitList)
					}
				}
			}
		})
	}
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
	var goc, _ = groupOfCommandFactory()
	var gc, _ = GroupCommandFactory()

	var gacHelp = `rrh group add [OPTIONS] <GROUPS...>
OPTIONS
    -d, --desc <DESC>        gives the description of the group.
    -o, --omit-list <FLAG>   gives the omit list flag of the group.
ARGUMENTS
    GROUPS                   gives group names.`

	var glcHelp = `rrh group list [OPTIONS]
OPTIONS
    -d, --desc             show description.
    -r, --repository       show repositories in the group.
    -o, --only-groupname   show only group name. This option is prioritized.`

	var grcHelp = `rrh group rm [OPTIONS] <GROUPS...>
OPTIONS
    -f, --force      force remove.
    -i, --inquiry    inquiry mode.
    -v, --verbose    verbose mode.
ARGUMENTS
    GROUPS           target group names.`

	var gucHelp = `rrh group update [OPTIONS] <GROUP>
OPTIONS
    -n, --name <NAME>        change group name to NAME.
    -d, --desc <DESC>        change description to DESC.
    -o, --omit-list <FLAG>   change omit-list of the group. FLAG must be "true" or "false".
ARGUMENTS
    GROUP               update target group names.`

	var gocHelp = `rrh group of <REPOSITORY_ID>
ARGUMENTS
    REPOSITORY_ID     show the groups of the repository.`

	var gcHelp = `rrh group <SUBCOMMAND>
SUBCOMMAND
    add       add new group.
    list      list groups (default).
    of        shows groups of the specified repository.
    rm        remove group.
    update    update group.`

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
	if goc.Help() != gocHelp {
		t.Error("help message did not match")
	}
	if grc.Help() != grcHelp {
		t.Error("help message did not match")
	}
}

func TestSynopsis(t *testing.T) {
	var gc, _ = GroupCommandFactory()
	if gc.Synopsis() != "add/list/update/remove groups and show groups of the repository." {
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
	var goc, _ = groupOfCommandFactory()
	if goc.Synopsis() != "show groups of the repository." {
		t.Error("synopsis did not match")
	}
}
