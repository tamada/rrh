package export

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/tamada/rrh/common"
)

func open(jsonName string) *common.Database {
	os.Setenv(common.RrhDatabasePath, fmt.Sprintf("../testdata/%s", jsonName))
	var config = common.OpenConfig()
	var db, _ = common.Open(config)
	return db
}

func TestHelpAndSynopsis(t *testing.T) {
	var export = ExportCommand{}
	var help = export.Help()
	var helpMessage = `rrh export [OPTIONS]
OPTIONS
    --no-indent      print result as no indented json,
    --no-hide-home   not replace home directory to '$HOME' keyword`
	if help != helpMessage {
		t.Error("help message did not match")
	}

	if export.Synopsis() != "export RRH database to stdout." {
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
	os.Setenv(common.RrhDatabasePath, "../testdata/broken.json")
	var export, _ = ExportCommandFactory()
	if val := export.Run([]string{}); val != 2 {
		t.Errorf("broken json successfully read!?: %d", val)
	}
}

func TestNullDB(t *testing.T) {
	os.Setenv(common.RrhDatabasePath, "../testdata/nulldb.json")
	var result, _ = common.CaptureStdout(func() {
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
	os.Setenv(common.RrhDatabasePath, "../testdata/nulldb.json")
	var result, _ = common.CaptureStdout(func() {
		var export, _ = ExportCommandFactory()
		export.Run([]string{"--no-indent"})
	})
	if strings.TrimSpace(result) != "{\"last_modified\":\"1970-01-01T09:00:00+09:00\",\"repositories\":[],\"groups\":[],\"relations\":[]}" {
		t.Errorf("nulldb data did not match: %s", result)
	}
}

func TestTmpDBNoIndent(t *testing.T) {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var result, _ = common.CaptureStdout(func() {
		var export, _ = ExportCommandFactory()
		export.Run([]string{"--no-indent"})
	})
	result = strings.TrimSpace(result)

	if !strings.HasPrefix(result, "{\"last_modified\":") &&
		!strings.HasSuffix(result, `"repositories":[{"repository_id":"repo1","repository_path":"path1","remotes":[]},{"repository_id":"repo2","repository_path":"path2","remotes":[]}],"groups":[{"group_name":"group1","group_desc":"desc1"},{"group_name":"group2","group_desc":"desc2"}],"relations":[{"repository_id":"repo1","group_name":"group1"}]}`) {
		t.Errorf("tmp.json was not matched.\ngot: %s", result)
	}
	// In example testing, how do I ignore the part of output, like below?
	// Output:
	// {"last_modified":".*",repositories":[{"repository_id":"repo1","repository_path":"path1","remotes":[]},{"repository_id":"repo2","repository_path":"path2","remotes":[]}],"groups":[{"group_name":"group1","group_desc":"desc1","group_items":["repo1"]},{"group_name":"group2","group_desc":"desc2","group_items":[]}]}
}
