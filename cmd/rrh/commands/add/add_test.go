package add

import (
	"os"
	"testing"

	"github.com/tamada/rrh"
)

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
		hasError    bool
		gCheckers   []groupChecker
		rCheckers   []repositoryChecker
		relCheckers []relationChecker
	}{
		{[]string{"-g", "group2", "../../../../testdata/helloworld"}, false,
			[]groupChecker{{"group2", true}},
			[]repositoryChecker{{"helloworld", true}},
			[]relationChecker{{"group2", "helloworld", true}},
		},
		{[]string{"../../../../testdata/fibonacci"}, false,
			[]groupChecker{{"no-group", true}},
			[]repositoryChecker{{"fibonacci", true}},
			[]relationChecker{{"no-group", "fibonacci", true}},
		},
		{[]string{"../../../../testdata/fibonacci", "../../../../testdata/helloworld", "../../../../not-exist-dir", "../../../../testdata/other/helloworld"}, true,
			[]groupChecker{{"no-group", true}},
			[]repositoryChecker{{"fibonacci", true}, {"helloworld", true}, {"not-exist-dir", false}},
			[]relationChecker{{"no-group", "fibonacci", true}, {"no-group", "helloworld", true}},
		},
		{[]string{"../../../../testdata/helloworld", "../../../../testdata/other/helloworld"}, true,
			[]groupChecker{},
			[]repositoryChecker{{"helloworld", true}},
			[]relationChecker{{"no-group", "helloworld", true}},
		},
		{[]string{"--repository-id", "hw", "../../../../testdata/other/helloworld"}, false,
			[]groupChecker{},
			[]repositoryChecker{{"hw", true}},
			[]relationChecker{{"no-group", "hw", true}},
		},
		{[]string{"--repository-id", "fails", "../../../../testdata/other/helloworld", "../../../../testdata/fibonacci"}, true,
			[]groupChecker{},
			[]repositoryChecker{},
			[]relationChecker{},
		},
		{[]string{}, true, // too few arguments
			[]groupChecker{},
			[]repositoryChecker{},
			[]relationChecker{},
		},
	}

	os.Setenv(rrh.ConfigPath, "../../../../testdata/config.json")
	for _, td := range testcases {
		var databaseFile = rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
			cmd := New()
			cmd.SetArgs(td.args)
			err := cmd.Execute()

			if err == nil && td.hasError || err != nil && !td.hasError {
				t.Errorf("%v: wont error %v but got %v", td.args, td.hasError, err)
			}

			var db, _ = rrh.Open(config)
			for _, checker := range td.gCheckers {
				if db.HasGroup(checker.groupName) != checker.existFlag {
					t.Errorf("%v: group %s wont: %v, got: %v", td.args, checker.groupName, checker.existFlag, !checker.existFlag)
				}
			}
			for _, checker := range td.rCheckers {
				if db.HasRepository(checker.repositoryID) != checker.existFlag {
					t.Errorf("%v: repository %s wont: %v, got: %v", td.args, checker.repositoryID, checker.existFlag, !checker.existFlag)
				}
			}
			for _, checker := range td.relCheckers {
				if db.HasRelation(checker.groupName, checker.repositoryID) != checker.existFlag {
					t.Errorf("%v: relation (%s, %s) wont: %v, got: %v", td.args, checker.repositoryID, checker.groupName, checker.existFlag, !checker.existFlag)
				}
			}
		})
		defer os.Remove(databaseFile)
	}
}
