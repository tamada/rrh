package export

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

func TestHelp(t *testing.T) {
	var export = ExportCommand{}
	var help = export.Help()
	var helpMessage = `rrh export [OPTIONS]
OPTIONS
    --no-indent    print result as no indented json (Default indented json)`
	if help != helpMessage {
		t.Error("help message was not match")
	}
}

func ExampleNullDB() {
	os.Setenv(common.RrhDatabasePath, "../testdata/nulldb.json")
	var export = ExportCommand{}
	export.Run([]string{})
	// Output:
	// {
	//   "last_modified": "0001-01-01T00:00:00Z",
	//   "repositories": [],
	//   "groups": []
	// }
}

func ExampleNullDBNoIndent() {
	os.Setenv(common.RrhDatabasePath, "../testdata/nulldb.json")
	var export = ExportCommand{}
	export.Run([]string{"--no-indent"})
	// Output: {"last_modified":"0001-01-01T00:00:00Z","repositories":[],"groups":[]}
}

func ExampleTmpNoIndent() {
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	var export = ExportCommand{}
	export.Run([]string{"--no-indent"})
	// Output: {"last_modified":"2019-02-22T17:22:14.055153+09:00","repositories":[{"repository_id":"repo1","repository_path":"path1","remotes":[]},{"repository_id":"repo2","repository_path":"path2","remotes":[]}],"groups":[{"group_name":"group1","group_desc":"desc1","group_items":["repo1"]},{"group_name":"group2","group_desc":"desc2","group_items":[]}]}
}
