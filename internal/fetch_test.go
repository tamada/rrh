package internal

import (
	"os"
	"testing"

	"github.com/tamada/rrh"
)

func TestPerformFetch(t *testing.T) {
	var dbFile = rrh.Rollback("../testdata/database.json", "../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		var fetch, _ = FetchCommandFactory()
		var status = fetch.Run([]string{"no-group"})
		if status != 0 {
			t.Errorf("status did not match, wont 0, got %d", status)
		}
	})
	defer os.Remove(dbFile)
}

func TestHelpOfFetchCommand(t *testing.T) {
	var fetchHelp = `rrh fetch [OPTIONS] [GROUPS...]
OPTIONS
    -r, --remote <REMOTE>   specify the remote name. Default is "origin."
ARGUMENTS
    GROUPS                  run "git fetch" command on each repository on the group.
                            if no value is specified, run on the default group.`
	var fetchAllHelp = `rrh fetch-all [OPTIONS]
OPTIONS
    -r, --remote <REMOTE>   specify the remote name. Default is "origin."`

	var fetch, _ = FetchCommandFactory()
	var fetchAll, _ = FetchAllCommandFactory()

	if fetch.Help() != fetchHelp {
		t.Errorf("help message of fetch command did not match")
	}
	if fetchAll.Help() != fetchAllHelp {
		t.Errorf("help message of fetch_all command did not match")
	}
}

func TestProgress(t *testing.T) {
	var progress = NewProgress(3)

	if progress.String() != "  0/  3" {
		t.Errorf("string representation of progress did not match, wont   0/  3, got %s", progress.String())
	}
	progress.Increment()
	if progress.String() != "  1/  3" {
		t.Errorf("string representation of progress did not match, wont   1/  3, got %s", progress.String())
	}
	progress.Increment()
	if progress.String() != "  2/  3" {
		t.Errorf("string representation of progress did not match, wont   2/  3, got %s", progress.String())
	}
	progress.Increment()
	if progress.String() != "  3/  3" {
		t.Errorf("string representation of progress did not match, wont   3/  3, got %s", progress.String())
	}
	progress.Increment()
	if progress.String() != "  3/  3" {
		t.Errorf("string representation of progress did not match, wont   3/  3, got %s", progress.String())
	}

	var invalidProgress = NewProgress(-10)
	if invalidProgress.total != 0 {
		t.Errorf("total of progress should be the positive value")
	}
}
