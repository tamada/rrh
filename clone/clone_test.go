package clone

import (
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

func TestCloneCommand_Run(t *testing.T) {
	os.Setenv(common.RrhConfigPath, "../testdata/config.json")
	os.Setenv(common.RrhDatabasePath, "../testdata/tmp.json")
	rollback(func() {
		var clone, _ = CloneCommandFactory()
		clone.Run([]string{"-d", "../testdata", "-v", "https://htamada@bitbucket.org/htamada/helloworld.git"})
		defer cleanup([]string{"../testdata/helloworld"})

		var config = common.OpenConfig()
		var db, _ = common.Open(config)
		if len(db.Repositories) != 3 {
			t.Fatal("helloworld was not registered.")
		}
		var repo = db.Repositories[2]
		var dir, _ = filepath.Abs("../testdata/helloworld")
		if repo.ID != "helloworld" || repo.Path != dir {
			t.Errorf("wont: helloworld (%s), got: %s (%s)", dir, repo.ID, repo.Path)
		}
		var stat, err = os.Stat(dir)
		if os.IsNotExist(err) || !stat.IsDir() {
			t.Errorf("%s not exist or not dir", dir)
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

	if clone.Synopsis() != `run "git clone"` {
		t.Error("synopsis did not match")
	}
}
