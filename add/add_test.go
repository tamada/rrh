package add

import (
	"os"
	"testing"

	"github.com/tamada/rrh/common"
)

func TestInvalidOptions(t *testing.T) {
	common.CaptureStdout(func() {
		var command, _ = AddCommandFactory()
		var flag = command.Run([]string{"--invalid-option"})
		if flag != 1 {
			t.Errorf("parse option failed.")
		}
	})
}

func TestHelpAndSynopsis(t *testing.T) {
	var command, _ = AddCommandFactory()
	if command.Synopsis() != "add repositories on the local path to RRH." {
		t.Error("synopsis did not match")
	}
	if command.Help() != `rrh add [OPTIONS] <REPOSITORY_PATHS...>
OPTIONS
    -g, --group <GROUP>    add repository to RRH database.
ARGUMENTS
    REPOSITORY_PATHS       the local path list of the git repositories` {
		t.Error("help did not match")
	}
}

func TestAdd(t *testing.T) {
	type groupChecker struct {
		groupName string
		existFlag bool
	}
	type repositoryChecker struct {
		repositoryID string
		existFlag    bool
	}
	type relationChecker struct {
		groupName    string
		repositoryID string
		existFlag    bool
	}
	var testcases = []struct {
		args        []string
		statusCode  int
		gCheckers   []groupChecker
		rCheckers   []repositoryChecker
		relCheckers []relationChecker
	}{
		{[]string{"--group", "group2", "../testdata/helloworld"}, 0,
			[]groupChecker{{"group2", true}},
			[]repositoryChecker{{"helloworld", true}},
			[]relationChecker{{"group2", "helloworld", true}},
		},
		{[]string{"../testdata/fibonacci"}, 0,
			[]groupChecker{{"no-group", true}},
			[]repositoryChecker{{"fibonacci", true}},
			[]relationChecker{{"no-group", "fibonacci", true}},
		},
		{[]string{"../testdata/fibonacci", "../testdata/helloworld", "../not-exist-dir", "../testdata/other/helloworld"}, 0,
			[]groupChecker{{"no-group", true}},
			[]repositoryChecker{{"fibonacci", true}, {"helloworld", true}, {"not-exist-dir", false}},
			[]relationChecker{{"no-group", "fibonacci", true}, {"no-group", "helloworld", true}},
		},
		{[]string{"../testdata/helloworld", "../testdata/other/helloworld"}, 0,
			[]groupChecker{},
			[]repositoryChecker{{"helloworld", true}},
			[]relationChecker{{"no-group", "helloworld", true}},
		},
	}

	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	for _, testcase := range testcases {
		common.Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
			var command, _ = AddCommandFactory()
			var status = command.Run(testcase.args)

			var config = common.OpenConfig()
			var db, _ = common.Open(config)
			if status != testcase.statusCode {
				t.Errorf("%v: status code did not match, wont: %d, got: %d", testcase.args, testcase.statusCode, status)
			}

			for _, checker := range testcase.gCheckers {
				if db.HasGroup(checker.groupName) != checker.existFlag {
					t.Errorf("%v: group wont: %v, got: %v", testcase.args, checker.existFlag, !checker.existFlag)
				}
			}
			for _, checker := range testcase.rCheckers {
				if db.HasRepository(checker.repositoryID) != checker.existFlag {
					t.Errorf("%v: repository wont: %v, got: %v", testcase.args, checker.existFlag, !checker.existFlag)
				}
			}
			for _, checker := range testcase.relCheckers {
				if db.HasRelation(checker.groupName, checker.repositoryID) != checker.existFlag {
					t.Errorf("%v: relation (%s, %s) wont: %v, got: %v", testcase.args, checker.repositoryID, checker.groupName, checker.existFlag, !checker.existFlag)
				}
			}
		})
	}
}

func TestAddToDifferentGroup(t *testing.T) {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	common.Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
		var command, _ = AddCommandFactory()
		command.Run([]string{"../testdata/fibonacci"})
		command.Run([]string{"-g", "group1", "../testdata/fibonacci"})

		var config = common.OpenConfig()
		var db, _ = common.Open(config)
		if !db.HasGroup("no-group") {
			t.Error("no-group: group not found")
		}
		if !db.HasRepository("fibonacci") {
			t.Error("fibonacci: repository not found")
		}
		if !db.HasRelation("no-group", "fibonacci") {
			t.Error("no-group, and fibonacci: the relation not found")
		}
		if !db.HasRelation("group1", "fibonacci") {
			t.Error("group1 and fibonacci: the relation not found")
		}
	})
}

func TestAddFailed(t *testing.T) {
	os.Setenv(common.RrhConfigPath, "../testdata/nulldb.json")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	os.Setenv(common.RrhAutoCreateGroup, "false")

	var add = AddCommand{}
	var config = common.OpenConfig()
	var db, _ = common.Open(config)

	var data = []struct {
		args      []string
		groupName string
	}{
		{[]string{"../not-exist-dir"}, "no-group"},
		{[]string{"../testdata/fibonacci"}, "not-exist-group"},
	}

	for _, datum := range data {
		var list = add.AddRepositoriesToGroup(db, datum.args, datum.groupName)
		if len(list) == 0 {
			t.Errorf("successfully add in invalid data: %v", datum)
		}
	}
}
