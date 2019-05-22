package add

import (
	"os"
	"testing"

	"github.com/tamada/rrh/common"
)

func TestInvalidOptions(t *testing.T) {
	common.CaptureStdout(func() {
		var command, _ = CommandFactory()
		var flag = command.Run([]string{"--invalid-option"})
		if flag != 1 {
			t.Errorf("parse option failed.")
		}
	})
}

func TestHelpAndSynopsis(t *testing.T) {
	var command, _ = CommandFactory()
	if command.Synopsis() != "add repositories on the local path to RRH." {
		t.Error("synopsis did not match")
	}
	if command.Help() != `rrh add [OPTIONS] <REPOSITORY_PATHS...>
OPTIONS
    -g, --group <GROUP>        add repository to RRH database.
    -r, --repository-id <ID>   specified repository id of the given repository path.
                               Specifying this option fails with multiple arguments.
ARGUMENTS
    REPOSITORY_PATHS           the local path list of the git repositories.` {
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
		{[]string{"-g", "group2", "../testdata/helloworld"}, 0,
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
		{[]string{"--repository-id", "hw", "../testdata/other/helloworld"}, 0,
			[]groupChecker{},
			[]repositoryChecker{{"hw", true}},
			[]relationChecker{{"no-group", "hw", true}},
		},
		{[]string{"--repository-id", "fails", "../testdata/other/helloworld", "../testdata/fibonacci"}, 0,
			[]groupChecker{},
			[]repositoryChecker{},
			[]relationChecker{},
		},
	}

	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	for _, testcase := range testcases {
		var databaseFile = common.Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
			var command, _ = CommandFactory()
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
		defer os.Remove(databaseFile)
	}
}

func TestAddToDifferentGroup(t *testing.T) {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	var databaseFile = common.Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
		var command, _ = CommandFactory()
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
	defer os.Remove(databaseFile)
}

func TestAddFailed(t *testing.T) {
	os.Setenv(common.RrhConfigPath, "../testdata/nulldb.json")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	os.Setenv(common.RrhAutoCreateGroup, "false")

	var add = Command{}
	var config = common.OpenConfig()
	var db, _ = common.Open(config)

	var data = []options{
		{args: []string{"../not-exist-dir"}, group: "no-group"},
		{args: []string{"../testdata/fibonacci"}, group: "not-exist-group"},
	}

	for _, datum := range data {
		var list = add.AddRepositoriesToGroup(db, &datum)
		if len(list) == 0 {
			t.Errorf("successfully add in invalid data: %v", datum)
		}
	}
}

func TestFindRemotes(t *testing.T) {
	var testdata = []struct {
		path      string
		errorFlag bool
		count     int
	}{
		{"../testdata/dummygit", true, 0},
	}
	for _, td := range testdata {
		var remotes, err = FindRemotes(td.path)
		if (err == nil) == td.errorFlag {
			t.Errorf("%s: error flag did not match, wont: %v, got: %v, %v", td.path, td.errorFlag, !td.errorFlag, err)
		}
		if err != nil && td.count != len(remotes) {
			t.Errorf("%s: remote count did not match, wont: %d, got: %d", td.path, td.count, len(remotes))
		}
	}
}
