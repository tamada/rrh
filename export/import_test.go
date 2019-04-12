package export

import (
	"os"
	"strings"
	"testing"

	"github.com/tamada/rrh/common"
)

func TestImport(t *testing.T) {
	type groupExistCheck struct {
		groupName string
		wontExist bool
	}
	type repoExistCheck struct {
		repositoryName string
		wontExist      bool
		path           string
	}
	type relationCheck struct {
		groupName      string
		repositoryName string
		wontRelation   bool
	}

	var testcases = []struct {
		args       []string
		statusCode int
		gChecks    []groupExistCheck
		rChecks    []repoExistCheck
		relChecks  []relationCheck
	}{
		{[]string{"../testdata/exported.json"}, 0,
			[]groupExistCheck{{"test", true}, {"group1", true}, {"group2", true}, {"group3", true}},
			[]repoExistCheck{{"repo1", true, "path1"}, {"repo2", true, "path2"}, {"fibonacci", true, "../testdata/fibonacci"}, {"helloworld", false, "../testdata/helloworld2"}},
			[]relationCheck{{"group1", "repo1", true}, {"group3", "repo2", true}, {"test", "fibonacci", true}, {"test", "helloworld", false}}},
		{[]string{"--overwrite", "../testdata/exported.json"}, 0,
			[]groupExistCheck{{"test", true}, {"group1", false}, {"group2", false}, {"group3", false}},
			[]repoExistCheck{{"repo1", false, "path1"}, {"repo2", false, "path2"}, {"fibonacci", true, "../testdata/fibonacci"}, {"helloworld", false, "../testdata/helloworld2"}},
			[]relationCheck{{"group1", "repo1", false}, {"group3", "repo2", false}, {"test", "fibonacci", true}, {"test", "helloworld", false}}},
		{[]string{"--auto-clone", "--overwrite", "--verbose", "../testdata/exported.json"}, 0,
			[]groupExistCheck{{"test", true}, {"group1", false}, {"group2", false}, {"group3", false}},
			[]repoExistCheck{{"repo1", false, "path1"}, {"repo2", false, "path2"}, {"fibonacci", true, "../testdata/fibonacci"}, {"helloworld", true, "../testdata/helloworld2"}},
			[]relationCheck{{"group1", "repo1", false}, {"group3", "repo2", false}, {"test", "fibonacci", true}, {"test", "helloworld", true}}},
	}

	for _, testcase := range testcases {
		os.Setenv(common.RrhConfigPath, "../testdata/config.json")
		var dbFile = common.Rollback("../testdata/tmp.json", "../testdata/config.json", func() {
			var command, _ = ImportCommandFactory()
			var statusCode = command.Run(testcase.args)
			if statusCode != testcase.statusCode {
				t.Errorf("%v: status code did not match: wont: %d, got: %d", testcase.args, testcase.statusCode, statusCode)
			}

			var db, _ = common.Open(common.OpenConfig())
			for _, gcheck := range testcase.gChecks {
				if db.HasGroup(gcheck.groupName) != gcheck.wontExist {
					t.Errorf("%v: group %s exist: wont: %v, got: %v", testcase.args, gcheck.groupName, gcheck.wontExist, !gcheck.wontExist)
				}
			}
			for _, rcheck := range testcase.rChecks {
				if db.HasRepository(rcheck.repositoryName) != rcheck.wontExist {
					t.Errorf("%v: repository %s exist: wont: %v, got: %v", testcase.args, rcheck.repositoryName, rcheck.wontExist, !rcheck.wontExist)
				}
				if db.HasRepository(rcheck.repositoryName) {
					var repo = db.FindRepository(rcheck.repositoryName)
					if repo.Path != rcheck.path {
						t.Errorf("%v: repository %s path did not match: wont: %s, got: %s", testcase.args, rcheck.repositoryName, rcheck.path, repo.Path)
					}
				}
			}
			for _, relcheck := range testcase.relChecks {
				if db.HasRelation(relcheck.groupName, relcheck.repositoryName) != relcheck.wontRelation {
					t.Errorf("%v: relation g(%s) and r(%s): wont: %v, got: %v", testcase.args, relcheck.groupName, relcheck.repositoryName, relcheck.wontRelation, !relcheck.wontRelation)
				}
			}
		})
		defer os.Remove(dbFile)
	}
	defer os.RemoveAll("../testdata/helloworld2")
}

func TestParsingFailOfArgs(t *testing.T) {
	var testcases = []struct {
		args []string
		wont string
	}{
		{[]string{}, "too few arguments"},
		{[]string{"a", "b", "c"}, "too many arguments: [a b c]"},
		{[]string{"-unknown-option"}, `rrh import [OPTIONS] <DATABASE_JSON>
OPTIONS
    --auto-clone    clone the repository, if paths do not exist.
    --overwrite     replace the local RRH database to the given database.
    -v, --verbose   verbose mode.
ARGUMENTS
    DATABASE_JSON   the exported RRH database.
flag provided but not defined: -unknown-option`},
	}

	for _, testcase := range testcases {
		var got = common.CaptureStdout(func() {
			var command, _ = ImportCommandFactory()
			command.Run(testcase.args)
		})
		if strings.TrimSpace(got) != testcase.wont {
			t.Errorf("%v: failed, wont: %s, got: %s", testcase.args, testcase.wont, got)
		}
	}
}

func TestSynopsisAndHelp(t *testing.T) {
	var command, _ = ImportCommandFactory()
	if command.Synopsis() != "import the given database." {
		t.Error("Synopsis did not match")
	}
	var helpMessage = `rrh import [OPTIONS] <DATABASE_JSON>
OPTIONS
    --auto-clone    clone the repository, if paths do not exist.
    --overwrite     replace the local RRH database to the given database.
    -v, --verbose   verbose mode.
ARGUMENTS
    DATABASE_JSON   the exported RRH database.`
	if command.Help() != helpMessage {
		t.Error("Help did not match")
	}
}
