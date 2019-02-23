package prune

import (
	"fmt"
	"os"
	"testing"

	"github.com/tamada/rrh/common"
)

func open(jsonName string) *common.Database {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	os.Setenv(common.RrhDatabasePath, fmt.Sprintf("../testdata/%s", jsonName))
	var config = common.OpenConfig()
	var db, _ = common.Open(config)
	return db
}

func TestSynopsis(t *testing.T) {
	var prune, _ = PruneCommandFactory()
	if prune.Synopsis() != "prune unnecessary repositories and groups." {
		t.Error("Synopsis message is not matched.")
	}
}
func TestHelp(t *testing.T) {
	var prune = PruneCommand{}
	if prune.Help() != "rrh prune" {
		t.Error("Help message is not matched.")
	}
}

func TestPrune(t *testing.T) {
	var db = open("tmp.json")
	db.Prune()
	if len(db.Repositories) != 1 && len(db.Groups) != 1 {
		t.Error("prune failed")
	}
}

func TestTruePrune(t *testing.T) {
	var db = open("tmp.json")
	var prune = PruneCommand{}
	prune.perform(db)

	if len(db.Repositories) != 0 && len(db.Groups) != 0 {
		t.Error("prune failed")
	}
}

func TestPruneCommandRunFailedByBrokenDBFile(t *testing.T) {
	os.Setenv(common.RrhDatabasePath, "../testdata/broken.json")
	var prune, _ = PruneCommandFactory()
	if prune.Run([]string{}) != 1 {
		t.Error("broken database read successfully.")
	}
}

func ExamplePruneCommand_Run() {
	var db = open("tmp.json")

	var prune, _ = PruneCommandFactory()
	prune.Run([]string{})
	// Output: Pruned 2 groups, 2 repositories

	db.StoreAndClose()
}
