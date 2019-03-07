package add

import (
	"os"
	"testing"

	"github.com/tamada/rrh/common"
)

func rollback(f func()) {
	var config = common.OpenConfig()
	var db, _ = common.Open(config)
	defer db.StoreAndClose()

	f()
}

func TestHelpAndSynopsis(t *testing.T) {
	var command, _ = AddCommandFactory()
	if command.Synopsis() != "add repositories on the local path to RRH" {
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

func TestAddToTheSpecifiedGroup(t *testing.T) {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	rollback(func() {
		var command, _ = AddCommandFactory()
		command.Run([]string{"--group", "group2", "../testdata/helloworld"})

		var config = common.OpenConfig()
		var db, _ = common.Open(config)
		if !db.HasGroup("group2") {
			t.Error("group2: group not found")
		}
		if !db.HasRepository("helloworld") {
			t.Error("helloworld: repository not found")
		}
		if !db.HasRelation("group2", "helloworld") {
			t.Error("gruop2, and helloworld: the relation not found")
		}
	})
}

func TestAddCommand_Run(t *testing.T) {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	rollback(func() {
		var command, _ = AddCommandFactory()
		command.Run([]string{"../testdata/helloworld"})

		var config = common.OpenConfig()
		var db, _ = common.Open(config)
		if !db.HasGroup("no-group") {
			t.Error("no-group: group not found")
		}
		if !db.HasRepository("helloworld") {
			t.Error("helloworld: repository not found")
		}
		if !db.HasRelation("no-group", "helloworld") {
			t.Error("no-group, and helloworld: the relation not found")
		}
	})
}
