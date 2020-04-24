package internal

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/tamada/rrh"
)

func open2(jsonName string) *rrh.Database {
	os.Setenv(rrh.DatabasePath, fmt.Sprintf("../testdata/%s", jsonName))
	var config = rrh.OpenConfig()
	var db, _ = rrh.Open(config)
	return db
}

func TestHelpAndSynopsisOfExportCommand(t *testing.T) {
	var export = ExportCommand{}
	var help = export.Help()
	var helpMessage = `rrh export [OPTIONS]
OPTIONS
    --no-indent      print result as no indented json
    --no-hide-home   not replace home directory to '${HOME}' keyword`
	if help != helpMessage {
		t.Error("help message did not match")
	}

	if export.Synopsis() != "export rrh database to stdout." {
		t.Error("synopsis did not match")
	}
}

func TestUnknownOptions(t *testing.T) {
	var export, _ = ExportCommandFactory()
	if export.Run([]string{"--unknown-option"}) != 1 {
		t.Error("--unknown-option was not failed.")
	}
}

func TestBrokenDatabase(t *testing.T) {
	var dbFile = rrh.Rollback("../testdata/broken.json", "../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
		var export, _ = ExportCommandFactory()
		if val := export.Run([]string{}); val != 2 {
			t.Errorf("broken json successfully read!?: %d", val)
		}
	})
	defer os.Remove(dbFile)
}

func TestNullDB(t *testing.T) {
	os.Setenv(rrh.DatabasePath, "../testdata/nulldb.json")
	var result = rrh.CaptureStdout(func() {
		var export, _ = ExportCommandFactory()
		export.Run([]string{})
	})
	var actually = `{
  "last_modified": "1970-01-01T09:00:00+09:00",
  "repositories": [],
  "groups": [],
  "relations": []
}`
	if strings.TrimSpace(result) != actually {
		t.Errorf("nulldb data did not match: wont: %s, got: %s", actually, strings.TrimSpace(result))
	}
}

func TestNullDBNoIndent(t *testing.T) {
	os.Setenv(rrh.DatabasePath, "../testdata/nulldb.json")
	var result = rrh.CaptureStdout(func() {
		var export, _ = ExportCommandFactory()
		export.Run([]string{"--no-indent"})
	})
	if strings.TrimSpace(result) != "{\"last_modified\":\"1970-01-01T09:00:00+09:00\",\"repositories\":[],\"groups\":[],\"relations\":[]}" {
		t.Errorf("nulldb data did not match: %s", result)
	}
}

func TestTmpDBNoIndent(t *testing.T) {
	var dbFile = rrh.Rollback("../testdata/test_db.json", "../testdata/config.json", func(config *rrh.Config, oldDB *rrh.Database) {
		var result = rrh.CaptureStdout(func() {
			var export, _ = ExportCommandFactory()
			export.Run([]string{"--no-indent"})
		})
		result = strings.TrimSpace(result)

		if !strings.HasPrefix(result, "{\"last_modified\":") ||
			!strings.HasSuffix(result, `"repositories":[{"repository_id":"repo1","repository_path":"path1","repository_desc":"","remotes":[]},{"repository_id":"repo2","repository_path":"path2","repository_desc":"","remotes":[{"name":"origin","url":"git@github.com:example/repo2.git"}]}],"groups":[{"group_name":"group1","group_desc":"desc1","omit_list":false},{"group_name":"group2","group_desc":"desc2","omit_list":false},{"group_name":"group3","group_desc":"desc3","omit_list":true}],"relations":[{"repository_id":"repo1","group_name":"group1"},{"repository_id":"repo2","group_name":"group3"}]}`) {
			t.Errorf("test_db.json was not matched.\ngot: %s", result)
		}
	})
	// In example testing, how do I ignore the part of output, like below?
	// Output:
	// {"last_modified":".*",repositories":[{"repository_id":"repo1","repository_path":"path1","remotes":[]},{"repository_id":"repo2","repository_path":"path2","remotes":[]}],"groups":[{"group_name":"group1","group_desc":"desc1","group_items":["repo1"]},{"group_name":"group2","group_desc":"desc2","group_items":[]}]}
	defer os.Remove(dbFile)
}
