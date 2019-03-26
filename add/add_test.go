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
		command.Run([]string{"../testdata/fibonacci"})

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
	})
}

func TestAddToDifferentGroup(t *testing.T) {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	rollback(func() {
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
