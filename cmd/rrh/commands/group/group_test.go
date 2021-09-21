package group

import (
	"os"

	"github.com/tamada/rrh"
)

func ExampleGroupCommand_Run() {
	dbFile := rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		cmd := New()
		cmd.SetArgs([]string{})
		cmd.SetOut(os.Stdout)
		cmd.Execute()
	})
	defer os.Remove(dbFile)
	// Output:
	// group1,1 repository
	// group2,0 repositories
	// group3,1 repository
}

func ExampleGroupListCommand_Run() {
	dbFile := rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		cmd := New()
		cmd.SetArgs([]string{"list", "--entry", "repo", "-e", "name,count,desc"})
		cmd.SetOut(os.Stdout)
		cmd.Execute()
	})
	defer os.Remove(dbFile)
	// Output:
	// group1,desc1,[repo1],1 repository
	// group2,desc2,[],0 repositories
	// group3,desc3,[repo2],1 repository
}

func ExampleGroupOfCommand_Run() {
	dbFile := rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		cmd := New()
		cmd.SetArgs([]string{"of", "repo1"})
		cmd.SetOut(os.Stdout)
		cmd.Execute()
	})
	defer os.Remove(dbFile)
	// Output:
	// repo1,[group1]
}

func ExampleGroupInfoCommand_Run() {
	dbFile := rrh.Rollback("../../../../testdata/test_db.json", "../../../../testdata/config.json", func(config *rrh.Config, db *rrh.Database) {
		cmd := New()
		cmd.SetArgs([]string{"info", "group1", "group2", "groupN"})
		cmd.SetOut(os.Stdout)
		cmd.Execute()
	})
	defer os.Remove(dbFile)
	// Output:
	// group1: desc1 (1 repository, abbrev: false)
	// group2: desc2 (0 repositories, abbrev: false)
	// groupN: group not found
}
