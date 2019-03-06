package clone

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/tamada/rrh/common"
)

func cleanup(dirs []string) {
	for _, dir := range dirs {
		os.RemoveAll(dir)
	}
}

func rollback(f func()) {
	var config = common.OpenConfig()
	var db, _ = common.Open(config)
	defer db.StoreAndClose()

	f()
}

func validate(repo common.Repository, repoID string, repoPath string) string {
	var dir, _ = filepath.Abs(repoPath)
	if repo.ID != repoID || repo.Path != dir {
		return fmt.Sprintf("wont: %s (%s), got: %s (%s)", repoID, dir, repo.ID, repo.Path)
	}
	var stat, err = os.Stat(dir)
	if os.IsNotExist(err) || !stat.IsDir() {
		return fmt.Sprintf("%s not exist or not dir", dir)
	}
	return ""
}

func contains(slice []string, checkItem string) bool {
	for _, item := range slice {
		if item == checkItem {
			return true
		}
	}
	return false
}

func TestCloneCommand_MultipleProjects(t *testing.T) {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")

	rollback(func() {
		var clone, _ = CloneCommandFactory()
		clone.Run([]string{"-d", "../testdata/hoge", "-g", "not-exist-group",
			"https://htamada@bitbucket.org/htamada/helloworld.git",
			"https://htamada@bitbucket.org/htamada/fibonacci.git"})
		defer cleanup([]string{"../testdata/hoge"})

		var config = common.OpenConfig()
		var db, _ = common.Open(config)
		if len(db.Repositories) != 4 {
			t.Fatal("helloworld and fibonacci were not registered.")
		}
		var hwRepo = db.Repositories[2]
		if message := validate(hwRepo, "helloworld", "../testdata/hoge/helloworld"); message != "" {
			t.Error(message)
		}
		var fiboRepo = db.Repositories[3]
		if message := validate(fiboRepo, "fibonacci", "../testdata/hoge/fibonacci"); message != "" {
			t.Error(message)
		}
		if !db.HasGroup("not-exist-group") || len(db.Groups) != 3 {
			t.Fatalf("not-exist-group: group not found: %v", db.Groups)
		}
		var group = db.FindGroup("not-exist-group")
		if !contains(group.Items, "helloworld") || !contains(group.Items, "fibonacci") {
			t.Errorf("%s: does not have helloworld or fibonacci", group.Name)
		}
	})
}

func TestCloneCommand_Run(t *testing.T) {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	rollback(func() {
		var clone, _ = CloneCommandFactory()
		clone.Run([]string{"-dest", "../testdata", "https://htamada@bitbucket.org/htamada/helloworld.git"})
		defer cleanup([]string{"../testdata/helloworld"})

		var config = common.OpenConfig()
		var db, _ = common.Open(config)
		if len(db.Repositories) != 3 {
			t.Fatal("helloworld was not registered.")
		}
		var repo = db.Repositories[2]
		if message := validate(repo, "helloworld", "../testdata/helloworld"); message != "" {
			t.Error(message)
		}
		var group = db.FindGroup("no-group")
		if len(group.Items) != 1 || !contains(group.Items, "helloworld") {
			t.Errorf("helloworld was not registered to the group \"no-group\": %v", group.Items)
		}
	})
}

func TestCloneCommand_SpecifyingId(t *testing.T) {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	rollback(func() {
		var clone, _ = CloneCommandFactory()
		clone.Run([]string{"-d", "../testdata/newid", "https://htamada@bitbucket.org/htamada/helloworld.git"})
		defer cleanup([]string{"../testdata/newid"})

		var config = common.OpenConfig()
		var db, _ = common.Open(config)
		if len(db.Repositories) != 3 {
			t.Fatal("newid was not registered.")
		}
		var repo = db.Repositories[2]
		if message := validate(repo, "newid", "../testdata/newid"); message != "" {
			t.Error(message)
		}
	})
}

func TestHelpAndSynopsis(t *testing.T) {
	var helpMessage = `rrh clone [OPTIONS] <REMOTE_REPOS...>
OPTIONS
    -g, --group <GROUP>   print managed repositories categoried in the group.
    -d, --dest <DEST>     specify the destination.
    -v, --verbose         verbose mode.
ARGUMENTS
    REMOTE_REPOS          repository urls`

	var clone, _ = CloneCommandFactory()
	if clone.Help() != helpMessage {
		t.Error("help message did not match")
	}

	if clone.Synopsis() != `run "git clone" and register it to a group` {
		t.Error("synopsis did not match")
	}
}
